package graph

type EdgeDict map[int]map[int]EdgeData
type Graph struct {
	//nodes_list []*GraphNode
	RootId int
	Nodes  map[int]*Node
	Edges  EdgeDict // edge[1<parent>][2<child>] = EdgeData{}
}
type Node struct {
	ID           int   // `json:"id"`
	Successors   []int //`json:"successors"`
	Predecessors []int // `json:"predecessors"`
	//Data         *ge.Essence `json:"data"`
}
type EdgeData map[string]string
type Edge struct {
	ParentId int
	ChildId  int
	//Data     EdgeData
}

func (g *Graph) AddNode(id int) {
	g.Nodes[id] = &Node{
		ID:           id,
		Successors:   make([]int, 0), // https://www.reddit.com/r/golang/comments/4rinrk/how_do_you_create_an_array_of_variable_length/
		Predecessors: make([]int, 0),
	}
}

func (g *Graph) AddEdge(edge Edge) {
	if _, ok := g.Nodes[edge.ParentId]; !ok {
		g.AddNode(edge.ParentId)
	}
	if _, ok := g.Nodes[edge.ChildId]; !ok {
		g.AddNode(edge.ChildId)
	}
	g.Nodes[edge.ParentId].Successors = append(g.Nodes[edge.ParentId].Successors, edge.ChildId)
	g.Nodes[edge.ChildId].Predecessors = append(g.Nodes[edge.ChildId].Predecessors, edge.ParentId)

	//childDict, ok := g.Edges[edge.ParentId]
	//if !ok {
	//	g.Edges[edge.ParentId] = map[int]EdgeData{}
	//	childDict = g.Edges[edge.ParentId]
	//}
	//childDict[edge.ChildId] = edge.Data  // todo
}

func (g *Graph) AddEdges(edges []Edge) {
	for _, edge := range edges {
		g.AddEdge(edge)
	}
}

func (g *Graph) nodesList() []*Node {
	nodes := make([]*Node, 0)
	for _, v := range g.Nodes {
		nodes = append(nodes, v)
	}
	return nodes
}

func (g *Graph) NodesIdsList() []int {
	nodes := make([]int, 0)
	for _, v := range g.Nodes {
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

func NewGraph(rootId int) *Graph {
	return &Graph{
		RootId: rootId,
		Nodes:  make(map[int]*Node),
		Edges:  make(EdgeDict),
	}
}
