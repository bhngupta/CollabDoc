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
	ApplyOperation(doc, op)
}
