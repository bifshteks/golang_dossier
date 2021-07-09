package essence

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	g "golang_dossier/graph"
)

type MarkDict map[int][]*Mark

func GetMarkDict(graph *g.Graph, db *sqlx.DB) MarkDict {
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
	ID          int           `json:"id"`
	CreatorId   sql.NullInt64 `json:"creator_id" db:"creator_id"`
	EssenceId   int           `json:"essence_id" db:"essence_id"`
	CreatedDate string        `json:"created_date" db:"created_date"`
	Color       string        `json:"color"`
	Comment     string        `json:"comment"`
}

func getMarksFromDB(graph *g.Graph, db *sqlx.DB) ([]*Mark, error) {
	var marksList []*Mark

	q, args, err := sqlx.In(`
		SELECT
			*
		from graph_essencemark
		where essence_id IN (?)
	`, graph.NodesIdsList())
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
