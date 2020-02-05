package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Note struct {
	Id       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func GetNote(ctx *gin.Context) {
	note_id := ctx.Param("id")

	ctx.JSON(
		http.StatusOK,
		Note{
			Id:       note_id,
			Question: "What is the answer to the ultimate question of life, the universe and everything?",
			Answer:   "42",
		},
	)
}

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/note/:id", GetNote)
	return router
}

func main() {
	router := SetupRouter()
	router.Run("127.0.0.1:8080")
}
