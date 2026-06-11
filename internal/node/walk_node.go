package node

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// var _ pipeline.Node = (*WalkNode)(nil)

// WalkNode рекурсивно обходит директорию и пишет пути файлов в выходной канал
type WalkNode struct {
	Dir     string
	inputs  []chan any
	outputs []chan any
}

func (n *WalkNode) SetInputs(inputs []chan any)   { n.inputs = inputs }
func (n *WalkNode) SetOutputs(outputs []chan any) { n.outputs = outputs }

// Run запускает рекурсивный обход директории
// По завершении закрывает все выходные каналы
func (n *WalkNode) Run(ctx context.Context) {
	filepath.WalkDir(n.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERR: %v\n", err)
			return nil
		}
		if d.IsDir() {
			return nil
		}

		for _, out := range n.outputs {
			select {
			case <-ctx.Done():
				return filepath.SkipAll
			case out <- path:
			}
		}

		return nil
	})
}
