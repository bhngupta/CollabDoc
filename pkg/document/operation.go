package document

import (
	"fmt"
	"sync"
)

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

var operationsMutex sync.Mutex

func ProcessPendingOperations(doc *Document, docID string, pendingOperations map[string][]Operation) {
	operationsMutex.Lock()
	defer operationsMutex.Unlock()

	fmt.Printf("Processing operations for document ID: %s\n", docID)
	if ops, ok := pendingOperations[docID]; ok {
		for _, op := range ops {
			HandleOperation(doc, op)
			fmt.Printf("Document state after applying operation: %+v\n", doc)
		}
		pendingOperations[docID] = []Operation{}
	}
	fmt.Printf("Final document state: %+v\n", doc)
}

func HandleOperation(doc *Document, op Operation) {
	fmt.Printf("Handling operation: %+v\n", op)
	err := ApplyOperation(doc, op)
	if err != nil {
		fmt.Println("ApplyOperation error:", err)
	}
	doc.Version++ // Increment the document version after applying the operation
	fmt.Printf("Document version after operation: %d\n", doc.Version)
	fmt.Printf("Document state after applying operation: %+v\n", doc)
}
