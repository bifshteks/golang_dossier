package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)
//type MarkAsDict struct {
//	essenceId int
//	marks []Mark
//}

func GetMarkDict(graph *Graph, db *sqlx.DB) map[int][]*Mark {
	markList, err := getMarksFromDB(graph, db)
	if err != nil {
		panic(err)
	}
	markDict := map[int][]*Mark{}
	for _, el := range markList {
		if _, ok := markDict[el.EssenceId]; !ok {
			markDict[el.EssenceId] = make([]*Mark, 0)
		}
		markDict[el.EssenceId] = append(markDict[el.EssenceId], el)
	}
	fmt.Println(len(markDict), markDict)
	return markDict
}

type Mark struct {
	ID int
	CreatorId sql.NullInt64 `db:"creator_id"`
	EssenceId int `db:"essence_id"`
	CreatedDate string `db:"created_date"`
	Color string
	Comment string
}


func getMarksFromDB(graph *Graph, db *sqlx.DB) ([]*Mark, error) {
	var marksList []*Mark

	q, args, err := sqlx.In(`
		SELECT
			*
		from graph_essencemark
		where essence_id IN (?)
	`, graph.nodesIdsList())
	if err != nil {
		fmt.Println("error on in")
		return nil, err
	}
	err = db.Select(&marksList, db.Rebind(q), args...)
	if err != nil {
		fmt.Println("error on mark select")
		return nil, err
	}
	return marksList, nil
}