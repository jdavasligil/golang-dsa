package bucket

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func TestLeakyBucket() {
	b := NewBucket[string](
	32,
	1,
	200 * time.Millisecond,
	time.Second,
	100 * time.Millisecond,
	1.05,
	2 * time.Second,
	5 * time.Second,
	)
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

		shutdownTimer := time.NewTimer(8 * time.Second)
		burstTimer := time.NewTicker(4 * time.Second)
		msgTimer := time.NewTicker(time.Second)

		for {
			select {
			case <-msgTimer.C:
				b.AddDrop(messages[idx])
				idx = (idx + 1) % len(messages)
			case <-burstTimer.C:
				for range 6 {
					b.AddDrop(messages[idx])
					idx = (idx + 1) % len(messages)
				}
			case <-shutdownTimer.C:
				err := b.Close()
				if err != nil {
					log.Fatal("Bucket already closed!")
				}
				return
			}
		}
	}()

	fmt.Println("Spinning up Consumer")
	// Consumer
	go func() {
		for {
			drop, err := b.AwaitDrop()
			fmt.Println(b.Status())

			if errors.Is(err, &BucketClosedError{}) {
				log.Println("Bucket has closed. Shutting down")
				dropsRemaining := b.Drain()
				log.Printf("Drops Remaining: %v\n", dropsRemaining )
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
