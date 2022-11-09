package main

import (
	"flag"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
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
	log.Print("before  connect")
	nc, err := nats.Connect("0.0.0.0:4222")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("before js after connect")
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		fmt.Println(err)
	}
	//
	//js.AddStream(&nats.StreamConfig{
	//	Name:     "FOO",
	//	Subjects: []string{"foo"},
	//})
	//go fetch(js)

	go test(js)
	//js.Publish("foo", []byte("Hello JS!"))

	// ordered push consumer
	//js.Subscribe("foo", func(msg *nats.Msg) {
	//	meta, _ := msg.Metadata()
	//	fmt.Printf("Stream Sequence  : %v\n", meta.Sequence.Stream)
	//	fmt.Printf("Consumer Sequence: %v\n", meta.Sequence.Consumer)
	//
	//}, nats.OrderedConsumer())
	//////////////////////////////////////

	//sub, _ := js.PullSubscribe("foo", "wq", nats.PullMaxWaiting(128))

	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()
	//
	//for {
	//	select {
	//	case <-ctx.Done():
	//		return
	//	default:
	//	}

	// Fetch will return as soon as any message is available rather than wait until the full batch size is available, using a batch size of more than 1 allows for higher throughput when needed.
	//msgs, _ := sub.Fetch(10)
	//for _, msg := range msgs {
	//	fmt.Printf("msg: %s \n", msg.Data)
	//	msg.Ack()
	//}
	//}

	//////////////////////////////////////

	//subj := args[0]

	//nc.Subscribe("foo", func(msg *nats.Msg) {
	//	fmt.Printf("Received a message: %s\n", string(msg.Data))
	//})
	nc.Flush()

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

func fetch(js nats.JetStream) {
	//sub, _ := js.PullSubscribe("ORDERS", "MONITOR")
	//msgs, _ := sub.Fetch(10)
	//fmt.Printf("%+v\n", msgs)
	//js.Subscribe("foo", func(m *nats.Msg) {
	//	fmt.Printf("Received a JetStream message: %s\n", string(m.Data))
	//})
	sub, _ := js.PullSubscribe("foo", "wq", nats.PullMaxWaiting(128))
	msgs, _ := sub.Fetch(10)
	for _, msg := range msgs {
		msg.Ack()
		fmt.Printf("Received a JetStream message: %s\n", string(msg.Data))
	}
}
func test(js nats.JetStream) {
	log.Print("before for")
	sub, _ := js.PullSubscribe("foo", "wq")
	for {

		log.Print("first of for")
		msgs, _ := sub.Fetch(10)
		for _, msg := range msgs {
			msg.Ack()
			fmt.Printf("Received a JetStream message: %s\n", string(msg.Data))
		}
		log.Print("end of for")
	}
}
