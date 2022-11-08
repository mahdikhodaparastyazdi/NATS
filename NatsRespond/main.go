package main

import (
	"flag"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
	"runtime"
	"time"
	//"fmt"
)

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'", i, m.Subject, string(m.Data))
}

func main() {

	args := flag.Args()
	if len(args) != 1 {

	}

	opts := []nats.Option{nats.Name("NATS Sample Subscriber")}
	opts = setupConnOptions(opts)

	nc, err := nats.Connect("0.0.0.0:4222")
	if err != nil {
		log.Fatal(err)
	}
	//subj := args[0]
	i := 0
	nc.QueueSubscribe("foo", "mahdi", func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		msg.Respond([]byte("ok"))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", "foo")

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	nc.Drain()
	log.Fatalf("Exiting")

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]", "foo")

	runtime.Goexit()
}

//for {
//	nc.Subscribe("help", func(m *nats.Msg) {
//		nc.Publish(m.Reply, []byte("I can help!"))
//	})
//	// nc.Subscribe("foo", func(m *nats.Msg) {
//	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
//	// })
//
//}
//}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}
