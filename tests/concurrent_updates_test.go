package tests

import (
	"CollabDoc/internal/server"
	"CollabDoc/pkg/document"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestConcurrentUpdates(t *testing.T) {
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

	// Fetch latest version for Client 2
	getMsg := server.Message{Type: "get", Op: document.Operation{DocID: "doc1"}}
	err = client2.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponse document.Document
	err = client2.ReadJSON(&getResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponse.ID)

	// Prepare update messages from both clients
	updateMsgClient1 := server.Message{Type: "operation", Op: document.Operation{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Title from Client 1", BaseVersion: createResponse.Version}}
	updateMsgClient2 := server.Message{Type: "operation", Op: document.Operation{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Title from Client 2", BaseVersion: getResponse.Version}}

	var wg sync.WaitGroup
	wg.Add(2)

	// Send update from Client 1
	go func() {
		defer wg.Done()
		err = client1.WriteJSON(updateMsgClient1)
		assert.NoError(t, err)

		var updateResponse map[string]interface{}
		err = client1.ReadJSON(&updateResponse)
		assert.NoError(t, err)
		assert.True(t, updateResponse["success"].(bool))
		updateMsgClient2.Op.BaseVersion = int(updateResponse["version"].(float64))
	}()

	// Send update from Client 2
	go func() {
		defer wg.Done()
		err = client2.WriteJSON(updateMsgClient2)
		assert.NoError(t, err)

		var updateResponse map[string]interface{}
		err = client2.ReadJSON(&updateResponse)
		assert.NoError(t, err)
		assert.True(t, updateResponse["success"].(bool))
	}()

	// Wait for both updates to complete
	wg.Wait()

	// Retrieve the document using Client 1
	err = client1.WriteJSON(getMsg)
	assert.NoError(t, err)

	err = client1.ReadJSON(&getResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponse.ID)

	// Verify the final content of the document
	finalContent := getResponse.Content
	expectedContent1 := "Title from Client 2Title from Client 1"
	expectedContent2 := "Title from Client 1Title from Client 2"
	assert.True(t, finalContent == expectedContent1 || finalContent == expectedContent2, "Final content is unexpected: %s", finalContent)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", finalContent)
}
