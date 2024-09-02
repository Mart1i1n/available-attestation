package library

import (
	"github.com/tsinghua-cel/strategy-gen/library/five"
	"github.com/tsinghua-cel/strategy-gen/library/four"
	"github.com/tsinghua-cel/strategy-gen/library/one"
	"github.com/tsinghua-cel/strategy-gen/library/three"
	"github.com/tsinghua-cel/strategy-gen/library/two"
	"github.com/tsinghua-cel/strategy-gen/types"
	"sync"
)

type Strategy interface {
	Run(param types.LibraryParams)
	Description() string
}

var (
	allStrategies = sync.Map{}
)

func Init() {
	allStrategies.Store("one", &one.One{})
	allStrategies.Store("two", &two.Two{})
	allStrategies.Store("three", &three.Three{})
	allStrategies.Store("four", &four.Four{})
	allStrategies.Store("five", &five.Five{})
}

func GetStrategy(name string) (Strategy, bool) {
	if v, ok := allStrategies.Load(name); ok {
		return v.(Strategy), true
	}
	return nil, false
}

func GetAllStrategies() map[string]Strategy {
	strategies := make(map[string]Strategy)
	allStrategies.Range(func(k, v interface{}) bool {
		strategies[k.(string)] = v.(Strategy)
		return true
	})
	return strategies
}
