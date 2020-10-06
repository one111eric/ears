/**
 *  Copyright (c) 2020  Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xmidt-org/ears/internal"

	"github.com/rs/zerolog/log"
)

var (
	ROUTE = `
	{
		"orgId" : "comcast",
		"appId" : "xfi",
		"userId" : "boris",
		"srcType" : "debug",
		"srcParams" :
		{
			"rounds" : 10,
			"intervalMS" : 250,
			"payload" : {
				"foo" : "bar"
			}
		},
		"dstType" : "debug",
		"dstParams" : {},
		"filterChain" : [
			{
				"type" : "match",
				"params" : {
					"pattern" : {
						"foo" : "bar"
					}
				}
			},
			{
				"type" : "filter",
				"params" : {
					"pattern" : {
						"hello" : "world"
					}
				}
			},
			{
				"type" : "split",
				"params" : {}
			},
			{
				"type" : "transform",
				"params" : {}
			}
		],
		"deliveryMode" : "at_least_once"
	}
	`
)

func main() {
	ctx := context.Background()
	var rtmgr internal.RoutingTableManager
	rtmgr = internal.NewInMemoryRoutingTableManager()
	log.Debug().Msg(fmt.Sprintf("ears has %d routes", rtmgr.GetRouteCount(ctx)))
	var rte internal.RoutingTableEntry
	err := json.Unmarshal([]byte(ROUTE), &rte)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	buf, err := json.MarshalIndent(rte, "", "\t")
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	fmt.Printf("%s\n", string(buf))
	err = rtmgr.AddRoute(ctx, &rte)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Debug().Msg(fmt.Sprintf("ears has %d routes", rtmgr.GetRouteCount(ctx)))
	allRoutes, err := rtmgr.GetAllRoutes(ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	if len(allRoutes) > 0 {
		log.Debug().Msg(fmt.Sprintf("first route has hash %s", allRoutes[0].Hash(ctx)))
	}
	time.Sleep(time.Duration(60) * time.Second)
}
