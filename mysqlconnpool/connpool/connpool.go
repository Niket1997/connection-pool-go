package mysqlconnpool

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectionPool struct to hold conn pool
type ConnectionPool struct {
	pool   chan *sql.Conn
	db     *sql.DB
	mu     sync.Mutex
	closed bool
}

// New creates a new conn pool
func New(dsn string, size int) (*ConnectionPool, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// test db connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	cp := &ConnectionPool{
		pool: make(chan *sql.Conn, size),
		db:   db,
	}

	for i := 0; i < size; i++ {
		if err := cp.NewConnection(); err != nil {
			return nil, err
		}
	}

	return cp, nil
}

// NewConnection creates a new database connection & adds it to connection pool
// if new connection is called async, then mutex helps
// in avoiding creation of connection with same id
func (cp *ConnectionPool) NewConnection() error {
	fmt.Println("adding a new connection to the pool")

	connection, err := cp.db.Conn(context.Background())
	if err != nil {
		return err
	}
	cp.pool <- connection
	return nil
}

// Get retrieves connection from connection pool
func (cp *ConnectionPool) Get() (*sql.Conn, error) {
	cp.mu.Lock()
	if cp.closed {
		cp.mu.Unlock()
		return nil, fmt.Errorf("connection pool is closed")
	}
	cp.mu.Unlock()

	connection := <-cp.pool
	fmt.Println("connection acquired")
	return connection, nil
}

// Release returns a connection back to the pool.
func (cp *ConnectionPool) Release(conn *sql.Conn) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	if cp.closed {
		// Pool is closed; close the connection directly
		err := conn.Close()
		if err != nil {
			return err
		}
		return nil
	}
	cp.pool <- conn
	fmt.Println("Connection released")
	return nil
}

// Close closes all connections in the pool and the underlying database.
func (cp *ConnectionPool) Close() error {
	cp.mu.Lock()
	if cp.closed {
		cp.mu.Unlock()
		return nil // Pool is already closed
	}
	cp.closed = true
	cp.mu.Unlock()

	// Close all connections in the pool
	for {
		select {
		case conn := <-cp.pool:
			err := conn.Close()
			if err != nil {
				return err
			}
		default:
			// No more connections in the pool
			close(cp.pool)
			return cp.db.Close()
		}
	}
}
