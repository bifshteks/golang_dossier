package essence

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	g "golang_dossier/graph"
	"log"
)

type Essence struct {
	ID                  int           `json:"id"`
	HorizontalID        uuid.UUID     `json:"horizontal_id" db:"horizontal_id"`
	Version             int           `json:"version"`
	CreatedDate         string        `json:"created_date" db:"created_date"`
	UpdatedDate         string        `json:"updated_date" db:"updated_date"`
	Name                string        `json:"name"`
	Slug                string        `json:"slug"`
	Value               string        `json:"value"`
	SchemaId            int           `json:"schema_id" db:"schema_id"`
	//RemovedFromGraph    bool          `json:"removed_from_graph" db:"removed_from_graph"`
	SchemaDefinitionIds pq.Int64Array `json:"schema_definition_ids" db:"schema_definition_ids"`
	DefinitionIds       pq.Int64Array `json:"definition_ids" db:"definition_ids"`
	//ParentIds           pq.Int64Array `json:"parent_ids" db:"parent_ids"`
	//ChildrenIds         pq.Int64Array `json:"children_ids" db:"children_ids"`
}

func getStepsDown(db *sqlx.DB, rootId int) ([]g.Edge, error) {
	const modelTable = "graph_essence"
	const essenceThroughTable = "graph_essencethrough"
	const SQL_ = `
        -- В recursive_buffer хранятся все м2м элементы-связки, элементы которых есть в указанном графе
        WITH RECURSIVE recursive_buffer AS (
            -- Первый селект - м2м элементы-связки, где ребенок - заданный корень
            SELECT
                id, parent_id, child_id
            FROM
                %[2]s -- through  --
            WHERE
                parent_id = %[3]d -- element_id  --

            -- используем UNION, а не UNION ALL, т.к. если мы получаем данные по графу, а не дереву,
            -- то может добавить несколько одинаковых записей с буфер (пройтись дважды по одному пути).
            UNION

            SELECT
                -- L.id нужен, потому-что для сырого запроса джанга всегда требует чтобы быть
                -- уникальный ключ для каждой строки
                --(django.db.models.query_utils.InvalidQuery: Raw query must include the primary key)
                L.id, L.parent_id, L.child_id
            FROM
                recursive_buffer AS P
                INNER JOIN %[2]s AS L -- through  --
                    ON P.child_id = L.parent_id
        )
        select
           ASD.parent_id, ASD.child_id
        from recursive_buffer as ASD
        --DONT NEED IT 
		--inner join %[1]s as M
        --    on child_id = m.id
        --order by M.name, M.id
	`
	var SQL = fmt.Sprintf(SQL_, modelTable, essenceThroughTable, rootId)

	rows, err := db.Query(SQL)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			fmt.Println("Could not close steps rows")
		}
	}()
	var edges []g.Edge
	for rows.Next() {
		var edge g.Edge
		err = rows.Scan(&edge.ParentId, &edge.ChildId)
		if err != nil {
			return nil, err
		}
		edges = append(edges, edge)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return edges, nil
}

func bulkUpdateGraphData(ids []int, graph *g.Graph, db *sqlx.DB) error {
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
		       --"graph_essence"."removed_from_graph", 
		       ARRAY_AGG("graph_schemadefinition"."defining_id" ) filter (where "graph_schemadefinition"."defining_id"  is not null) AS "schema_definition_ids",
		       ARRAY_AGG("graph_essencedefinition"."defining_id" ) filter (where "graph_essencedefinition"."defining_id"  is not null) AS "definition_ids",
		       --ARRAY_AGG("graph_essencethrough"."child_id" ) filter (where "graph_essencethrough"."child_id"  is not null) AS "children_ids", 
		       --ARRAY_AGG(T9."parent_id" ) filter (where T9."parent_id"  is not null) AS "parent_ids" 
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

    rows, err := db.Queryx(db.Rebind(q), args...)
    if err != nil {
    	return err
	}
    for rows.Next() {
    	var essence Essence
        err := rows.StructScan(&essence)
        if err != nil {
            log.Fatalln(err)
        }
        essences = append(essences, &essence)
        if essence.SchemaDefinitionIds == nil {
        	essence.SchemaDefinitionIds = []int64{}
		}
		if essence.DefinitionIds == nil {
			essence.DefinitionIds = []int64{}
		}
		//if essence.ParentIds == nil {
		//	essence.ParentIds = []int64{}
		//}
		//if essence.ChildrenIds == nil {
		//	essence.ChildrenIds = []int64{}
		//}
    }
	fillGraphData(essences, graph)
	return nil
}

func fillGraphData(essences []*Essence, graph *g.Graph) {
	for _, essence := range essences {
		graph.Nodes[essence.ID].Data = essence
	}
}

func ProcessSteps(steps []g.Edge, rootId int, db *sqlx.DB) (*g.Graph, error) {

	graph := g.NewGraph(rootId)
	graph.AddEdges(steps)

	//fmt.Println("got node", len(graph.nodesList()), graph.nodes[198188])

	var essenceIdToBulkGetData []int
	i := 0
	for _, nodeId := range graph.NodesIdsList() {
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

	fmt.Println("end data for ess") //, graph.nodes[198188].data)

	return graph, nil
}

func GetEssenceGraph(db *sqlx.DB, rootId int) *g.Graph {
	steps, err := getStepsDown(db, rootId)
	if err != nil {
		fmt.Println("error on get steps down", err)
		panic(err)
	}
	graph, err := ProcessSteps(steps, rootId, db)
	if err != nil {
		fmt.Println("error from procesSteps")
		panic(err)
		//return nil, err
	}
	return graph
}
