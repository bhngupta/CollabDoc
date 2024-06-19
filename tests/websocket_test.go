package tests

import (
	"CollabDoc/internal/server"
	"CollabDoc/pkg/document"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketServer(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(server.HandleConnections))
	defer testServer.Close()

	// Convert the test server URL to WebSocket URL
	u := "ws" + testServer.URL[4:] + "/ws"

	// Create a WebSocket client
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer ws.Close()

	// Test create document
	createMsg := server.Message{Type: "create", Op: document.Operation{DocID: "doc1"}}
	err = ws.WriteJSON(createMsg)
	assert.NoError(t, err)

	var createResponse document.Document
	err = ws.ReadJSON(&createResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", createResponse.ID)

	// Test update document
	updateMsg := server.Message{Type: "operation", Op: document.Operation{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Collaborative Document", BaseVersion: createResponse.Version}}
	err = ws.WriteJSON(updateMsg)
	assert.NoError(t, err)

	var updateResponse map[string]interface{}
	err = ws.ReadJSON(&updateResponse)
	assert.NoError(t, err)
	assert.True(t, updateResponse["success"].(bool))

	// Test get document
	getMsg := server.Message{Type: "get", Op: document.Operation{DocID: "doc1"}}
	err = ws.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponse document.Document
	err = ws.ReadJSON(&getResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponse.ID)
	assert.Equal(t, "Collaborative Document", getResponse.Content)

	// Close the WebSocket connection gracefully
	err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	assert.NoError(t, err)
}
