package document

import "time"

func ApplyOperation(doc *Document, op Operation) {
	switch op.OpType {
	case "insert":
		doc.Content = doc.Content[:op.Pos] + op.Content + doc.Content[op.Pos:]
	case "delete":
		doc.Content = doc.Content[:op.Pos] + doc.Content[op.Pos+op.Length:]
	case "update":
		doc.Content = op.Content
	}
	doc.Version++
	doc.UpdatedAt = time.Now().Format(time.RFC3339)
}
