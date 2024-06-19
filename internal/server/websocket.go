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
	defer func() {
		ws.Close()
		fmt.Println("Client disconnected:", ws.RemoteAddr())
	}()

	fmt.Println("Client connected:", ws.RemoteAddr())

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
			if doc == nil {
				fmt.Println("Document not found:", msg.Op.DocID)
				operationsMutex.Unlock()
				break
			}

			if _, ok := pendingOperations[msg.Op.DocID]; !ok {
				pendingOperations[msg.Op.DocID] = []document.Operation{}
			}

			pendingOperations[msg.Op.DocID] = append(pendingOperations[msg.Op.DocID], msg.Op)
			operationsMutex.Unlock()

			document.ProcessPendingOperations(doc, msg.Op.DocID, pendingOperations)

			err = persist.SaveState(ss)
			if err != nil {
				fmt.Println("Error saving state:", err)
			}

			err = ws.WriteJSON(map[string]interface{}{"success": true, "version": doc.Version})

		case "create":
			doc := ss.CreateDocument(msg.Op.DocID)
			operationsMutex.Lock()
			pendingOperations[msg.Op.DocID] = []document.Operation{}
			operationsMutex.Unlock()

			fmt.Printf("Created document: %+v\n", doc)
			err = ws.WriteJSON(doc)

		case "get":
			doc, exists := ss.GetDocument(msg.Op.DocID)
			if !exists {
				doc = ss.CreateDocument(msg.Op.DocID)
				operationsMutex.Lock()
				pendingOperations[msg.Op.DocID] = []document.Operation{}
				operationsMutex.Unlock()
				fmt.Printf("Created document as it didn't exist: %+v\n", doc)
			}

			fmt.Printf("Client connected. ClientID: %p, DocID: %s\n", ws, msg.Op.DocID)
			err = ws.WriteJSON(doc)

		case "heartbeat":
			err = ws.WriteJSON(map[string]interface{}{"type": "heartbeat_ack"})
			if err != nil {
				fmt.Println("Error: Sending Heartbeat")
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

func StartWebSocketServer() {
	http.HandleFunc("/ws", HandleConnections)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
