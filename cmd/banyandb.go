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

package main

import (
	"context"
	"github.com/containous/flaeg"
	"github.com/hanahmily/banyandb/config"
	"github.com/hanahmily/banyandb/log"
	"github.com/hanahmily/banyandb/server"
	bStorage "github.com/hanahmily/banyandb/storage"
	fmtLog "log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	banyanConfig := config.NewBanyanConfig()
	banyanDefaultPointerConfig := config.NewBanyanDefaultPointersConfig()
	banyanCmd := &flaeg.Command{
		Name:                  "BanyanDB",
		Description:           `BanyanDB is an APM database`,
		Config:                banyanConfig,
		DefaultPointersConfig: banyanDefaultPointerConfig,
		Run: func() error {
			runCmd(banyanConfig)
			return nil
		},
	}
	f := flaeg.New(banyanCmd, os.Args[1:])
	if err := f.Run(); err != nil {
		fmtLog.Printf("Error %s", err.Error())
		os.Exit(1)
	}
}

func runCmd(globalConfig *config.BanyanConfig) {
	storage, err := bStorage.NewStorage(globalConfig.Dir)
	if err != nil {
		log.Fatal(err)
	}
	storage.Open()
	newCtx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-signals:
			cancel()
		}
	}()

	svr := server.NewServer(globalConfig.Server, storage)
	svr.Start(newCtx)
	svr.Wait()
	storage.Close()
	log.Info("Shut down")
	log.Exit()
}
