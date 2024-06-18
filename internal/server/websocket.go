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
	ss              *document.StateSynchronizer
	persist         *persistence.Persistence
	operations      = make(map[string][]document.Operation)
	operationsMutex sync.Mutex
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
			BroadcastOperation(transformedOp)

			// Save state
			err = persist.SaveState(ss)
			if err != nil {
				fmt.Println("Error saving state:", err)
			}

		case "create":
			doc := ss.CreateDocument(msg.Op.DocID)
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

func HandleOperation(doc *document.Document, op document.Operation, opsSinceBase []document.Operation) document.Operation {
	for _, prevOp := range opsSinceBase {
		op = document.TransformOperation(op, prevOp)
	}
	document.ApplyOperation(doc, op)
	return op
}

func BroadcastOperation(op document.Operation) {
	// Send the operation to all connected clients
	// Implement this function based on your WebSocket library
}

func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleConnections)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
