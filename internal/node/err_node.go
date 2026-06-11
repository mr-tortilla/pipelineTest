package node

import (
	"context"
	"fmt"
	"os"
)

// var _ pipeline.Node = (*ErrNode)(nil)

// ErrNode читает ошибки из входного канала и печатает их в stderr
type ErrNode struct {
	inputs  []chan any
	outputs []chan any
}

func (n *ErrNode) SetInputs(inputs []chan any)   { n.inputs = inputs }
func (n *ErrNode) SetOutputs(outputs []chan any) { n.outputs = outputs }

// Run читает ошибки из входного канала и печатает их в stderr
func (n *ErrNode) Run(ctx context.Context) {
	in := n.inputs[0]

	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-in:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err.(error))
		}
	}
}
