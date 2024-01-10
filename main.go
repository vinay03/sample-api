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

type TestServerDummyResponse struct {
	Message   string            `json:"message"`
	Headers   map[string]string `json:"_headers"`
	ReplicaId int               `json:"replicaId"`
	Host      string            `json:"host"`
}

type TestServerDummyDelayedResponse struct {
	Message   string `json:"message"`
	ReplicaId int    `json:"replicaId"`
}

func GetNumberedHandler(ReplicaNumber int) func(*gin.Context) {
	return func(c *gin.Context) {
		response := TestServerDummyResponse{
			Message:   fmt.Sprintf("Response to URI '%v' from Replica #%v", c.Request.URL, ReplicaNumber),
			Host:      c.Request.Host,
			ReplicaId: ReplicaNumber,
		}
		response.Headers = make(map[string]string)
		for name, values := range c.Request.Header {
			for _, value := range values {
				response.Headers[name] = value
			}
		}

		c.JSON(http.StatusOK, response)

	}
}

func GetDelayedHandler(ReplicaNumber int, defaultDelayInterval time.Duration) func(*gin.Context) {
	return func(c *gin.Context) {
		delayInterval := defaultDelayInterval

		if c.Request.Method == "POST" {
			inputData := struct {
				Delay int `json:"delay"`
			}{}

			err := c.BindJSON(&inputData)
			if err == nil {
				if inputData.Delay > 0 && inputData.Delay <= 20 {
					delayInterval = time.Duration(inputData.Delay * int(time.Second))
				}
			}
		}

		if delayInterval > 0 {
			log.Printf("Starting wait for '%v'\n", delayInterval)
			time.Sleep(delayInterval)
			log.Printf("Wait of '%v' ended\n", delayInterval)
		}

		response := TestServerDummyResponse{
			Message:   fmt.Sprintf("Response to URI '%v' from Replica #%v", c.Request.URL, ReplicaNumber),
			Host:      c.Request.Host,
			ReplicaId: ReplicaNumber,
		}
		response.Headers = make(map[string]string)
		for name, values := range c.Request.Header {
			for _, value := range values {
				response.Headers[name] = value
			}
		}

		c.JSON(http.StatusOK, response)
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
	serverPortStart := 8090
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

		healthHandlerFunc := GetDelayedHandler(ReplicaNumber, 0*time.Second)
		router.GET("/health", healthHandlerFunc)
		router.POST("/health", healthHandlerFunc)

		delayedHandlerFunc := GetDelayedHandler(ReplicaNumber, 0*time.Second)
		router.GET("/delayed", delayedHandlerFunc)
		router.POST("/delayed", delayedHandlerFunc)

		fmt.Printf("API for replica #%v started listening at localhost:%v\n", ReplicaNumber, serverPortStart+ReplicaNumber)
		go router.Run(fmt.Sprintf("localhost:%v", serverPortStart+ReplicaNumber))
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done // Will block here until user hits ctrl+c
}
