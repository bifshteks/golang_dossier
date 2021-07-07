package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

const EssenceRootId = 649 //176945 //649

func doJob() (*Graph, MarkDict, EdgeDict) {
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
	markDict := GetMarkDict(graph, db) //
	linksDict := graph.edges            //
	return graph, markDict, linksDict
	//fmt.Println("marksDict", markDict, "\n\n")
	//fmt.Println("linksDict", linksDict, "\n\n")
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		graph, _, _ := doJob()
		c.JSON(200, gin.H{
			//"message": marks,
			//"qwe": links,
			"graph": JsonizeGraph(graph.nodes[graph.rootId], graph),
		})
	})
	if err := r.Run(); err != nil {
		fmt.Println("Could not run server")
	}
}