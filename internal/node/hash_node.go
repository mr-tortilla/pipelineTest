package node

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"

	"github.com/mr-tortilla/pipelineLibrary"
)

var _ pipeline.Node = (*HashNode)(nil)

// Result содержит результат вычисления MD5 хеша файла
type Result struct {
	Path string
	Hash []byte
}

// HashNode вычисляет MD5 хеш для каждого файла из входного канала
type HashNode struct {
	In     <-chan string // входной канал с путями файлов
	Out    chan<- Result // выходной канал с результатами
	ErrOut chan<- error  // канал ошибок
}

// Run читает пути из входного канала и вычисляет MD5 хеш для каждого файла
func (n *HashNode) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case path, ok := <-n.In:
			if !ok {
				return
			}
			hash, err := hashFile(path)
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case n.ErrOut <- fmt.Errorf("%s: %w", path, err):
				}
				continue
			}
			select {
			case <-ctx.Done():
				return
			case n.Out <- Result{Path: path, Hash: hash}:
			}
		}
	}
}

// hashFile вычисляет MD5 хеш файла по указанному пути
func hashFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
