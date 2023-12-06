package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func GetNumberedHandler(ReplicaNumber int) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Response to URI '%v' from Replica #%v", c.Request.URL, ReplicaNumber),
		})
	}
}

func getHealthHandlerFunc() func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"active": true,
		})
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	serverPortStart := 8080
	// gin.SetMode(gin.ReleaseMode)

	replicas := 1
	var err error
	if len(argsWithoutProg) > 0 {
		replicas, err = strconv.Atoi(argsWithoutProg[0])
		if err != nil {
			replicas = 1
		}
	}
	fmt.Printf("Starting %v replicas\n", replicas)
	for ReplicaNumber := 1; ReplicaNumber <= replicas; ReplicaNumber++ {
		router := gin.Default()
		router.Use(gin.Recovery())

		handlerFunc := GetNumberedHandler(ReplicaNumber)
		router.GET("/", handlerFunc)

		healthHandlerFunc := getHealthHandlerFunc()
		router.GET("/health", healthHandlerFunc)

		fmt.Printf("API for replica #%v started\n", ReplicaNumber)
		go router.Run(fmt.Sprintf("localhost:%v", serverPortStart+ReplicaNumber))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done // Will block here until user hits ctrl+c
}
