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
	"context"
	"github.com/gin-gonic/gin"
	"github.com/houyi-tracing/houyi/idl/api_v1"
	"github.com/houyi-tracing/houyi/pkg/routing"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

const (
	EqualTo              = "=="
	NotEqualTo           = "!="
	GreaterThan          = ">"
	GreaterThanOrEqualTo = ">="
	LessThan             = "<"
	LessThanOrEqualTo    = "<="
)

type Tag struct {
	Name     string      `json:"name"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

type HttpHandler struct {
	Logger            *zap.Logger
	StrategyManagerEp routing.Endpoint
}

func (h *HttpHandler) GetServices(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	conn, err := grpc.Dial(h.StrategyManagerEp.String(), grpc.WithInsecure())
	if conn == nil || err != nil {
		h.Logger.Error("failed to dial to strategy manager", zap.Error(err))
	} else {
		defer conn.Close()
	}

	client := api_v1.NewTraceGraphManagerClient(conn)
	resp, err := client.GetServices(context.TODO(), &api_v1.GetServicesRequest{})
	if resp != nil && err == nil {
		if services := resp.GetServices(); services != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": services,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": []string{},
			})
		}
	} else {
		h.Logger.Error("failed to get services", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": []string{},
		})
	}
}

func (h *HttpHandler) GetOperations(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	svc := c.Query("service")
	if svc != "" {
		h.Logger.Debug("GetOperations", zap.String("service name", svc))
	} else {
		h.Logger.Debug("GetOperations: service is empty")
	}

	conn, err := grpc.Dial(h.StrategyManagerEp.String(), grpc.WithInsecure())
	if conn == nil || err != nil {
		h.Logger.Error("failed to dial to strategy manager", zap.Error(err))
	} else {
		defer conn.Close()
	}

	client := api_v1.NewTraceGraphManagerClient(conn)
	resp, err := client.GetOperations(context.TODO(), &api_v1.GetOperationsRequest{Service: svc})
	if resp != nil && err == nil {
		if operations := resp.GetOperations(); operations != nil {
			c.JSON(http.StatusOK, gin.H{
				"result": operations,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": []string{},
			})
		}
	} else {
		h.Logger.Error("failed to get operations", zap.Error(err))
	}
}

func (h *HttpHandler) GetTags(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	conn, err := grpc.Dial(h.StrategyManagerEp.String(), grpc.WithInsecure())
	if conn == nil || err != nil {
		h.Logger.Error("failed to dial to strategy manager", zap.Error(err))
	} else {
		defer conn.Close()
	}

	client := api_v1.NewEvaluatorManagerClient(conn)
	resp, err := client.GetTags(context.TODO(), &api_v1.GetTagsRequest{})
	if resp != nil && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": convertToJsonTags(resp.GetTags()),
		})
	} else {
		h.Logger.Error("failed to get tags", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": []*api_v1.EvaluatingTag{},
		})
	}
}

func (h *HttpHandler) GetCausalRelations(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	svc := c.Query("service")
	op := c.Query("operation")
	if svc == "" {
		h.Logger.Debug("GetTags: Empty service")
		c.Status(http.StatusBadRequest)
		return
	}
	if op == "" {
		h.Logger.Debug("GetTags: Empty operation")
		c.Status(http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial(h.StrategyManagerEp.String(), grpc.WithInsecure())
	if conn == nil || err != nil {
		h.Logger.Error("failed to dial to strategy manager", zap.Error(err))
	} else {
		defer conn.Close()
	}

	client := api_v1.NewTraceGraphManagerClient(conn)
	resp, err := client.Traces(context.TODO(), &api_v1.Operation{
		Service:   svc,
		Operation: op,
	})
	if resp != nil && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"result": resp.GetEntries(),
		})
	} else {
		h.Logger.Error("failed to get causal relations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": []*api_v1.TraceNode{},
		})
	}
}

func (h *HttpHandler) UpdateTags(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	tags := make([]Tag, 0)

	err := c.BindJSON(&tags)
	if err == nil {
		conn, err := grpc.Dial(h.StrategyManagerEp.String(), grpc.WithInsecure())
		if conn == nil || err != nil {
			h.Logger.Error("failed to dial to strategy manager", zap.Error(err))
		} else {
			defer conn.Close()
		}

		client := api_v1.NewEvaluatorManagerClient(conn)
		resp, err := client.UpdateTags(context.TODO(), &api_v1.UpdateTagsRequest{
			Tags: convertToTags(tags),
		})
		if resp != nil && err == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "OK",
			})
		} else {
			h.Logger.Error("failed to update tags", zap.Error(err))
		}
	} else {
		h.Logger.Error("failed to parse JSON from request's body", zap.Error(err))
	}
}

func convertToJsonTags(tags []*api_v1.EvaluatingTag) []Tag {
	ret := make([]Tag, 0)
	for _, t := range tags {
		newTag := Tag{}
		newTag.Name = t.TagName

		switch t.OperationType {
		case api_v1.EvaluatingTag_EQUAL_TO:
			newTag.Operator = EqualTo
		case api_v1.EvaluatingTag_NOT_EQUAL_TO:
			newTag.Operator = NotEqualTo
		case api_v1.EvaluatingTag_GREATER_THAN:
			newTag.Operator = GreaterThan
		case api_v1.EvaluatingTag_GREATER_THAN_OR_EQUAL_TO:
			newTag.Operator = GreaterThanOrEqualTo
		case api_v1.EvaluatingTag_LESS_THAN:
			newTag.Operator = LessThan
		case api_v1.EvaluatingTag_LESS_THAN_OR_EQUAL_TO:
			newTag.Operator = LessThanOrEqualTo
		}

		switch t.ValueType {
		case api_v1.EvaluatingTag_INTEGER:
			newTag.Value = t.GetIntegerVal()
		case api_v1.EvaluatingTag_STRING:
			newTag.Value = t.GetStringVal()
		case api_v1.EvaluatingTag_FLOAT:
			newTag.Value = t.GetFloatVal()
		case api_v1.EvaluatingTag_BOOLEAN:
			newTag.Value = t.GetBooleanVal()
		}

		ret = append(ret, newTag)
	}
	return ret
}

func convertToTags(tags []Tag) []*api_v1.EvaluatingTag {
	ret := make([]*api_v1.EvaluatingTag, 0)
	for _, t := range tags {
		newTag := &api_v1.EvaluatingTag{}
		newTag.TagName = t.Name

		switch t.Operator {
		case EqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_EQUAL_TO
		case NotEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_NOT_EQUAL_TO
		case GreaterThan:
			newTag.OperationType = api_v1.EvaluatingTag_GREATER_THAN
		case GreaterThanOrEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_GREATER_THAN_OR_EQUAL_TO
		case LessThan:
			newTag.OperationType = api_v1.EvaluatingTag_LESS_THAN
		case LessThanOrEqualTo:
			newTag.OperationType = api_v1.EvaluatingTag_LESS_THAN_OR_EQUAL_TO
		}

		switch val := t.Value.(type) {
		case int64:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: val}
		case int32:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int16:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int8:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case int:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint64:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint32:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint16:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case uint8:
			newTag.ValueType = api_v1.EvaluatingTag_INTEGER
			newTag.Value = &api_v1.EvaluatingTag_IntegerVal{IntegerVal: int64(val)}
		case float64:
			newTag.ValueType = api_v1.EvaluatingTag_FLOAT
			newTag.Value = &api_v1.EvaluatingTag_FloatVal{FloatVal: val}
		case float32:
			newTag.ValueType = api_v1.EvaluatingTag_FLOAT
			newTag.Value = &api_v1.EvaluatingTag_FloatVal{FloatVal: float64(val)}
		case string:
			newTag.ValueType = api_v1.EvaluatingTag_STRING
			newTag.Value = &api_v1.EvaluatingTag_StringVal{StringVal: val}
		case bool:
			newTag.ValueType = api_v1.EvaluatingTag_BOOLEAN
			newTag.Value = &api_v1.EvaluatingTag_BooleanVal{BooleanVal: val}
		}

		ret = append(ret, newTag)
	}
	return ret
}
