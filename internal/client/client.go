// internal/client/main.go

package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type  string `json:"type"`
	DocID string `json:"doc_id,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func StartClient() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	fmt.Printf("Connecting to %s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	// Create document
	createMsg := Message{Type: "create", DocID: "doc2"}
	err = c.WriteJSON(createMsg)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// Read response for create
	var createResponse map[string]interface{}
	err = c.ReadJSON(&createResponse)
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	fmt.Printf("Create response: %v\n", createResponse)

	// Update document
	updateMsg := Message{Type: "update", DocID: "doc2", Key: "title", Value: "Collaborative Dosent"}
	err = c.WriteJSON(updateMsg)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// Read response for update
	var updateResponse map[string]interface{}
	err = c.ReadJSON(&updateResponse)
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	fmt.Printf("Update response: %v\n", updateResponse)

	// Get document
	getMsg := Message{Type: "get", DocID: "doc1"}
	err = c.WriteJSON(getMsg)
	if err != nil {
		log.Println("Write error:", err)
		return
	}

	// Read response for get
	for {
		var getResponse map[string]interface{}
		err = c.ReadJSON(&getResponse)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		fmt.Printf("Get response: %v\n", getResponse)
	}
}
