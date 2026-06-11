package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/mr-tortilla/pipelineLibrary"
	"github.com/mr-tortilla/pipelineTest/internal/node"
)

const defaultParallelism = 10

func main() {
	dir, parallelism := parseArgs()

	// каналы между нодами
	paths := make(chan string)
	results := make(chan node.Result)
	errs := make(chan error)

	// нода обхода директории
	walk := &node.WalkNode{Dir: dir, Out: paths, ErrOut: errs}

	// группа нод вычисления MD5
	hashNodes := make([]pipeline.Node, parallelism)
	for i := 0; i < parallelism; i++ {
		hashNodes[i] = &node.HashNode{In: paths, Out: results, ErrOut: errs}
	}
	hashGroup := pipeline.NewNodeGroup(
		func() {
			close(results)
			close(errs)
		},
		hashNodes...,
	)

	// нода печати результатов
	print := &node.PrintNode{In: results}

	// нода печати ошибок
	errNode := &node.ErrNode{In: errs}

	// собираем пайплайн
	p := pipeline.New()
	p.Add(walk, hashGroup, print, errNode)

	// обработка прерывания
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	p.Exec(ctx)
	p.Wait()
}

func parseArgs() (string, int) {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: md5walk <directory> [parallelism]\n")
		os.Exit(1)
	}

	dir := os.Args[1]

	parallelism := defaultParallelism
	if len(os.Args) >= 3 {
		n, err := strconv.Atoi(os.Args[2])
		if err != nil || n < 1 {
			fmt.Fprintf(os.Stderr, "parallelism must be a positive number\n")
			os.Exit(1)
		}
		parallelism = n
	}

	return dir, parallelism
}
