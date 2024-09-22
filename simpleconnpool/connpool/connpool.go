package connpool

import (
	"connection-pool-go/simpleconnpool/conn"
	"errors"
	"fmt"
	"sync"
)

// ConnectionPool struct to hold conn pool
type ConnectionPool struct {
	pool chan *conn.Connection
	mu   sync.Mutex
	id   int
}

// New creates a new conn pool
func New(size int) *ConnectionPool {
	cp := &ConnectionPool{
		pool: make(chan *conn.Connection, size),
		id:   0,
	}

	for i := 0; i < size; i++ {
		cp.NewConnection()
	}

	return cp
}

// NewConnection creates a new connection & adds it to connection pool
// if new connection is called async, then mutex helps
// in avoiding creation of connection with same id
func (cp *ConnectionPool) NewConnection() {
	fmt.Println("adding a new connection to the pool")
	cp.mu.Lock()
	defer cp.mu.Unlock()

	connection := conn.New(cp.id)
	cp.id++
	cp.pool <- connection
}

// Get retrieves connection from connection pool
func (cp *ConnectionPool) Get() (*conn.Connection, error) {
	connection := <-cp.pool
	if !connection.IsActive() {
		return nil, errors.New("connection is inactive")
	}
	fmt.Printf("connection %d acquired\n", connection.GetId())
	return connection, nil
}

// Release releases the connection
func (cp *ConnectionPool) Release(connection *conn.Connection) {
	cp.pool <- connection
	fmt.Printf("release connection %d\n", connection.GetId())
}

// Close closes the connection pool
func (cp *ConnectionPool) Close() {
	close(cp.pool)
	fmt.Println("closed connection pool")
}
