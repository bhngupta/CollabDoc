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

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		switch msg.Type {
		case "operation":
			operationsMutex.Lock()
			doc := ss.Documents[msg.Op.DocID]
			opsSinceBase := operations[msg.Op.DocID][msg.Op.BaseVersion:]
			transformedOp := HandleOperation(doc, msg.Op, opsSinceBase)
			operations[msg.Op.DocID] = append(operations[msg.Op.DocID], transformedOp)
			operationsMutex.Unlock()

			// Broadcast the transformed operation
			BroadcastOperation(transformedOp)

			// Save state
			persist.SaveState(ss)
		case "create":
			doc := ss.CreateDocument(msg.Op.DocID)
			err = ws.WriteJSON(doc)
		case "get":
			doc, exists := ss.GetDocument(msg.Op.DocID)
			if exists {
				err = ws.WriteJSON(doc)
			} else {
				err = ws.WriteJSON(map[string]bool{"exists": false})
			}
		default:
			fmt.Println("Unknown message type")
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
