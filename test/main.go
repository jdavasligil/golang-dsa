package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	bucket "github.com/jdavasligil/golang-dsa/leaky_bucket"
)

func main() {
	b, err := bucket.NewBucket[string](&bucket.BucketOptions{
		Capacity:        32,
		LowLatency:      false,
		DropBias:        0.90,
		DropInterval:    333 * time.Millisecond,
		MinDropInterval: 100 * time.Millisecond,
		MaxDropInterval: time.Second,
		MaxWaitTime:     5 * time.Second,
		UpdateInterval:  24 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	messages := []string{
		"Best of luck on the runs!",
		"Good morning everyone",
		"what about gnasty gnorc? Im going after HIM",
		"whats the command to see the colored text options?",
		"its been peacefu... indeed",
		"Stewart Copeland was a genius!",
		"What's your sum of best?",
		"glhf",
		":NotLikeThis:",
		"new run new hype",
		"LOL",
		"you just gotta believe!",
		"what do you think of enter the dragonfly or heros tail?",
		"poor frog having to walk up the mountain with a sign every reset",
		"ketchup packets is vegetables",
		"whats up gl",
		" :CatJam: :CatJam: :CatJam: :CatJam: :CatJam: :CatJam:",
		" :CatJam: :CatJam: :CatJam: :CatJam: :CatJam: :CatJam:",
		" :CatJam: :CatJam: :CatJam: :CatJam: :CatJam: :CatJam:",
		"ggs",
		"GG",
		"ggs",
		"GG",
		"GG",
		":GG: :GG: :GG:",
		"I feel fired up Bob!",
		"another run?",
	}

	done := make(chan struct{})

	fmt.Println("Spinning up Producer")
	// Producer
	go func() {
		var idx int

		shutdownTimer := time.NewTimer(60 * time.Second)
		burstTimer := time.NewTicker(6 * time.Second)
		msgTimer := time.NewTicker(time.Second)

		for {
			select {
			case <-msgTimer.C:
				b.AddDrop(messages[idx])
				idx = (idx + 1) % len(messages)
			case <-burstTimer.C:
				for range 12 {
					b.AddDrop(messages[idx])
					idx = (idx + 1) % len(messages)
				}
			case <-shutdownTimer.C:
				log.Println("SHUTDOWN!")
				b.Close()
				log.Println("CLOSED!")
				return
			}
		}
	}()

	fmt.Println("Spinning up Consumer")
	// Consumer
	go func() {
		for {
			drop, err := b.AwaitDrop()
			//fmt.Println(b.Status())

			if errors.Is(err, &bucket.BucketClosedError{}) {
				log.Println("Bucket has closed. Shutting down")
				dropsRemaining := b.Drain()
				log.Printf("Drops Remaining: %v\n", dropsRemaining)
				close(done)
				return
			} else if err != nil {
				log.Println(err)
			} else {
				fmt.Println(drop)
			}
		}
	}()

	<-done
}
