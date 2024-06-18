package tests

import (
	"CollabDoc/pkg/document"
	"testing"

	"github.com/stretchr/testify/assert"
	"sync"
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
		Length:  len(doc.Content),
	}

	// Prepare the second conflicting update operation
	op2 := document.Operation{
		DocID:   "doc1",
		OpType:  "update",
		Pos:     0,
		Content: "Title from Client 2",
		Length:  len("Title from Client 1"),
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
	expectedContent1 := "Title from Client 1"
	expectedContent2 := "Title from Client 2"
	assert.True(t, doc.Content == expectedContent1 || doc.Content == expectedContent2, "Final content is unexpected: %s", doc.Content)

	// Log the final content for visibility
	t.Logf("Final content of the document: %s", doc.Content)
}
