// internal/server/websocket.go

package server

import (
	"CollabDoc/pkg/document"
	"CollabDoc/pkg/persistence"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Type  string `json:"type"`
	DocID string `json:"doc_id,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

var (
	ss      *document.StateSynchronizer
	persist *persistence.Persistence
)

func init() {
	persist = persistence.NewPersistence("state.json")
	var err error
	ss, err = persist.LoadState()
	if err != nil {
		fmt.Println("Error loading state:", err)
		ss = document.NewStateSynchronizer()
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
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
		case "create":
			doc := ss.CreateDocument(msg.DocID)
			err = ws.WriteJSON(doc)
		case "update":
			success := ss.UpdateDocument(msg.DocID, msg.Key, msg.Value)
			if success {
				persist.SaveState(ss)
			}
			err = ws.WriteJSON(map[string]bool{"success": success})
		case "get":
			doc, exists := ss.GetDocument(msg.DocID)
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

func StartWebSocketServer() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
