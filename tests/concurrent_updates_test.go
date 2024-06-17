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
	createMsg := server.Message{Type: "create", DocID: "doc1"}
	err = client1.WriteJSON(createMsg)
	assert.NoError(t, err)

	var createResponse document.Document
	err = client1.ReadJSON(&createResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", createResponse.ID)

	// Prepare update messages from both clients
	updateMsgClient1 := server.Message{Type: "update", DocID: "doc1", Key: "title", Value: "Title from Client 1"}
	updateMsgClient2 := server.Message{Type: "update", DocID: "doc1", Key: "title", Value: "Title from Client 2"}

	var wg sync.WaitGroup
	wg.Add(2)

	// Send update from Client 1
	go func() {
		defer wg.Done()
		err = client1.WriteJSON(updateMsgClient1)
		assert.NoError(t, err)

		var updateResponse map[string]bool
		err = client1.ReadJSON(&updateResponse)
		assert.NoError(t, err)
		assert.True(t, updateResponse["success"])
	}()

	// Send update from Client 2
	go func() {
		defer wg.Done()
		err = client2.WriteJSON(updateMsgClient2)
		assert.NoError(t, err)

		var updateResponse map[string]bool
		err = client2.ReadJSON(&updateResponse)
		assert.NoError(t, err)
		assert.True(t, updateResponse["success"])
	}()

	// Wait for both updates to complete
	wg.Wait()

	// Retrieve the document using Client 1
	getMsg := server.Message{Type: "get", DocID: "doc1"}
	err = client1.WriteJSON(getMsg)
	assert.NoError(t, err)

	var getResponse document.Document
	err = client1.ReadJSON(&getResponse)
	assert.NoError(t, err)
	assert.Equal(t, "doc1", getResponse.ID)

	// Verify the final title of the document
	finalTitle := getResponse.Content["title"]
	assert.Contains(t, []string{"Title from Client 1", "Title from Client 2"}, finalTitle)

	// Log the final title for visibility
	t.Logf("Final title of the document: %s", finalTitle)
}
