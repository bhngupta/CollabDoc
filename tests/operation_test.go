package tests

import (
	"CollabDoc/pkg/document"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperations(t *testing.T) {
	// Initialize StateSynchronizer
	ss := document.NewStateSynchronizer()

	// Create a document
	doc := ss.CreateDocument("doc1")
	t.Logf("Created document: %+v\n", doc)

	// Define operations
	operations := []document.Operation{
		{DocID: "doc1", OpType: "insert", Pos: 0, Content: "Hello", BaseVersion: doc.Version},
		{DocID: "doc1", OpType: "insert", Pos: 5, Content: " World", BaseVersion: doc.Version + 1},
		{DocID: "doc1", OpType: "insert", Pos: 11, Content: "!", BaseVersion: doc.Version + 2},
		{DocID: "doc1", OpType: "update", Pos: 6, Content: "beautiful ", BaseVersion: doc.Version + 3},
		{DocID: "doc1", OpType: "delete", Pos: 0, Length: 5, BaseVersion: doc.Version + 4}, // Deletes "Hello"
	}

	// Apply operations
	for i, op := range operations {
		success := ss.UpdateDocument("doc1", op)
		assert.True(t, success)
		finalDoc, exists := ss.GetDocument("doc1")
		assert.True(t, exists)
		t.Logf("After operation %d (%s): %s", i+1, op.OpType, finalDoc.Content)
	}

	// Retrieve the document
	finalDoc, exists := ss.GetDocument("doc1")
	assert.True(t, exists)

	// Expected final content: " beautiful World!"
	expectedContent := " beautiful World!"

	// Verify final content
	assert.Equal(t, expectedContent, finalDoc.Content)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", finalDoc.Content)
}
