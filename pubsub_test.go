package main

import (
	"fmt"
	"github.com/mhconradt/pingpong-conversations-api/pubsub"
	"log"
)

func sub() {
	s, err := pubsub.NewSubscriber("my-topic", "sub-1")
	if err != nil {
		log.Fatal(err)
	}
	for msg := range s.Messages {
		m := string(msg)
		fmt.Println(m)
		if m == "exit" {
			s.Close()
		}
	}
}

func pub() {
	p, err := pubsub.NewProducer("my-topic")
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte("oh no")
	err = p.Send(msg)
	if err != nil {
		fmt.Println("error sending oh no: ", err)
	}
}

/*
func main() {
	if len(os.Args) < 2 {
		log.Fatal("need to specify publish (pub) or subscribe (sub)")
	}

}
*/
