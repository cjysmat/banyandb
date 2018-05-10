/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"context"
	"github.com/hanahmily/banyandb/config"
	"github.com/hanahmily/banyandb/log"
	"os"
)

type Server struct {
	signals  chan os.Signal
	stopChan chan bool
	config   *config.ServerConfig
}

func NewServer(config *config.ServerConfig) *Server {
	server := new(Server)
	server.signals = make(chan os.Signal, 1)
	server.stopChan = make(chan bool, 1)
	server.config = config
	return server
}

func (s *Server) Start(ctx context.Context) {
	defer log.Info("Server started")
	go func() {
		<-ctx.Done()
		log.Info("Stopping server gracefully")
		s.Stop()
	}()
}

func (s *Server) Wait() {
	<-s.stopChan
}

func (s *Server) Stop() {
	defer log.Info("Server stopped")
	s.stopChan <- true
}
