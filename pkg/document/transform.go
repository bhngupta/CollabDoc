package document

import "fmt"

// TransformOperation adjusts the position of the operation based on previous operations
func TransformOperation(op, prevOp Operation) Operation {
	fmt.Printf("Transforming operation: %+v based on previous operation: %+v\n", op, prevOp)

	if prevOp.OpType == "insert" && op.Pos >= prevOp.Pos {
		op.Pos += len(prevOp.Content)
	} else if prevOp.OpType == "delete" && op.Pos >= prevOp.Pos {
		op.Pos -= prevOp.Length
	}

	fmt.Printf("Transformed operation: %+v\n", op)
	return op
}
