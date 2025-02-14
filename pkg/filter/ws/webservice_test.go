// Copyright 2021 Comcast Cable Communications Management, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ws_test

import (
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/xmidt-org/ears/internal/pkg/appsecret"
	"github.com/xmidt-org/ears/pkg/event"
	"github.com/xmidt-org/ears/pkg/filter/ws"
	"github.com/xmidt-org/ears/pkg/tenant"
	"github.com/xorcare/pointer"
	"reflect"
	"testing"
)

func TestFilterWsBasic(t *testing.T) {
	ctx := context.Background()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatalf("failed to load test configuration %s\n", err.Error())
	}
	secrets := appsecret.NewConfigVault(v)
	f, err := ws.NewFilter(tenant.Id{AppId: "myapp", OrgId: "myorg"}, "ws", "myws", ws.Config{
		ToPath:                 ".value",
		Url:                    "http://echo.jsontest.com/key/value/one/two",
		Method:                 "GET",
		EmptyPathValueRequired: pointer.Bool(true),
	}, secrets)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	eventStr := `{"foo":"bar"}`
	var obj interface{}
	err = json.Unmarshal([]byte(eventStr), &obj)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	e, err := event.New(ctx, obj, event.FailOnNack(t))
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	evts := f.Filter(e)
	if len(evts) != 1 {
		t.Fatalf("wrong number of encoded events: %d\n", len(evts))
	}
	expectedEventStr := `{ "foo" : "bar", "value": { "one": "two", "key": "value" } }`
	var res interface{}
	err = json.Unmarshal([]byte(expectedEventStr), &res)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	if !reflect.DeepEqual(evts[0].Payload(), res) {
		pl, _ := json.MarshalIndent(evts[0].Payload(), "", "\t")
		t.Fatalf("wrong payload in encoded event: %s\n", pl)
	}
}

func TestFilterWsUrlEval(t *testing.T) {
	ctx := context.Background()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatalf("failed to load test configuration %s\n", err.Error())
	}
	secrets := appsecret.NewConfigVault(v)
	f, err := ws.NewFilter(tenant.Id{AppId: "myapp", OrgId: "myorg"}, "ws", "myws", ws.Config{
		ToPath: ".value",
		Url:    "{.url}",
		Method: "GET",
	}, secrets)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	eventStr := `{"foo":"bar", "url":"http://echo.jsontest.com/key/value/one/two"}`
	var obj interface{}
	err = json.Unmarshal([]byte(eventStr), &obj)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	e, err := event.New(ctx, obj, event.FailOnNack(t))
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	evts := f.Filter(e)
	if len(evts) != 1 {
		t.Fatalf("wrong number of encoded events: %d\n", len(evts))
	}
	expectedEventStr := `{ "foo" : "bar", "url":"http://echo.jsontest.com/key/value/one/two", "value": { "one": "two", "key": "value" } }`
	var res interface{}
	err = json.Unmarshal([]byte(expectedEventStr), &res)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	if !reflect.DeepEqual(evts[0].Payload(), res) {
		pl, _ := json.MarshalIndent(evts[0].Payload(), "", "\t")
		t.Fatalf("wrong payload in encoded event: %s\n", pl)
	}
}

func TestFilterWsEmptyPathValueNotMet(t *testing.T) {
	ctx := context.Background()
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		t.Fatalf("failed to load test configuration %s\n", err.Error())
	}
	secrets := appsecret.NewConfigVault(v)
	f, err := ws.NewFilter(tenant.Id{AppId: "myapp", OrgId: "myorg"}, "ws", "myws", ws.Config{
		ToPath:                 ".value",
		Url:                    "http://echo.jsontest.com/key/value/one/two",
		Method:                 "GET",
		EmptyPathValueRequired: pointer.Bool(true),
	}, secrets)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	eventStr := `{"foo":"bar", "value": "alreadyHere"}`
	var obj interface{}
	err = json.Unmarshal([]byte(eventStr), &obj)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	e, err := event.New(ctx, obj, event.FailOnNack(t))
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	evts := f.Filter(e)
	if len(evts) != 1 {
		t.Fatalf("wrong number of encoded events: %d\n", len(evts))
	}
	expectedEventStr := `{ "foo" : "bar", "value": "alreadyHere" }`
	var res interface{}
	err = json.Unmarshal([]byte(expectedEventStr), &res)
	if err != nil {
		t.Fatalf("encode test failed: %s\n", err.Error())
	}
	if !reflect.DeepEqual(evts[0].Payload(), res) {
		pl, _ := json.MarshalIndent(evts[0].Payload(), "", "\t")
		t.Fatalf("wrong payload in encoded event: %s\n", pl)
	}
}
