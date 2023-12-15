package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func GetNumberedHandler(ReplicaNumber int) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Response to URI '%v' from Replica #%v", c.Request.URL, ReplicaNumber),
		})
	}
}

func GetDelayedHandler(ReplicaNumber int) func(*gin.Context) {
	return func(c *gin.Context) {
		log.Println("Starting wait...")
		time.Sleep(20 * time.Second)
		log.Println("ending wait...")
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
	gin.SetMode(gin.ReleaseMode)

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
		router.GET("/:path", handlerFunc)

		healthHandlerFunc := getHealthHandlerFunc()
		router.GET("/health", healthHandlerFunc)

		delayedHandlerFunc := GetDelayedHandler(ReplicaNumber)
		router.GET("/delayed", delayedHandlerFunc)

		fmt.Printf("API for replica #%v started\n", ReplicaNumber)
		go router.Run(fmt.Sprintf("localhost:%v", serverPortStart+ReplicaNumber))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done // Will block here until user hits ctrl+c
}
