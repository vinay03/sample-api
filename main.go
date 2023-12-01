package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	argsWithoutProg := os.Args[1:]
	serverPort := "8080"

	if len(argsWithoutProg) > 0 {
		serverPort = argsWithoutProg[0]
	}
	ReplicaNumber := "0"
	handlerFunc := func(c *gin.Context) {
		ReplicaNumber = os.Getenv("ReplicaNumber")
		c.JSON(http.StatusOK, fmt.Sprintf("Response to URI '%v' from Replica #%v", c.Request.URL, ReplicaNumber))
	}

	router.GET("/", handlerFunc)
	router.Run("localhost:" + serverPort)
	fmt.Printf("API for replica #%v started", ReplicaNumber)
}
