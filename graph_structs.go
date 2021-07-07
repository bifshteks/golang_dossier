package main

type EdgeDict map[int]map[int]EdgeData
type Graph struct {
	//nodes_list []*GraphNode
	nodes map[int]*GraphNode
	edges EdgeDict // edge[1<parent>][2<child>] = EdgeData{}
}
type GraphNode struct {
	ID           int
	successors   []int
	predecessors []int
	data         *Essence
}
type EdgeData map[string]string
type Edge struct {
	ParentId int
	ChildId  int
	Data     EdgeData
}

func (g *Graph) AddNode(id int) {
	g.nodes[id] = &GraphNode{
		ID:           id,
		successors:   make([]int, 0), // https://www.reddit.com/r/golang/comments/4rinrk/how_do_you_create_an_array_of_variable_length/
		predecessors: make([]int, 0),
	}
}

func (g *Graph) AddEdge(edge Edge) {
	if _, ok := g.nodes[edge.ParentId]; !ok {
		g.AddNode(edge.ParentId)
	}
	if _, ok := g.nodes[edge.ChildId]; !ok {
		g.AddNode(edge.ChildId)
	}
	g.nodes[edge.ParentId].successors = append(g.nodes[edge.ParentId].successors, edge.ChildId)
	g.nodes[edge.ChildId].predecessors = append(g.nodes[edge.ChildId].predecessors, edge.ParentId)

	childDict, ok := g.edges[edge.ParentId]
	if !ok {
		g.edges[edge.ParentId] = map[int]EdgeData{}
		childDict = g.edges[edge.ParentId]
	}
	childDict[edge.ChildId] = edge.Data
}

func (g *Graph) AddEdges(edges []Edge) {
	for _, edge := range edges {
		g.AddEdge(edge)
	}
}

func (g *Graph) nodesList() []*GraphNode {
	nodes := make([]*GraphNode, 0)
	for _, v := range g.nodes {
		nodes = append(nodes, v)
	}
	return nodes
}

func (g *Graph) nodesIdsList() []int {
	nodes := make([]int, 0)
	for _, v := range g.nodes {
		nodes = append(nodes, v.ID)
	}
	return nodes
}

//func (g *Graph) edgesList() []*Edge {
//	edges := make([]*Edge, 0)
//	for _, edge := range g.edges {
//		edges = append(edges, edge)
//	}
//	return edges
//}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[int]*GraphNode),
		edges: make(EdgeDict),
	}
}
