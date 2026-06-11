package node

import (
	"context"
	"fmt"

	"github.com/mr-tortilla/pipelineLibrary"
)

var _ pipeline.Node = (*PrintNode)(nil)

// PrintNode печатает результаты вычисления MD5 хешей в stdout
type PrintNode struct {
	In <-chan Result // входной канал с результатами
}

// Run читает результаты из входного канала и печатает их в stdout
func (n *PrintNode) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case result, ok := <-n.In:
			if !ok {
				return
			}
			fmt.Printf("%x  %s\n", result.Hash, result.Path)
		}
	}
}
