package document

type Document struct {
	ID      string            `json:"id"`
	Content map[string]string `json:"content"`
}

type StateSynchronizer struct {
	Documents        map[string]*Document
	ConflictResolver *ConflictResolver
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
	doc := &Document{ID: id, Content: make(map[string]string)}
	ss.Documents[id] = doc
	return doc
}

func (ss *StateSynchronizer) UpdateDocument(id, key, value string) bool {
	doc, exists := ss.Documents[id]
	if !exists {
		return false
	}
	ss.ConflictResolver.ResolveConflict(doc, key, value)
	return true
}
