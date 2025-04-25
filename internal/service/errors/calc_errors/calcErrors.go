package calc_errors

import "fmt"

type SyntaxError struct {
	Message string
} // position?

func (e *SyntaxError) What() string {
	return fmt.Sprintf("Syntax error: %s", e.Message)
}

type EvaluationError struct {
	NodeType string
	Reason   string
	Context  string
}

func (e *EvaluationError) What() string {
	if e.NodeType != "" {
		return fmt.Sprintf("Evaluation error in %s node: %s (%s)", e.NodeType, e.Reason, e.Context)
	}

	return fmt.Sprintf("Evaluation error: %s (%s)", e.Reason, e.Context)
}
