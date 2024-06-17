package document

func TransformOperation(op, prevOp Operation) Operation {
	if op.Pos >= prevOp.Pos && prevOp.OpType == "insert" {
		op.Pos += len(prevOp.Content)
	} else if op.Pos >= prevOp.Pos && prevOp.OpType == "delete" {
		op.Pos -= prevOp.Length
	}
	return op
}
