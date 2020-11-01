package spendingview

import (
	"sync"
)

type SpendingView struct {
	Current map[string]float64 `json:"current"`
	Limits  map[string]float64 `json:"limits"`
	Annual  map[string]float64 `json:"annual"`
}

type Container struct {
	view *SpendingView
	lock *sync.RWMutex
}

func (c *Container) Read() *SpendingView {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.view
}

func (c *Container) Write(view *SpendingView) {
	c.lock.Lock()
	c.view = view
	c.lock.Unlock()
}
