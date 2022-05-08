package main

import (
	"bytes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/sonntuet1997/numberlink-sat/numberlink"
	"github.com/sonntuet1997/numberlink-sat/solver"
)

type Data struct {
	Data string `form:"data" json:"data" xml:"data"  binding:"required"`
}

func init() {

}

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/solve", func(c *gin.Context) {
		var json Data
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		board := numberlink.NewFromString(json.Data)
		runtime := solver.SolveWithCustomSolver(board, "cadical -q", "normal")
		var buff bytes.Buffer
		board.Print(&buff)
		c.JSON(http.StatusOK, gin.H{"result": buff.String(), "runtime": runtime})
	})
	err := router.Run(":80")
	if err != nil {
		return
	}
}
