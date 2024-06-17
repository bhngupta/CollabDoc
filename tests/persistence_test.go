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
	ss.UpdateDocument("doc1", "title", "Collaborative Document")
	err := persist.SaveState(ss)
	assert.NoError(t, err)

	// Load state and verify
	loadedSS, err := persist.LoadState()
	assert.NoError(t, err)

	doc, exists := loadedSS.GetDocument("doc1")
	assert.True(t, exists)
	assert.Equal(t, "Collaborative Document", doc.Content["title"])
}
