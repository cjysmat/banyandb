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

package graph

import (
	"context"
	"fmt"
	"github.com/hanahmily/banyandb/log"
	"github.com/hanahmily/banyandb/query/graph/schema"
	"github.com/hanahmily/banyandb/storage"
)

type Query struct {
	S *storage.Storage
}

func (l *Query) Mutation_createLogEntity(ctx context.Context, logMeta schema.LogMetaInput) (string, error) {
	log.Infof("Creating Log Metadata: %v", logMeta)
	meta := l.S.CreateLogMeta(logMeta.Name)
	var err error = nil
	for _, v := range logMeta.LogItems {
		err = meta.AddLogMetaItem(v.Name, v.Type.String())
		if err != nil {
			break
		}
	}
	return fmt.Sprintf("Create Log entity %s successfully", logMeta.Name), meta.Finish(err)
}

func (l *Query) Query_log(ctx context.Context) ([]schema.LogItem, error) {
	return []schema.LogItem{}, nil
}
