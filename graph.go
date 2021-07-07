package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Essence struct {
	ID                  int
	HorizontalID        uuid.UUID `db:"horizontal_id"`
	Version             int
	CreatedDate         string `db:"created_date"`
	UpdatedDate         string `db:"updated_date"`
	Name                string
	Slug                string
	Value               string
	SchemaId            int           `db:"schema_id"`
	RemovedFromGraph    bool          `db:"removed_from_graph"`
	SchemaDefinitionIds pq.Int64Array `db:"schema_definition_ids"`
	DefinitionIds       pq.Int64Array `db:"definition_ids"`
	ParentIds           pq.Int64Array `db:"parent_ids"`
	ChildrenIds         pq.Int64Array `db:"children_ids"`
}

func bulkUpdateGraphData(ids []int, graph *Graph, db *sqlx.DB) error {
	var essences []*Essence
	q, args, err := sqlx.In(`
		SELECT 
		       "graph_essence"."id", 
		       "graph_essence"."horizontal_id", 
		       "graph_essence"."version", 
		       "graph_essence"."created_date", 
		       "graph_essence"."updated_date", 
		       "graph_essence"."name", 
		       "graph_essence"."slug", 
		       "graph_essence"."value", 
		       "graph_essence"."schema_id", 
		       "graph_essence"."removed_from_graph", 
		       ARRAY_AGG("graph_schemadefinition"."defining_id" ) filter (where "graph_schemadefinition"."defining_id"  is not null) AS "schema_definition_ids", 
		       ARRAY_AGG("graph_essencedefinition"."defining_id" ) filter (where "graph_essencedefinition"."defining_id"  is not null) AS "definition_ids",
		       ARRAY_AGG("graph_essencethrough"."child_id" ) filter (where "graph_essencethrough"."child_id"  is not null) AS "children_ids", 
		       ARRAY_AGG(T9."parent_id" ) filter (where T9."parent_id"  is not null) AS "parent_ids" 
		FROM "graph_essence" 
		    INNER JOIN "graph_schema" ON ("graph_essence"."schema_id" = "graph_schema"."id") 
		    LEFT OUTER JOIN "graph_schemadefinition" ON ("graph_schema"."id" = "graph_schemadefinition"."defined_id") 
		    LEFT OUTER JOIN "graph_essencedefinition" ON ("graph_essence"."id" = "graph_essencedefinition"."defined_id") 
		    LEFT OUTER JOIN "graph_essencethrough" ON ("graph_essence"."id" = "graph_essencethrough"."parent_id") 
		    LEFT OUTER JOIN "graph_essencethrough" T9 ON ("graph_essence"."id" = T9."child_id") 
		WHERE "graph_essence"."id" IN (?) GROUP BY "graph_essence"."id"
	`, ids)
	if err != nil {
		return err
	}
	err = db.Select(&essences, db.Rebind(q), args...)
	if err != nil {
		return err
	}
	fillGraphData(essences, graph)
	return nil
}

func fillGraphData(essences []*Essence, graph *Graph) {
	for _, essence := range essences {
		graph.nodes[essence.ID].data = essence
	}
}

func ProcessSteps(steps []Edge, db *sqlx.DB) (*Graph, error) {

	graph := NewGraph()
	graph.AddEdges(steps)

	//fmt.Println("got node", len(graph.nodesList()), graph.nodes[198188])

	var essenceIdToBulkGetData []int
	i := 0
	for _, nodeId := range graph.nodesIdsList() {
		if i > 1000 {
			err := bulkUpdateGraphData(essenceIdToBulkGetData, graph, db)
			if err != nil {
				fmt.Println("Error from bulk", err)
				break
			}
			essenceIdToBulkGetData = []int{}
			i = 0
		}
		i++
		essenceIdToBulkGetData = append(essenceIdToBulkGetData, nodeId)
	}

	// if something in essenceIdToBulkGetData left not updated
	err := bulkUpdateGraphData(essenceIdToBulkGetData, graph, db)
	if err != nil {
		fmt.Println("Error from bulk2", err)
		return nil, err
	}

	fmt.Println("end data for ess")//, graph.nodes[198188].data)

	return graph, nil
}

func GetEssenceGraph(db *sqlx.DB) (*Graph) {
	const essenceRootId = 649
	steps, err := getStepsDown(essenceRootId, db)
	if err != nil {
		fmt.Println("error on get steps down", err)
		panic(err)
	}
	graph, err := ProcessSteps(steps, db)
	if err != nil {
		fmt.Println("error from procesSteps")
		panic(err)
		//return nil, err
	}
	return graph
}
