package adapter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//Adapter adapts to cosmos sdk modules, and add extra features via hooks
type Adapter interface {
	InitGenesis(ctx sdk.Context) error
	Name() string
}

type Manager struct {
	Adapters         map[string]Adapter
	OrderInitGenesis []string
}

func NewManager(adapters ...Adapter) *Manager {
	adapterMap := make(map[string]Adapter)
	adaptersStr := make([]string, 0, len(adapters))
	for _, adapter := range adapters {
		adapterMap[adapter.Name()] = adapter
		adaptersStr = append(adaptersStr, adapter.Name())
	}
	return &Manager{
		Adapters:         adapterMap,
		OrderInitGenesis: adaptersStr,
	}
}

func (m Manager) InitGenesis(ctx sdk.Context) error {
	for _, moduleName := range m.OrderInitGenesis {
		if m.Adapters[moduleName] == nil {
			continue
		}
		if err := m.Adapters[moduleName].InitGenesis(ctx); err != nil {
			return err
		}
	}
	return nil
}
