package internal

import (
	"context"
	"errors"
)

type (
	InMemoryRoutingTableManager struct {
		//TODO: add index by source plugin
		routingTableIndex RoutingTableIndex
	}
)

// NewInMemoryRoutingTableManager creates a new local in memory routing table cache
func NewInMemoryRoutingTableManager() *InMemoryRoutingTableManager {
	mgr := new(InMemoryRoutingTableManager)
	mgr.routingTableIndex = make(map[string]*RoutingTableEntry)
	return mgr
}

func (mgr *InMemoryRoutingTableManager) AddRoute(ctx *context.Context, entry *RoutingTableEntry) error {
	if entry == nil {
		return errors.New("missing routing table entry")
	}
	if err := entry.Validate(); err != nil {
		return err
	}
	mgr.routingTableIndex[entry.Hash()] = entry
	return nil
}

func (mgr *InMemoryRoutingTableManager) RemoveRoute(ctx *context.Context, entry *RoutingTableEntry) error {
	if entry == nil {
		return errors.New("missing routing table entry")
	}
	if err := entry.Validate(); err != nil {
		return err
	}
	delete(mgr.routingTableIndex, entry.Hash())
	return nil
}

func (mgr *InMemoryRoutingTableManager) ReplaceAllRoutes(ctx *context.Context, entries []*RoutingTableEntry) error {
	mgr.routingTableIndex = make(map[string]*RoutingTableEntry)
	for _, entry := range entries {
		if err := entry.Validate(); err != nil {
			return err
		}
		delete(mgr.routingTableIndex, entry.Hash())
	}
	return nil
}

func (mgr *InMemoryRoutingTableManager) Validate() error {
	for _, entry := range mgr.routingTableIndex {
		if err := entry.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (mgr *InMemoryRoutingTableManager) Hash() string {
	hash := ""
	for _, entry := range mgr.routingTableIndex {
		hash = hash + entry.Hash()
	}
	return hash
}

func (mgr *InMemoryRoutingTableManager) GetAllRoutes(ctx *context.Context) ([]*RoutingTableEntry, error) {
	tbl := make([]*RoutingTableEntry, len(mgr.routingTableIndex))
	idx := 0
	for _, entry := range mgr.routingTableIndex {
		tbl[idx] = entry
	}
	return tbl, nil
}

func (mgr *InMemoryRoutingTableManager) GetRoutesBySourcePlugin(ctx *context.Context, plugin *Plugin) ([]*RoutingTableEntry, error) {
	tbl := make([]*RoutingTableEntry, len(mgr.routingTableIndex))
	for _, entry := range mgr.routingTableIndex {
		if entry.Source.Hash() == plugin.Hash() {
			tbl = append(tbl, entry)
		}
	}
	return tbl, nil
}

func (mgr *InMemoryRoutingTableManager) GetRoutesByDestinationPlugin(ctx *context.Context, plugin *Plugin) ([]*RoutingTableEntry, error) {
	tbl := make([]*RoutingTableEntry, len(mgr.routingTableIndex))
	for _, entry := range mgr.routingTableIndex {
		if entry.Destination.Hash() == plugin.Hash() {
			tbl = append(tbl, entry)
		}
	}
	return tbl, nil
}

func (mgr *InMemoryRoutingTableManager) GetRoutesForEvent(ctx *context.Context, event *Event) ([]*RoutingTableEntry, error) {
	return mgr.GetRoutesBySourcePlugin(ctx, event.Source)
}
