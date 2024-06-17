// pkg/document/state.go

package document

type Document struct {
	ID      string            `json:"id"`      // This ID can be the same as the URL
	Content map[string]string `json:"content"` // key-value to represent document content
}

type StateSynchronizer struct {
	Documents map[string]*Document
}

func NewStateSynchronizer() *StateSynchronizer {
	return &StateSynchronizer{Documents: make(map[string]*Document)}
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
	doc.Content[key] = value
	return true
}
