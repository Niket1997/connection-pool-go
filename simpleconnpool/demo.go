package simpleconnpool

import (
	"connection-pool-go/simpleconnpool/connpool"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Demo() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	fmt.Println("Hello world!!")

	poolSize := 3

	cp := connpool.New(poolSize)
	defer cp.Close()

	var wg sync.WaitGroup
	numGoRoutines := 6

	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("goroutine %d requesting connection\n", id)
			connection, err := cp.Get()
			if err != nil {
				fmt.Printf("goroutine %d failed to acquire connection: %v\n\n", id, err)
				return
			}

			// Simulate random failure
			if rand.Intn(10) < 3 {
				fmt.Printf("goroutine %d: connection %d failed\n", id, connection.GetId())
				connection.MarkInactive()
			} else {
				fmt.Printf("goroutine %d: using connection %d\n", id, connection.GetId())
				connection.SendRequest()
			}

			// release the connection if it is active
			if connection.IsActive() {
				cp.Release(connection)
			} else {
				fmt.Printf("goroutine %d: connection %d discarded...", id, connection.GetId())
				cp.NewConnection()
			}

		}(i)
	}

	wg.Wait()
	fmt.Println("all goroutines have completed")
}
