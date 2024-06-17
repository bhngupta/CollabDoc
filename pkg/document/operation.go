package document

type Operation struct {
	DocID       string `json:"docID"`
	OpType      string `json:"opType"`  // "insert", "delete", "update"
	Pos         int    `json:"pos"`     // Position in the document content where the operation is applied
	Length      int    `json:"length"`  // Length of text to delete (used for delete operations)
	Content     string `json:"content"` // Content to insert or update
	ClientID    string `json:"clientID"`
	OpID        string `json:"opID"`
	BaseVersion int    `json:"baseVersion"`
}
