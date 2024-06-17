package tests

import (
	"CollabDoc/pkg/document"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConflictResolution(t *testing.T) {
	cr := document.ConflictResolver{}
	doc := &document.Document{
		ID:      "doc1",
		Content: map[string]string{"title": "Initial Title"},
	}

	// Apply a conflicting update
	cr.ResolveConflict(doc, "title", "Resolved Title")
	assert.Equal(t, "Resolved Title", doc.Content["title"])
}
