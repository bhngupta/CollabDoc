package tests

import (
	"CollabDoc/pkg/document"
	"CollabDoc/pkg/persistence"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPersistence(t *testing.T) {
	filePath := "test_state.json"
	defer os.Remove(filePath)

	// Initialize StateSynchronizer and Persistence
	ss := document.NewStateSynchronizer()
	persist := persistence.NewPersistence(filePath)

	// Create a document and save state
	ss.CreateDocument("doc1")
	op := document.Operation{
		DocID:   "doc1",
		OpType:  "update",
		Pos:     0,
		Content: "Collaborative Document",
	}
	ss.UpdateDocument("doc1", op)
	err := persist.SaveState(ss)
	assert.NoError(t, err)

	// Load state and verify
	loadedSS, err := persist.LoadState()
	assert.NoError(t, err)

	doc, exists := loadedSS.GetDocument("doc1")
	assert.True(t, exists)
	assert.Equal(t, "Collaborative Document", doc.Content)

	// Create another document, update it, and save state
	ss.CreateDocument("doc2")
	op2 := document.Operation{
		DocID:   "doc2",
		OpType:  "insert",
		Pos:     0,
		Content: "New Document",
	}
	ss.UpdateDocument("doc2", op2)
	err = persist.SaveState(ss)
	assert.NoError(t, err)

	// Load state again and verify both documents
	loadedSS, err = persist.LoadState()
	assert.NoError(t, err)

	doc1, exists := loadedSS.GetDocument("doc1")
	assert.True(t, exists)
	assert.Equal(t, "Collaborative Document", doc1.Content)

	doc2, exists := loadedSS.GetDocument("doc2")
	assert.True(t, exists)
	assert.Equal(t, "New Document", doc2.Content)
}
