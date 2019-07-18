package socket_server

import (
	"sync"
)

type ConnMgr struct {
	connections map[string]*Connection
	mutex       *sync.RWMutex
}

func NewConnMgr() *ConnMgr {
	return &ConnMgr{
		connections: make(map[string]*Connection),
		mutex:       new(sync.RWMutex),
	}
}

func (cm *ConnMgr) Put(id string, c *Connection) {
	cm.mutex.Lock()
	cm.connections[id] = c
	cm.mutex.Unlock()
}

func (cm *ConnMgr) Get(id string) *Connection {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	if c, ok := cm.connections[id]; ok {
		return c
	}

	return nil
}

func (cm *ConnMgr) Del(id string) {
	cm.mutex.Lock()
	delete(cm.connections, id)
	cm.mutex.Unlock()
}
