package document

import (
	"fmt"
	"time"
)

func ApplyOperation(doc *Document, op Operation) error {
	fmt.Printf("Applying operation: %+v on document: %+v\n", op, doc)
	if op.Pos < 0 || op.Pos > len(doc.Content) {
		fmt.Printf("Invalid operation position: %d, document length: %d\n", op.Pos, len(doc.Content))
		return fmt.Errorf("invalid operation position: %d", op.Pos)
	}

	switch op.OpType {
	case "insert":
		doc.Content = doc.Content[:op.Pos] + op.Content + doc.Content[op.Pos:]
	case "delete":
		if op.Pos+op.Length > len(doc.Content) {
			fmt.Printf("Invalid operation length: %d, document length: %d\n", op.Length, len(doc.Content))
			return fmt.Errorf("invalid operation length: %d", op.Length)
		}
		doc.Content = doc.Content[:op.Pos] + doc.Content[op.Pos+op.Length:]
	case "update":
		if op.Pos+op.Length > len(doc.Content) {
			fmt.Printf("Invalid operation length: %d, document length: %d\n", op.Length, len(doc.Content))
			return fmt.Errorf("invalid operation length: %d", op.Length)
		}
		doc.Content = doc.Content[:op.Pos] + op.Content + doc.Content[op.Pos+op.Length:]
	case "backspace":
		if op.Pos == 0 {
			fmt.Println("Cannot perform backspace at position 0")
			return fmt.Errorf("cannot perform backspace at position 0")
		}
		doc.Content = doc.Content[:op.Pos-1] + doc.Content[op.Pos:]
	default:
		fmt.Printf("Unknown operation type: %s\n", op.OpType)
		return fmt.Errorf("unknown operation type: %s", op.OpType)
	}
	doc.UpdatedAt = time.Now().Format(time.RFC3339)
	fmt.Printf("Document after operation: %+v\n", doc)
	return nil
}
