package document

import (
	"sync"
)

type ConflictResolver struct {
	mutex sync.Mutex
}

func (cr *ConflictResolver) ResolveConflict(doc *Document, op Operation) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	// Apply the operation directly to the document content
	err := ApplyOperation(doc, op)
	if err != nil {
		// Handle the error appropriately (e.g., logging, returning error, etc.)
		return
	}
}
