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
		Content: "Initial Title",
	}

	// Apply a conflicting update
	op := document.Operation{
		DocID:   "doc1",
		OpType:  "update",
		Pos:     0,
		Content: "Resolved Title",
	}
	cr.ResolveConflict(doc, op)
	assert.Equal(t, "Resolved Title", doc.Content)
}
