package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"strconv"
	"time"
)

//const EssenceRootId = 176945 //649

func doJob(essenceRootId int) (*Graph, MarkDict, EdgeDict) {
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
	graph := GetEssenceGraph(db, essenceRootId)
	markDict := GetMarkDict(graph, db) //
	linksDict := graph.edges            //
	return graph, markDict, linksDict
	//fmt.Println("marksDict", markDict, "\n\n")
	//fmt.Println("linksDict", linksDict, "\n\n")
}

func main() {
	r := gin.Default()
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"PUT", "PATCH", "GET", "OPTIONS"},
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        //AllowOriginFunc: func(origin string) bool {
        //    return origin == "https://github.com"
        //},
        MaxAge: 12 * time.Hour,
    }))
	r.GET("/api/essences/:id/", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		graph, marks, _ := doJob(id)
		c.JSON(200, gin.H{
			"marks": marks,
			//"links": links,
			"tree": JsonizeGraph(graph.nodes[graph.rootId], graph),
		})
	})
	if err := r.Run(); err != nil {
		fmt.Println("Could not run server")
	}
}