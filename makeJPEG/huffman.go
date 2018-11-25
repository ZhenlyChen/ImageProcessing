package makeJPEG

import (
	"sort"
)

type Node struct {
	Value  int
	Weight int
	Left   *Node
	Right  *Node
}

type Nodes []Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Weight < n[j].Weight
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

var Root *Node

func huffmanEncode(priorityMap map[int]int) map[int]string {
	nodes := makeSortedNodes(priorityMap)
	hfmRoot := makeHuffmanTree(nodes)
	return encodeTree(hfmRoot)
}

func makeSortedNodes(priorityMap map[int]int) []Node {
	var nodeList Nodes
	for k, v := range priorityMap {
		nodeList = append(nodeList, Node{k, v, nil, nil})
	}
	sort.Sort(nodeList)
	return nodeList
}

func makeHuffmanTree(nodes Nodes) *Node {
	if len(nodes) < 2 {
		return &nodes[0]
	}
	for len(nodes) > 1 {
		a := nodes[0]
		b := nodes[1]
		father := Node{-1, a.Weight + b.Weight, &a, &b}
		oldLen := len(nodes)
		for i := 2; i < len(nodes); i++ {
			if father.Weight < nodes[i].Weight {
				nodes := append(nodes, Node{})
				for j := i + 1; j < len(nodes); j++ {
					nodes[j] = nodes[j - 1]
				}
				nodes[i] = father
				break
			}
		}
		if oldLen == len(nodes) {
			nodes = append(nodes, father)
		}
		nodes = nodes[2:]
	}
	return &nodes[0]
}

func encodeTree(root *Node) map[int]string {
	var initialCode string
	encodeMap := make(map[int]string)
	root.traverse(initialCode, func(value int, code string) {
		encodeMap[value] = code
	})
	return encodeMap
}

func (n Node) traverse(code string, visit func(int, string)) {
	if left := n.Left; left != nil {
		left.traverse(code + "0", visit)
	} else {
		visit(n.Value, code)
		return
	}
	n.Right.traverse(code + "1", visit)
}