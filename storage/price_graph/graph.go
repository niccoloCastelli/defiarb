package price_graph

import (
	"fmt"
	"github.com/pkg/errors"
)

func NewNode(address string, name string) *Node {
	return &Node{Address: address, Name: name}
}

type Node struct {
	Address   string
	Name      string
	edgesFrom []*Edge
	edgesTo   []*Edge
}

type Edge struct {
	Address  string
	Name     string
	Exchange string
	From     *Node
	To       *Node
	Weight   float64
}

func NewGraph() *Graph {
	return &Graph{
		nodes: map[string]*Node{},
		edges: map[string]*Edge{},
	}
}

type Graph struct {
	nodes map[string]*Node
	edges map[string]*Edge
}

func (g *Graph) SetWeight(edgeKey string, weight float64) error {
	edge, ok := g.edges[edgeKey]
	if !ok {
		return errors.New("edge not found")
	}
	edge.Weight = weight
	return nil
}

func (g *Graph) AddNode(node Node) {
	if _, ok := g.nodes[node.Address]; !ok {
		node.edgesFrom = []*Edge{}
		node.edgesTo = []*Edge{}
		g.nodes[node.Address] = &node
	}
}

func (g *Graph) GetNode(address string) (*Node, error) {
	if n, ok := g.nodes[address]; ok {
		return n, nil
	}
	return nil, errors.New("node not found")
}

func (g *Graph) AddEdge(fromAddr string, toAddr string, weight float64, address string, name string, exchange string) error {
	nodeFrom, ok := g.nodes[fromAddr]
	if !ok {
		return errors.New("node from not found")
	}
	nodeTo, ok := g.nodes[toAddr]
	if !ok {
		return errors.New("node to not found")
	}
	edge := &Edge{
		Address:  address,
		Name:     name,
		From:     nodeFrom,
		To:       nodeTo,
		Weight:   weight,
		Exchange: exchange,
	}
	edgeKey := fmt.Sprintf("%s@%s->%s", address, fromAddr, toAddr)
	g.edges[edgeKey] = edge
	nodeFrom.edgesFrom = append(nodeFrom.edgesFrom, edge)
	nodeTo.edgesTo = append(nodeFrom.edgesTo, edge)
	return nil
}

func EdgeKey(address string, fromAddr string, toAddr string) string {
	return fmt.Sprintf("%s@%s->%s", address, fromAddr, toAddr)
}
