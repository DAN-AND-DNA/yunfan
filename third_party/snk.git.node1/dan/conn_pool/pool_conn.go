package conn_pool

import (
	"net"
	"sync"
)

type Pool_conn struct {
	is_init bool
	net.Conn
	owner_pool Conn_pool
	is_bad     bool
	m_conn     sync.RWMutex
}

func (this *Pool_conn) Close() error {
	if !this.is_init {
		return nil
	}

	this.m_conn.RLock()
	defer this.m_conn.RUnlock()

	if this.is_bad {
		// discard
		this.Conn.Close()
		return nil
	}

	this.owner_pool.put(this.Conn)
	return nil
}

func (this *Pool_conn) Set_bad() {
	if !this.is_init {
		return
	}

	this.m_conn.Lock()
	defer this.m_conn.Unlock()
	this.is_bad = true
}
