package mysqlconnpool

import (
	mysqlconnpool "connection-pool-go/mysqlconnpool/connpool"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func Demo() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	dsn := "root:password@tcp(127.0.0.1:3306)/demo?parseTime=true"

	poolSize := 3

	cp, err := mysqlconnpool.New(dsn, poolSize)
	if err != nil {
		fmt.Printf("Error creating connection pool: %v\n", err)
		return
	}
	defer cp.Close()

	var wg sync.WaitGroup
	numGoRoutines := 6

	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			conn, err := cp.Get()
			if err != nil {
				fmt.Printf("Goroutine %d: Failed to get connection: %v\n", id, err)
				return
			}
			defer cp.Release(conn)

			// Perform a database operation.
			var now time.Time
			err = conn.QueryRowContext(context.Background(), "SELECT NOW()").Scan(&now)
			if err != nil {
				fmt.Printf("Goroutine %d: Query failed: %v\n", id, err)
				return
			}

			fmt.Printf("Goroutine %d: Current time from DB: %v\n", id, now)
			// Simulate work.
			time.Sleep(2 * time.Second)
		}(i)
	}

	wg.Wait()
	fmt.Println("all goroutines have completed")
}
