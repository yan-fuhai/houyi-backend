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
	"flag"
	"github.com/houyi-tracing/houyi/ports"
	"github.com/spf13/viper"
)

const (
	strategyManagerAddr = "strategy.manager.addr"
	strategyManagerPort = "strategy.manager.port"
	httpListenPort      = "http.listen.port"
)

const (
	DefaultStrategyManagerAddr = "strategy-manager"
	DefaultStrategyManagerPort = ports.StrategyManagerGrpcListenPort
	DefaultHttpListenPort      = 80
)

type Flags struct {
	StrategyManagerAddr string
	StrategyManagerPort int
	HttpListenPort      int
}

func AddFlags(flags *flag.FlagSet) {
	flags.String(strategyManagerAddr, DefaultStrategyManagerAddr, "Address of strategy manager.")
	flags.Int(strategyManagerPort, DefaultStrategyManagerPort, "Port to serve gRPC for strategy manager.")
	flags.Int(httpListenPort, DefaultHttpListenPort, "Port to serve HTTP for backend.")
}

func (opts *Flags) InitFromViper(v *viper.Viper) *Flags {
	opts.StrategyManagerAddr = v.GetString(strategyManagerAddr)
	opts.StrategyManagerPort = v.GetInt(strategyManagerPort)
	opts.HttpListenPort = v.GetInt(httpListenPort)
	return opts
}
