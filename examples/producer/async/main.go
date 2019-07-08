/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/apache/rocketmq-client-go/internal/producer"
	"github.com/apache/rocketmq-client-go/primitive"
)

// Package main implements a async producer to send message.
func main() {
	nameServerAddr := []string{"127.0.0.1:9876"}
	p, _ := producer.NewProducer(nameServerAddr, primitive.WithRetry(2))
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		err := p.SendAsync(context.Background(), &primitive.Message{
			Topic:      "TopicTest",
			Body:       []byte("Hello RocketMQ Go Client!"),
			Properties: map[string]string{"id": strconv.Itoa(i)},
		}, func(ctx context.Context, result *primitive.SendResult, e error) {
			if e != nil {
				fmt.Printf("receive message error: %s\n", err)
			} else {
				fmt.Printf("send message success: result=%s\n", result.String())
			}
			wg.Done()
		})

		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		}
	}
	wg.Wait()
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shundown producer error: %s", err.Error())
	}
}