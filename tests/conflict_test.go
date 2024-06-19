package tests

import (
	"CollabDoc/pkg/document"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentConflictResolution(t *testing.T) {
	cr := document.ConflictResolver{}
	doc := &document.Document{
		ID:      "doc1",
		Content: "Initial Title",
	}

	// Prepare the first update operation
	op1 := document.Operation{
		DocID:   "doc1",
		OpType:  "update",
		Pos:     0,
		Content: "Title from Client 1",
		Length:  len("Initial Title"), // Update the length to replace the initial title
	}

	// Prepare the second conflicting update operation
	op2 := document.Operation{
		DocID:   "doc1",
		OpType:  "update",
		Pos:     0,
		Content: "Title from Client 2",
		Length:  len("Initial Title"), // Update the length to replace the initial title
	}

	// Simulate concurrent application of operations
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		cr.ResolveConflict(doc, op1)
	}()

	go func() {
		defer wg.Done()
		cr.ResolveConflict(doc, op2)
	}()

	wg.Wait()

	// Verify the final content of the document

	expectedContent1 := "Title from Client 1ient 2"
	expectedContent2 := "Title from Client 2ient 1"
	assert.True(t, doc.Content == expectedContent1 || doc.Content == expectedContent2, "Final content is unexpected: %s", doc.Content)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", doc.Content)
}
