package document

import "time"

type StateSynchronizer struct {
	Documents        map[string]*Document `json:"documents"`
	ConflictResolver *ConflictResolver    `json:"-"`
}

func NewStateSynchronizer() *StateSynchronizer {
	return &StateSynchronizer{
		Documents:        make(map[string]*Document),
		ConflictResolver: &ConflictResolver{}, // Properly initialize ConflictResolver
	}
}

func (ss *StateSynchronizer) GetDocument(id string) (*Document, bool) {
	doc, exists := ss.Documents[id]
	return doc, exists
}

func (ss *StateSynchronizer) CreateDocument(id string) *Document {
	doc := &Document{
		ID:        id,
		Title:     "Untitled Document",
		Content:   "",
		Version:   1,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}
	ss.Documents[id] = doc
	return doc
}

func (ss *StateSynchronizer) UpdateDocument(id string, op Operation) bool {
	doc, exists := ss.Documents[id]
	if !exists {
		return false
	}
	ss.ConflictResolver.ResolveConflict(doc, op)
	return true
}
