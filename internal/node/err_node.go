package node

import (
	"context"
	"fmt"
	"os"

	"github.com/mr-tortilla/pipelineLibrary"
)

var _ pipeline.Node = (*ErrNode)(nil)

// ErrNode читает ошибки из входного канала и печатает их в stderr
type ErrNode struct {
	In <-chan error // входной канал с ошибками
}

// Run читает ошибки из входного канала и печатает их в stderr
func (n *ErrNode) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-n.In:
			if !ok {
				return
			}
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
		}
	}
}
