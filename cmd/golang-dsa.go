// main.go
// Playground to test DSA and patterns

package main

import (
	"fmt"
	"log"

	"github.com/jdavasligil/golang-dsa/pkg/queue"
)

func main() {
    fmt.Println("Creating a new queue...")
    q := queue.NewQueue[int](5)
    fmt.Println("Done.")

    fmt.Println()

    fmt.Println("Queue empty? ", q.IsEmpty())

    fmt.Println()

    fmt.Println("Enqueue 1..5")
    for i := 1; i < 6; i++ {
        q.Print()
        err := q.Enqueue(i); if err != nil {
            log.Fatalln(err)
        }
    }
    q.Print()

    fmt.Println()

    fmt.Println("Overwrite test")
    for i := 6; i < 12; i++ {
        q.Print()
        q.EnqueueOver(i)
    }
    q.Print()
}
