package node

import (
	"context"
	"fmt"
)

// var _ pipeline.Node = (*PrintNode)(nil)

// PrintNode печатает результаты вычисления MD5 хешей в stdout
type PrintNode struct {
	inputs  []chan any
	outputs []chan any
}

func (n *PrintNode) SetInputs(inputs []chan any)   { n.inputs = inputs }
func (n *PrintNode) SetOutputs(outputs []chan any) { n.outputs = outputs }

// Run читает результаты из входного канала и печатает их в stdout
func (n *PrintNode) Run(ctx context.Context) {
	in := n.inputs[0]

	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-in:
			if !ok {
				return
			}
			result := val.(Result)
			fmt.Printf("%x  %s\n", result.Hash, result.Path)
		}
	}
}
