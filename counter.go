package utils

import (
	"sync"
)

// Map exported
type Map struct {
	sync.Map
	mutex sync.RWMutex
}

// Incr map, 试用读多写少场景
func (m *Map) Incr(key string, value int) {
	m.mutex.Lock()
	if count, ok := m.LoadOrStore(key, value); ok {
		m.Store(key, count.(int)+value)
	}
	m.mutex.Unlock()
}

type Counter struct {
	sync.RWMutex
	Data map[string]int
}

func NewCounter() *Counter {
	return &Counter{Data: make(map[string]int)}
}

func (c *Counter) Get(key string) int {
	c.RLock()
	count := c.Data[key]
	c.RUnlock()
	return count
}

func (c *Counter) Set(key string, value int) {
	c.Lock()
	c.Data[key] = value
	c.Unlock()
}

func (c *Counter) Incr(key string, value int) int {
	c.Lock()
	c.Data[key] += value
	count := c.Data[key]
	c.Unlock()
	return count
}

func (c *Counter) Delete(key string) {
	c.Lock()
	delete(c.Data, key)
	c.Unlock()
}

func (c *Counter) Len() int {
	return len(c.Data)
}
