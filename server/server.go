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
	"github.com/hanahmily/banyandb/query"
)

type Endpoint interface {
	Start(config *config.ServerConfig) error
	Stop() error
}

type Server struct {
	signals  chan os.Signal
	stopChan chan bool
	config   *config.ServerConfig
	endpoints []Endpoint
}

func NewServer(config *config.ServerConfig) *Server {
	server := new(Server)
	server.signals = make(chan os.Signal, 1)
	server.stopChan = make(chan bool, 1)
	server.config = config
	server.endpoints = append(server.endpoints, &query.Query{})

	return server
}

func (s *Server) setupEndpoints() error {
	for _, v := range s.endpoints {
		err := v.Start(s.config)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) tearDownEndpoints() {
	for _, v := range s.endpoints {
		err := v.Stop()
		if err != nil {
			log.Error(err)
		}
	}
}

func (s *Server) Start(ctx context.Context) {
	defer log.Info("Server started")
	go func() {
		<-ctx.Done()
		log.Info("Stopping server gracefully")
		s.Stop()
	}()
	err := s.setupEndpoints()
	if err != nil {
		log.Error(err)
		s.Stop()
	}
}

func (s *Server) Wait() {
	<-s.stopChan
}

func (s *Server) Stop() {
	defer log.Info("Server stopped")
	s.tearDownEndpoints()
	s.stopChan <- true
}
