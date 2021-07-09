package main

import (
	g "golang_dossier/graph"
	ge "golang_dossier/graph/essence"
)
type JsonGraph struct {
	ge.Essence
	Children []*JsonGraph `json:"children"`
}

func JsonizeGraph(node *g.Node, graph *g.Graph) *JsonGraph {
	var jsonized JsonGraph
	jsonized.ID = node.Data.ID
	jsonized.HorizontalID = node.Data.HorizontalID
	jsonized.Version = node.Data.Version
	jsonized.CreatedDate = node.Data.CreatedDate
	jsonized.UpdatedDate = node.Data.UpdatedDate
	jsonized.Name = node.Data.Name
	jsonized.Slug = node.Data.Slug
	jsonized.Value = node.Data.Value
	jsonized.SchemaId = node.Data.SchemaId
	//jsonized.RemovedFromGraph = node.Data.RemovedFromGraph
	jsonized.SchemaDefinitionIds = node.Data.SchemaDefinitionIds
	jsonized.DefinitionIds = node.Data.DefinitionIds
	//jsonized.ParentIds = node.Data.ParentIds
	//jsonized.ChildrenIds = node.Data.ChildrenIds
	jsonized.Children = make([]*JsonGraph, 0)

	for _, childId := range node.Successors {
		jsonized.Children = append(jsonized.Children, JsonizeGraph(graph.Nodes[childId], graph))
	}
	return &jsonized
}
