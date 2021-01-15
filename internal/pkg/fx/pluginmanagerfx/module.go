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

package pluginmanagerfx

import (
	"fmt"
	p "github.com/xmidt-org/ears/internal/pkg/plugin"
	"github.com/xmidt-org/ears/pkg/plugin/manager"
	"github.com/xmidt-org/ears/pkg/plugins/debug"
	"github.com/xmidt-org/ears/pkg/plugins/match"

	pkgplugin "github.com/xmidt-org/ears/pkg/plugin"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		ProvidePluginManager,
	),
)

type PluginIn struct {
	fx.In
}

type PluginOut struct {
	fx.Out

	PluginManager p.Manager
}

func ProvidePluginManager(in PluginIn) (PluginOut, error) {
	out := PluginOut{}

	mgr, err := manager.New()
	if err != nil {
		return out, fmt.Errorf("could not provide plugin manager: %w", err)
	}

	// Go ahead and register some default plugins
	defaultPlugins := []struct {
		name   string
		plugin pkgplugin.Pluginer
	}{
		{
			name:   "debug",
			plugin: debug.NewPluginVersion("debug", "", ""),
		},

		{
			name:   "match",
			plugin: match.NewPluginVersion("match", "", ""),
		},
	}

	for _, plug := range defaultPlugins {
		err = mgr.RegisterPlugin(plug.name, plug.plugin)
		if err != nil {
			return out, fmt.Errorf("could register %s plugin: %w", plug.name, err)
		}
	}

	m, err := p.NewManager(p.WithPluginManager(mgr))
	if err != nil {
		return out, fmt.Errorf("could not provide plugin manager: %w", err)
	}

	out.PluginManager = m

	return out, nil

}
