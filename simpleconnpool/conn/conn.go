package conn

import (
	"connection-pool-go/simpleconnpool/constants"
	"fmt"
	"time"
)

// Connection struct to hold conn metadata
type Connection struct {
	id       int
	isActive bool
	// Add other metadata
}

// New creates a new conn
func New(id int) *Connection {
	return &Connection{
		id:       id,
		isActive: true,
	}
}

// SendRequest sends a new request on a conn
func (c *Connection) SendRequest() {
	fmt.Printf("sending request on conn %d\n", c.id)
	time.Sleep(constants.SLEEP_TIME_FOR_CONNECTION)
	fmt.Printf("completed request on conn %d\n", c.id)
}

// GetId return connection id
func (c *Connection) GetId() int {
	return c.id
}

// IsActive check if connection is active
func (c *Connection) IsActive() bool {
	return c.isActive
}

// MarkInactive mark connection as inactive
func (c *Connection) MarkInactive() {
	c.isActive = false
}
