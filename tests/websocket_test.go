package tests

import (
	"CollabDoc/internal/server"
	"CollabDoc/pkg/document"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketServer(t *testing.T) {
	// Start the WebSocket server
	go server.StartWebSocketServer()

	// Create a WebSocket client
	u := "ws://localhost:8080/ws"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Test create document
	createMsg := server.Message{Type: "create", DocID: "doc1"}
	err = ws.WriteJSON(createMsg)
	assert.NoError(t, err)

	var createResponse document.Document
	err = ws.ReadJSON(&createResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", createResponse.ID)

	// Test update document
	updateMsg := server.Message{Type: "update", DocID: "doc1", Key: "title", Value: "Collaborative Document"}
	err = ws.WriteJSON(updateMsg)
	assert.NoError(t, err)

	var updateResponse map[string]bool
	err = ws.ReadJSON(&updateResponse)
	assert.NoError(t, err)
	assert.True(t, updateResponse["success"])

	// Test get document
	getMsg := server.Message{Type: "get", DocID: "doc1"}
	err = ws.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponse document.Document
	err = ws.ReadJSON(&getResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponse.ID)
	assert.Equal(t, "Collaborative Document", getResponse.Content["title"])
}
