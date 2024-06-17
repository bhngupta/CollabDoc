package document

import (
	"sync"
)

type ConflictResolver struct {
	mutex sync.Mutex
}

func (cr *ConflictResolver) ResolveConflict(doc *Document, key, newValue string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	// Simple conflict resolution strategy: last write wins
	doc.Content[key] = newValue
}
