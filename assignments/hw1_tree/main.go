package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type treeNode struct {
	File     os.FileInfo
	Children []treeNode
}

func (n treeNode) Name() string {
	if n.File.IsDir() {
		return n.File.Name()
	} else {
		return fmt.Sprintf("%s (%s)", n.File.Name(), n.Size())
	}
}

func (n treeNode) Size() string {
	if n.File.Size() > 0 {
		return fmt.Sprintf("%db", n.File.Size())
	} else {
		return "empty"
	}
}

func getNodes(path string, withFiles bool) ([]treeNode, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var nodes []treeNode
	for _, file := range files {
		if !withFiles && !file.IsDir() {
			continue
		}

		node := treeNode{
			File: file,
		}

		if file.IsDir() {
			children, err := getNodes(path+string(os.PathSeparator)+file.Name(), withFiles)
			if err != nil {
				return nil, err
			}

			node.Children = children
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

func printNodes(out io.Writer, nodes []treeNode, parentPrefix string) {
	var (
		lastIdx     = len(nodes) - 1
		prefix      = "├───"
		childPrefix = "│\t"
	)

	for i, node := range nodes {
		if i == lastIdx {
			prefix = "└───"
			childPrefix = "\t"
		}

		fmt.Fprint(out, parentPrefix, prefix, node.Name(), "\n")

		if node.File.IsDir() {
			printNodes(out, node.Children, parentPrefix+childPrefix)
		}
	}
}

func dirTree(out io.Writer, path string, printFiles bool) (err error) {
	nodes, err := getNodes(path, printFiles)
	if err != nil {
		return
	}

	printNodes(out, nodes, "")
	return
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
