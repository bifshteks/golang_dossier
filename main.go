package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func getStepsDown(rootId int, db *sqlx.DB) ([]Edge, error) {
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
	var edges []Edge
	//qwe, err := rows.Columns()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("columns", qwe)
	for rows.Next() {
		var edge Edge
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

func main() {
	db, err := DbConnect()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			fmt.Println("Could not close db connection", err)
		}
	}()
	fmt.Println("before execSql")
	graph := GetEssenceGraph(db)
	_ = GetMarkDict(graph, db) // markDict
	_ = graph.edges            // linksDict
	//fmt.Println("marksDict", markDict, "\n\n")
	//fmt.Println("linksDict", linksDict, "\n\n")
}
