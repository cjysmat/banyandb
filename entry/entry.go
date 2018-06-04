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

package entry

import (
	"github.com/hanahmily/banyandb/config"
	"net"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/hanahmily/banyandb/entry/grpc"
	"github.com/hanahmily/banyandb/log"
)

func NewEntry() {

}

type Entry struct {
	s *grpc.Server
}

func (e *Entry) Start(config *config.ServerConfig) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.EntryAddr, config.EntryPort))
	if err != nil {
		return err
	}
	e.s = grpc.NewServer()
	pb.RegisterLogServiceServer(e.s, &pb.Logger{})
	go func() {
		log.Error(e.s.Serve(lis))
	}()
	return nil
}

func (e *Entry) Stop() error {
	e.s.GracefulStop()
	return nil
}
