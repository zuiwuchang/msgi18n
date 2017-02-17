package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
)

const (
	VERSION = "1.0.0"

	TYPE_NODE  = 1
	TYPE_OTHER = 2
)

type Document struct {
	Sections []*Section
}

func (d *Document) Show() {
	for _, section := range d.Sections {
		section.Show()
	}
}
func (d *Document) HasKey(key string) bool {
	for _, section := range d.Sections {
		if section.HasKey(key) {
			return true
		}
	}
	return false
}

type Section struct {
	Name  string
	Nodes []*Node
	Keys  map[string]*Node
}

func (s *Section) HasKey(key string) bool {
	_, ok := s.Keys[key]
	return ok
}
func (s *Section) Show() {
	fmt.Printf("******	%v	******\n", s.Name)
	fmt.Printf("%v Nodes\n", len(s.Nodes))
	fmt.Printf("%v Keys\n", len(s.Keys))

}
func (s *Section) Push(node *Node) {
	if s.Nodes == nil {
		s.Nodes = make([]*Node, 0, 100)
	}
	s.Nodes = append(s.Nodes, node)

	if node.Type == TYPE_NODE {
		if s.Keys == nil {
			s.Keys = make(map[string]*Node)
		}
		if _, ok := s.Keys[node.Name]; ok {
			fmt.Println("repeat key", node.Name)
		}
		s.Keys[node.Name] = node
	}
}

type Node struct {
	Name  string
	Value string
	Type  int
}

func main() {
	n := len(os.Args)
	if n < 3 {
		fmt.Println("msgi18n", VERSION)
		fmt.Println("writing by king")
		fmt.Println("zuiwuchang@gmail.com")
		return
	}
	docSrc := &Document{}
	fmt.Println("******	read src file	******")
	fmt.Println(os.Args[1])
	if e := readDocument(os.Args[1], docSrc); e != nil {
		fmt.Println(e)
		return
	}
	docSrc.Show()

	for i := 2; i < len(os.Args); i++ {
		fmt.Printf("******	write dist file %v	******\n", i-2)
		dist := os.Args[i]
		fmt.Println(dist)
		if e := writeDocument(dist, docSrc); e != nil {
			fmt.Println(e)
		}
	}
}
func writeDocument(path string, docSrc *Document) error {
	docDist := &Document{}
	if e := readDocument(path, docDist); e != nil {
		return e
	}
	docDist.Show()

	f, e := os.Create(path)
	if e != nil {
		return e
	}
	defer f.Close()

	for _, section := range docDist.Sections {
		if section.Name != "" {
			f.WriteString(section.Name)
			f.WriteString("\n")
		}
		for _, node := range section.Nodes {
			if node.Type == TYPE_NODE {
				f.WriteString(node.Name)
				f.WriteString("=")
				f.WriteString(node.Value)
			} else {
				f.WriteString(node.Name)
			}
			f.WriteString("\n")
		}
	}

	for _, section := range docSrc.Sections {
		for _, node := range section.Nodes {
			if node.Type != TYPE_NODE {
				continue
			}

			if docDist.HasKey(node.Name) {
				continue
			}
			f.WriteString(node.Name)
			f.WriteString("=\n")
		}
	}
	return nil
}
func readDocument(path string, doc *Document) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	doc.Sections = make([]*Section, 0, 10)
	var section *Section
	for {
		b, _, e := reader.ReadLine()
		if e != nil {
			break
		}
		str := strings.TrimSpace(string(b))
		if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
			if section != nil {
				doc.Sections = append(doc.Sections, section)
			}
			section = &Section{Name: str}
			continue
		}

		pos := bytes.IndexByte(b, '=')
		if pos == -1 {
			section.Push(&Node{Name: str, Type: TYPE_OTHER})
			continue
		}

		if section == nil {
			section = &Section{}
		}

		key := b[0:pos]
		if len(key) > 0 {
			section.Push(&Node{Name: string(key), Value: string(b[pos+1:]), Type: TYPE_NODE})
		} else {
			section.Push(&Node{Name: str, Type: TYPE_OTHER})
		}
	}
	if section != nil {
		doc.Sections = append(doc.Sections, section)
	}
	return nil
}
