package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type nodeType int

const (
	FOLDER nodeType = iota
	FILE
)

type Exporter interface {
	Export() string
}

type Node struct {
	Class  nodeType
	Name   string
	Path   string
	Doc    *Doc
	Parent interface{}
	// Scope   string // for now
}

func NewNode(class nodeType, name string, path string, doc *Doc, parent interface{}) Node {
	return Node{
		Class:  class,
		Name:   name,
		Path:   path,
		Doc:    doc,
		Parent: parent,
	}
}

func (n *Node) Export() string {
	return n.Name
}

type File struct {
	Node
	Contents string
}

func (f File) Export() string {
	return f.Contents + "\n"
}

func NewFile(name string, path string, doc *Doc, parent interface{}) (File, error) {
	contents, err := os.ReadFile(filepath.Join(doc.RootPath, path))
	if err != nil {
		return File{}, err
	}
	return File{Node: NewNode(FILE, name, path, doc, parent), Contents: string(contents)}, nil
}

type Folder struct {
	Node
	Children map[string]interface{}
}

func (f Folder) Export() string {
	s := ""
	for _, n := range f.Children {
		s += n.(Exporter).Export()
	}
	return s
}

func NewFolder(name string, path string, doc *Doc, parent interface{}) *Folder {
	return &Folder{
		Node:     NewNode(FOLDER, name, path, doc, parent),
		Children: make(map[string]interface{}),
	}
}

type Doc struct {
	*Folder
	RootPath string
}

func (doc *Doc) Export() {
	fmt.Println("Export for", doc.RootPath)
	fmt.Println(doc.Folder.Export())
}

func NewDoc(path string) (Doc, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Doc{}, err
	}
	name := "(root)" // for now
	doc := Doc{RootPath: absPath}
	doc.Folder = NewFolder(".", name, &doc, nil)

	return doc, nil
}

func (doc *Doc) FindParentFolder(path string) *Folder {
	pathArray := strings.Split(filepath.Clean(path), string(os.PathSeparator))
	parentArray := pathArray[:len(pathArray)-1] // pop

	currFolder := doc.Folder

	for len(parentArray) > 0 {
		// shift: x, a = a[0], a[1:]
		child, ok := doc.Folder.Children[parentArray[0]]
		if !ok {
			break
		}
		currFolder = child.(*Folder)
		parentArray = parentArray[1:]
	}

	return currFolder
}

func (doc *Doc) addDirectory(path string) {
	name := filepath.Base(path)
	node := doc.FindParentFolder(path)
	node.Children[name] = NewFolder(name, path, doc, node)
}

func (doc *Doc) addFile(path string) error {
	name := filepath.Base(path)
	if name == "_index.bm" {
		// do some clever stuff here :-)
		return nil
	}
	node := doc.FindParentFolder(path)
	file, err := NewFile(name, path, doc, node)
	if err != nil {
		return err
	}
	node.Children[name] = file
	return nil
}

func main() {
	root := "../book"

	doc, err := NewDoc(root)
	if err != nil {
		log.Fatal(err)
	}

	fileSystem := os.DirFS(root)

	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if path == "." {
			return nil
		}
		if d.IsDir() {
			doc.addDirectory(path)
			return nil
		}
		err = doc.addFile(path)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	doc.Export()
}
