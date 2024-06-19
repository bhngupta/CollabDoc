package tests

import (
	"CollabDoc/internal/server"
	"CollabDoc/pkg/document"
	"CollabDoc/pkg/persistence"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestStreamOfOperations(t *testing.T) {
	// Setup persistence file path
	filePath := "test_state.json"
	defer os.Remove(filePath)
	persist := persistence.NewPersistence(filePath)

	// Initialize StateSynchronizer and save initial state
	ss := document.NewStateSynchronizer()
	err := persist.SaveState(ss)
	assert.NoError(t, err)

	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(server.HandleConnections))
	defer testServer.Close()

	// Convert the test server URL to WebSocket URL
	u := "ws" + testServer.URL[4:] + "/ws"

	// Connect WebSocket client
	client, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer client.Close()

	// Create a document
	createMsg := server.Message{Type: "create", Op: document.Operation{DocID: "doc1"}}
	err = client.WriteJSON(createMsg)
	assert.NoError(t, err)

	var createResponse document.Document
	err = client.ReadJSON(&createResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", createResponse.ID)

	// Define operations with BaseVersion set to 0 initially
	operations := []document.Operation{
		{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Hello World"},
		{DocID: "doc1", OpType: "update", Pos: 6, Length: 5, Content: "Universe"},
		{DocID: "doc1", OpType: "update", Pos: 6, Length: 8, Content: "Everyone"},
		{DocID: "doc1", OpType: "update", Pos: 0, Length: 5, Content: "Hi"},
		{DocID: "doc1", OpType: "update", Pos: 3, Length: 8, Content: "Galaxy"},
		{DocID: "doc1", OpType: "update", Pos: 3, Length: 6, Content: "Beautiful Galaxy"},
		{DocID: "doc1", OpType: "insert", Pos: 19, Content: "!"},
		{DocID: "doc1", OpType: "update", Pos: 0, Length: 2, Content: "Hello"},
		{DocID: "doc1", OpType: "delete", Pos: 6, Length: 10},
		{DocID: "doc1", OpType: "delete", Pos: 0, Length: 5},
		{DocID: "doc1", OpType: "backspace", Pos: 8},
		{DocID: "doc1", OpType: "backspace", Pos: 7},
	}

	var currentVersion int = createResponse.Version // Initial version is set after document creation
	// Apply operations
	for i, op := range operations {
		// Send the operation to the server
		op.BaseVersion = currentVersion
		msg := server.Message{Type: "operation", Op: op}
		err = client.WriteJSON(msg)
		assert.NoError(t, err)

		// Read the server's response
		var response map[string]interface{}
		err = client.ReadJSON(&response)
		assert.NoError(t, err)
		assert.True(t, response["success"].(bool))

		// Update the current version based on the server's response
		currentVersion = int(response["version"].(float64))

		// Log the document state after each operation
		t.Logf("After operation %d (%s): %s", i+1, op.OpType, response)
	}

	// Retrieve the final document state
	getMsg := server.Message{Type: "get", Op: document.Operation{DocID: "doc1"}}
	err = client.WriteJSON(getMsg)
	assert.NoError(t, err)

	var finalDoc document.Document
	err = client.ReadJSON(&finalDoc)
	assert.NoError(t, err)

	// Expected final content: " Galax"
	expectedContent := " Galax"

	// Verify final content
	assert.Equal(t, expectedContent, finalDoc.Content)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", finalDoc.Content)
}
