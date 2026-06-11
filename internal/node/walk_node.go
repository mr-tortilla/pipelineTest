package node

import (
	"context"
	"os"
	"path/filepath"

	"github.com/mr-tortilla/pipelineLibrary"
)

// убеждаемся что WalkNode реализует интерфейс Node
var _ pipeline.Node = (*WalkNode)(nil)

// WalkNode рекурсивно обходит директорию и пишет пути файлов в Out.
// Ошибки пишет в ErrOut.
type WalkNode struct {
	Dir    string        // директория для обхода
	Out    chan<- string // выходной канал с путями файлов
	ErrOut chan<- error  // канал ошибок
}

// Run запускает рекурсивный обход директории
// По завершении закрывает выходной канал и канал ошибок
func (n *WalkNode) Run(ctx context.Context) {
	defer close(n.Out)

	err := filepath.WalkDir(n.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			select {
			case <-ctx.Done():
				return filepath.SkipAll
			case n.ErrOut <- err:
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		select {
		case <-ctx.Done():
			return filepath.SkipAll
		case n.Out <- path:
		}

		return nil
	})
	if err != nil {
		n.ErrOut <- err
		return
	}
}
