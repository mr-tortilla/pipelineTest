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

	paths := make(chan any)
	results := make(chan any)
	errs := make(chan any)

	walkNode := &node.WalkNode{Dir: dir}
	printNode := &node.PrintNode{}
	errNode := &node.ErrNode{}

	p := pipeline.NewPipeline()

	// соединяем walkNode со всеми HashNode через один канал
	for i := 0; i < parallelism; i++ {
		hashNode := &node.HashNode{}
		p.Connect(walkNode, hashNode, paths)
		p.Connect(hashNode, printNode, results)
		p.Connect(hashNode, errNode, errs)
		p.Add(hashNode)
	}

	p.Add(walkNode, printNode, errNode)

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
