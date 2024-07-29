package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

var totalCount = 0

func TestAddRateLimiting(t *testing.T) {
	var wg sync.WaitGroup
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go makeCall(i, &wg)
	}
	wg.Wait()
	if totalCount > MaxConnections {
		t.Fatalf("More requests executed than max allowed - %d, Max allowed - %d", totalCount, MaxConnections)
	}
	fmt.Printf("Total success count - %d", totalCount)
}

func makeCall(count int, wg *sync.WaitGroup) {
	res, err := http.Get("http://localhost:8080/v1/async")
	if err != nil {
		fmt.Printf("error in making call %e", err)
		wg.Done()
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Failed call number %d with status code %d\n", count, res.StatusCode)
		wg.Done()
		return
	}
	fmt.Printf("Succeeded call number %d\n", count)
	totalCount++
	wg.Done()
}
