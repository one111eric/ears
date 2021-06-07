// Copyright 2020 Comcast Cable Communications Management, LLC
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

package split

import (
	"errors"
	"fmt"
	"github.com/xmidt-org/ears/internal/pkg/rtsemconv"
	"github.com/xmidt-org/ears/pkg/event"
	"github.com/xmidt-org/ears/pkg/filter"
	"go.opentelemetry.io/otel"
)

// a SplitFilter splits and event into two or more events
func NewFilter(config interface{}) (*Filter, error) {
	cfg, err := NewConfig(config)
	if err != nil {
		return nil, &filter.InvalidConfigError{
			Err: err,
		}
	}
	cfg = cfg.WithDefaults()
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	f := &Filter{
		config: *cfg,
	}
	return f, nil
}

// Filter splits an event containing an array into multiple events
func (f *Filter) Filter(evt event.Event) []event.Event {
	if f == nil {
		evt.Nack(&filter.InvalidConfigError{
			Err: fmt.Errorf("<nil> pointer filter")})
		return nil
	}
	if evt.Trace() {
		tracer := otel.Tracer(rtsemconv.EARSTracerName)
		_, span := tracer.Start(evt.Context(), "splitFilter")
		defer span.End()
	}
	events := []event.Event{}
	obj, _, _ := evt.GetPathValue(f.config.Path)
	if obj == nil {
		evt.Nack(errors.New("nil object at path " + f.config.Path))
		return []event.Event{}
	}
	arr, ok := obj.([]interface{})
	if !ok {
		evt.Nack(errors.New("split on non array type"))
		return nil
	}
	for _, p := range arr {
		nevt, err := evt.Clone(evt.Context())
		if err != nil {
			evt.Nack(err)
			return nil
		}
		err = nevt.SetPayload(p)
		if err != nil {
			evt.Nack(err)
			return nil
		}
		events = append(events, nevt)
	}
	evt.Ack()
	return events
}

func (f *Filter) Config() interface{} {
	if f == nil {
		return Config{}
	}
	return f.config
}

func (f *Filter) Name() string {
	return ""
}

func (f *Filter) Plugin() string {
	return "split"
}
