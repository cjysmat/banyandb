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

package storage

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/hanahmily/banyandb/log"
	"os"
)

type ItemType string

const (
	ItemTypeString ItemType = "STRING"
	ItemTypeInt32  ItemType = "INT32"
	ItemTypeInt64  ItemType = "INT64"
	ItemTypeFloat  ItemType = "FLOAT"
	ItemTypeDouble ItemType = "DOUBLE"
	ItemTypeBool   ItemType = "BOOL"
	ItemTypeBytes  ItemType = "BYTES"
)

func NewStorage(dir string) (*Storage, error) {
	rootDir := dir
	createDirIfNotExist(rootDir)
	metaDir := fmt.Sprintf("%s/meta", dir)
	createDirIfNotExist(metaDir)
	return &Storage{
		rootDir: rootDir,
		metaDir: metaDir,
	}, nil
}

type Storage struct {
	rootDir string
	metaDir string
	metaDb  *badger.DB
}

func (s *Storage) Open() error {
	log.Info("Opening storage")
	err := s.openMeta()
	if err != nil {
		return err
	}
	log.Info("Storage has been opened")
	return nil
}

func (s *Storage) Close() error {
	log.Info("Closing storage")
	return s.metaDb.Close()
}

func (s *Storage) CreateLogMeta(name string) *LogMeta {
	db := s.metaDb
	return &LogMeta{entityName: name, txn: db.NewTransaction(true)}
}

func (s *Storage) openMeta() error {
	opts := badger.DefaultOptions
	opts.Dir = s.metaDir
	opts.ValueDir = s.metaDir
	opts.SyncWrites = true
	db, err := badger.Open(opts)
	if err != nil {
		return err
	}
	s.metaDb = db
	return nil
}

type LogMeta struct {
	entityName string
	txn        *badger.Txn
}

func (lm *LogMeta) AddLogMetaItem(itemName string, itemType string) error {
	return lm.txn.Set(metaKey(lm.entityName, itemName), []byte(itemType))
}

func (lm *LogMeta) Finish(err error) error {
	if err != nil {
		lm.txn.Discard()
		return err
	}
	return lm.txn.Commit(nil)
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func metaKey(name string, item string) []byte {
	return []byte(fmt.Sprintf("%s.%s", name, item))
}
