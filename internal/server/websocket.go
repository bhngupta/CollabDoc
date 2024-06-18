package server

import (
	"CollabDoc/pkg/document"
	"CollabDoc/pkg/persistence"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type string             `json:"type"`
	Op   document.Operation `json:"operation,omitempty"`
}

var (
	ss                *document.StateSynchronizer
	persist           *persistence.Persistence
	operations        = make(map[string][]document.Operation)
	pendingOperations = make(map[string][]document.Operation)
	operationsMutex   sync.Mutex
)

func init() {
	persist = persistence.NewPersistence("state.json")
	var err error
	ss, err = persist.LoadState()
	if err != nil {
		fmt.Println("Error loading state:", err)
		ss = document.NewStateSynchronizer()
	} else {
		// Ensure ConflictResolver is initialized even when loading from persistence
		if ss.ConflictResolver == nil {
			ss.ConflictResolver = &document.ConflictResolver{}
		}
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer ws.Close()

	var docID string

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		if docID == "" && msg.Op.DocID != "" {
			docID = msg.Op.DocID
			fmt.Printf("Client connected from %s with Document ID: %s\n", ws.RemoteAddr(), docID)
		}

		switch msg.Type {
		case "operation":
			fmt.Printf("Received operation for Document ID %s: %+v\n", docID, msg.Op)

			operationsMutex.Lock()
			doc := ss.Documents[msg.Op.DocID]
			if doc == nil {
				fmt.Println("Document not found:", msg.Op.DocID)
				operationsMutex.Unlock()
				break
			}

			// Initialize operations slice if it doesn't exist
			if _, ok := operations[msg.Op.DocID]; !ok {
				operations[msg.Op.DocID] = []document.Operation{}
			}

			// Initialize pendingOperations slice if it doesn't exist
			if _, ok := pendingOperations[msg.Op.DocID]; !ok {
				pendingOperations[msg.Op.DocID] = []document.Operation{}
			}

			fmt.Printf("Operations array before this: %+v\n", operations)
			// Ensure the base version is within valid range
			if msg.Op.BaseVersion < 1 || msg.Op.BaseVersion > len(operations[msg.Op.DocID])+1 {
				fmt.Println("Invalid BaseVersion range:", msg.Op.BaseVersion, "length:", len(operations[msg.Op.DocID]))
				operationsMutex.Unlock()
				break
			}

			fmt.Printf("BaseVersion: %d, Document Content: %s, Operations: %+v\n", msg.Op.BaseVersion, doc.Content, operations[msg.Op.DocID])
			// Queue the operation
			pendingOperations[msg.Op.DocID] = append(pendingOperations[msg.Op.DocID], msg.Op)
			operationsMutex.Unlock()

			// Process the operations in the queue
			processPendingOperations(doc, msg.Op.DocID)
			doc, docExists := ss.Documents[msg.Op.DocID]
			if !docExists || doc == nil {
				// Create document if it does not exist
				doc = ss.CreateDocument(msg.Op.DocID)
				ss.Documents[msg.Op.DocID] = doc
				fmt.Printf("Document with ID %s created.\n", msg.Op.DocID)
			}

			opsSinceBase := operations[msg.Op.DocID][msg.Op.BaseVersion:]
			transformedOp := HandleOperation(doc, msg.Op, opsSinceBase)
			operations[msg.Op.DocID] = append(operations[msg.Op.DocID], transformedOp)
			operationsMutex.Unlock()

			// Log the transformed operation
			fmt.Printf("Transformed operation for Document ID %s: %+v\n", docID, transformedOp)

			// Broadcast the transformed operation
			//BroadcastOperation(transformedOp)

			// Save state
			err = persist.SaveState(ss)
			if err != nil {
				fmt.Println("Error saving state:", err)
			}

			if err != nil {
				fmt.Println("Error saving state:", err)
			}

			// Respond with success and the updated version
			err = ws.WriteJSON(map[string]interface{}{"success": true, "version": doc.Version})
		case "create":
			doc := ss.CreateDocument(msg.Op.DocID)
			fmt.Printf("Created document: %+v\n", doc)
			// Initialize the operations slice for the new document
			operationsMutex.Lock()
			operations[msg.Op.DocID] = []document.Operation{}
			pendingOperations[msg.Op.DocID] = []document.Operation{}
			operationsMutex.Unlock()

			fmt.Printf("Created document: %+v\n", doc)
			err = ws.WriteJSON(doc)

		case "get":
			doc, exists := ss.GetDocument(msg.Op.DocID)
			if exists {
				err = ws.WriteJSON(map[string]interface{}{
					"id":      docID,
					"content": doc.Content, // Adjust according to your document structure
				})
				if err != nil {
					fmt.Println("Error sending document content:", err)
				}
			} else {
				err = ws.WriteJSON(map[string]bool{"exists": false})
				if err != nil {
					fmt.Println("Error sending document not found:", err)
				}
			}

		default:
			fmt.Println("Unknown message type:", msg.Type)
		}

		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func processPendingOperations(doc *document.Document, docID string) {
	operationsMutex.Lock()
	defer operationsMutex.Unlock()

	fmt.Printf("Processing operations for document ID: %s\n", docID)
	if ops, ok := pendingOperations[docID]; ok {
		for _, op := range ops {
			opsSinceBase := operations[docID][op.BaseVersion-1:] // Include all operations since base version
			fmt.Printf("Ops since base: %+v\n", opsSinceBase)
			transformedOp := HandleOperation(doc, op, opsSinceBase)
			// Append the transformed operation to the history
			operations[docID] = append(operations[docID], transformedOp)
			// Print document state after each operation
			fmt.Printf("Document state after applying operation: %+v\n", doc)
		}
		// Clear the pending operations queue
		pendingOperations[docID] = []document.Operation{}
	}
	// Print final document state after all operations
	fmt.Printf("Final document state: %+v\n", doc)
}

func HandleOperation(doc *document.Document, op document.Operation, opsSinceBase []document.Operation) document.Operation {
	fmt.Printf("Handling operation: %+v\n", op)
	for _, prevOp := range opsSinceBase {
		op = document.TransformOperation(op, prevOp)
	}
	err := document.ApplyOperation(doc, op)
	if err != nil {
		fmt.Println("ApplyOperation error:", err)
		return op
	}
	doc.Version++ // Increment the document version after applying the operation
	fmt.Printf("Document version after operation: %d\n", doc.Version)
	fmt.Printf("Document state after applying operation: %+v\n", doc)
	return op
}

func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleConnections)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
