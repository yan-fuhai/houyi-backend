// Copyright (c) 2021 The Houyi Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"go.uber.org/zap"
)

type HttpServerParams struct {
	Logger              *zap.Logger
	StrategyManagerAddr string
	StrategyManagerPort int
	HttpListenPort      int
}

type HttpServer struct {
	Logger            *zap.Logger
	HttpHandler       *HttpHandler
	StrategyManagerEp routing.Endpoint
	HttpListenPort    int
	engine            *gin.Engine
}

func NewHttpServer(params *HttpServerParams) *HttpServer {
	return &HttpServer{
		Logger: params.Logger,
		StrategyManagerEp: routing.Endpoint{
			Addr: params.StrategyManagerAddr,
			Port: params.StrategyManagerPort,
		},
		HttpListenPort: params.HttpListenPort,
		engine:         gin.Default(),
	}
}

func (s *HttpServer) StartHttpServer() error {
	handler := HttpHandler{
		Logger:            s.Logger,
		StrategyManagerEp: s.StrategyManagerEp,
	}
	s.engine.GET("/getServices", handler.GetServices)
	s.engine.GET("/getOperations", handler.GetOperations)
	s.engine.GET("/getCausalRelations", handler.GetCausalRelations)
	s.engine.GET("/getTags", handler.GetTags)
	s.engine.POST("/updateTags", handler.UpdateTags)
	s.engine.Run(fmt.Sprintf(":%d", s.HttpListenPort))

	return nil
}
