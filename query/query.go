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

package query

import (
	"fmt"
	"github.com/hanahmily/banyandb/config"
	"github.com/hanahmily/banyandb/log"
	"github.com/hanahmily/banyandb/query/graph"
	"github.com/hanahmily/banyandb/storage"
	"github.com/vektah/gqlgen/handler"
	"net/http"
)

type Query struct {
	svr     *http.Server
	storage *storage.Storage
}

func (q *Query) Start(config *config.ServerConfig) error {
	addr := fmt.Sprintf("%s:%d", config.QueryAddr, config.QueryPort)
	log.Infof("Query is listening on %s", addr)
	q.svr = &http.Server{Addr: addr}
	http.Handle("/", handler.Playground("Query", "/query"))
	http.Handle("/query", handler.GraphQL(graph.MakeExecutableSchema(&graph.Query{S: q.storage})))
	go func() {
		log.Error(q.svr.ListenAndServe())
	}()
	return nil
}

func (q *Query) Stop() error {
	log.Info("Query is shutting down")
	return q.svr.Shutdown(nil)
}
