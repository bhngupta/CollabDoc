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
	clients           = make(map[string]map[*websocket.Conn]bool)
	clientMutex       sync.Mutex
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
		// Ensure documents map is initialized even when loading from persistence
		if ss.Documents == nil {
			ss.Documents = make(map[string]*document.Document)
		}
	}

	clients = make(map[string]map[*websocket.Conn]bool)
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	docID := r.URL.Query().Get("docID")
	if docID == "" {
		http.Error(w, "docID is required", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer func() {
		ws.Close()
		clientMutex.Lock()
		delete(clients[docID], ws)
		if len(clients[docID]) == 0 {
			delete(clients, docID)
		}
		clientMutex.Unlock()
		fmt.Println("Client disconnected:", ws.RemoteAddr())
	}()

	clientMutex.Lock()
	if clients[docID] == nil {
		clients[docID] = make(map[*websocket.Conn]bool)
		fmt.Printf("Created new client map for document ID: %s\n", docID)
	}
	clients[docID][ws] = true
	clientMutex.Unlock()

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

			clientMutex.Lock()
			for client := range clients[msg.Op.DocID] {
				if client != ws {
					if err := client.WriteJSON(doc); err != nil {
						fmt.Println("Broadcast error:", err)
						client.Close()
						delete(clients[msg.Op.DocID], client)
					}
				}
			}
			clientMutex.Unlock()

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
