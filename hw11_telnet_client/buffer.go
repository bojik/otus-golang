package main

import (
	"bytes"
	"sync"
)

type Buffer struct {
	bu bytes.Buffer
	mu sync.RWMutex
}

func (c *Buffer) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bu.Write(p)
}

func (c *Buffer) WriteString(s string) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.bu.WriteString(s)
}

func (c *Buffer) Read(p []byte) (int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bu.Read(p)
}

func (c *Buffer) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.bu.String()
}
