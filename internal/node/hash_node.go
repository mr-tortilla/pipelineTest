package node

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

// var _ pipeline.Node = (*HashNode)(nil)

// Result содержит результат вычисления MD5 хеша файла
type Result struct {
	Path string
	Hash []byte
}

// HashNode вычисляет MD5 хеш для каждого файла из входного канала
type HashNode struct {
	inputs  []chan any
	outputs []chan any
}

func (n *HashNode) SetInputs(inputs []chan any)   { n.inputs = inputs }
func (n *HashNode) SetOutputs(outputs []chan any) { n.outputs = outputs }

// Run читает пути из входного канала и вычисляет MD5 хеш для каждого файла
func (n *HashNode) Run(ctx context.Context) {
	in := n.inputs[0]
	out := n.outputs[0]
	err := n.outputs[1]

	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-in:
			if !ok {
				return
			}
			path := val.(string)
			hash, e := hashFile(path)
			if e != nil {
				select {
				case <-ctx.Done():
					return
				case err <- fmt.Errorf("%s: %w", path, e):
				}
				continue
			}
			select {
			case <-ctx.Done():
				return
			case out <- Result{Path: path, Hash: hash}:
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
