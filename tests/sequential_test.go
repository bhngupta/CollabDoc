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

func TestSequentialUpdates(t *testing.T) {
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

	// Connect two WebSocket clients
	client1, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer client1.Close()

	client2, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.NoError(t, err)
	defer client2.Close()

	// Create a document using Client 1
	createMsg := server.Message{Type: "create", Op: document.Operation{DocID: "doc1"}}
	err = client1.WriteJSON(createMsg)
	assert.NoError(t, err)

	var createResponse document.Document
	err = client1.ReadJSON(&createResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", createResponse.ID)
	t.Logf("Created document: %+v\n", createResponse)

	// Client 1 sends an initial update
	updateMsgClient1 := server.Message{Type: "operation", Op: document.Operation{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Title from Client 1", BaseVersion: createResponse.Version}}
	err = client1.WriteJSON(updateMsgClient1)
	assert.NoError(t, err)

	var updateResponse1 map[string]interface{}
	err = client1.ReadJSON(&updateResponse1)
	assert.NoError(t, err)
	assert.True(t, updateResponse1["success"].(bool))
	t.Logf("Client 1 update response: %+v\n", updateResponse1)

	// Fetch the latest version after Client 1's update
	getMsg := server.Message{Type: "get", Op: document.Operation{DocID: "doc1"}}
	err = client1.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponseClient1 document.Document
	err = client1.ReadJSON(&getResponseClient1)
	assert.NoError(t, err)
	t.Logf("Client 1 get document response: %+v\n", getResponseClient1)

	// Ensure Client 2 has the updated version before sending its operation
	err = client2.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponseClient2 document.Document
	err = client2.ReadJSON(&getResponseClient2)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponseClient2.ID)
	t.Logf("Client 2 get document response: %+v\n", getResponseClient2)

	// Ensure Client 2 has the updated version before sending its operation
	assert.Equal(t, getResponseClient1.Version, getResponseClient2.Version, "Client 2 should have the updated document version before its operation")

	// Client 2 sends an update at position 5 after Client 1's update
	updateMsgClient2 := server.Message{Type: "operation", Op: document.Operation{DocID: "doc1", OpType: "insert", Pos: 5, Content: "Title from Client 2", BaseVersion: getResponseClient2.Version}}
	err = client2.WriteJSON(updateMsgClient2)
	assert.NoError(t, err)

	var updateResponse2 map[string]interface{}
	err = client2.ReadJSON(&updateResponse2)
	assert.NoError(t, err)
	assert.True(t, updateResponse2["success"].(bool))
	t.Logf("Client 2 update response: %+v\n", updateResponse2)

	// Retrieve the document using Client 1 to verify final content
	err = client1.WriteJSON(getMsg)
	assert.NoError(t, err)

	var finalGetResponse document.Document
	err = client1.ReadJSON(&finalGetResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", finalGetResponse.ID)
	t.Logf("Final document state: %+v\n", finalGetResponse)

	// Verify the final content of the document
	finalContent := finalGetResponse.Content
	expectedContent := "TitleTitle from Client 2 from Client 1"
	assert.Equal(t, expectedContent, finalContent, "Final content is unexpected: %s", finalContent)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", finalContent)
}
