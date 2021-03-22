package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

const (
	StrategyManagerAddr = "192.168.31.77"
	StrategyManagerGrpcPort = 18760
)

var (
	logger, _ = zap.NewProduction()
)

func GetServices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"result": []string{
			"a", "b", "c", "d", "e", "f",
		},
	})
}

func GetTags(c *gin.Context) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", StrategyManagerAddr, StrategyManagerGrpcPort), grpc.WithInsecure(), grpc.WithBlock())
	if conn == nil || err != nil {
		logger.Error("", zap.Error(err))

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "can not dial to strategy manager",
		})
		return
	} else {
		defer conn.Close()
	}

	client := api_v1.NewEvaluatorManagerClient(conn)
	resp, err := client.GetTags(context.TODO(), &api_v1.GetTagsRequest{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get tags from strategy manager",
		})
	}
	fmt.Println(resp)
}

func main() {
	r := gin.Default()
	r.GET("/getServices", GetServices)
	r.GET("/get_tags", GetTags)
	r.Run() // listen and serve on 0.0.0.0:8080
}
