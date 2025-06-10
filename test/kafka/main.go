package main

import (
	"log"
	"net/url"

	"github.com/hotstar/golib/kafka"
)

func main() {
	clientId := "test-id"
	endpoint := "localhost:9092"
	u, err := url.Parse(endpoint)
	log.Println(clientId)
	if err != nil {
		log.Fatalf("err:%s", err)
	}
	// p := kafka.NewProducer(kafka.ProducerConfig{
	// 	EndPoint: *u,
	// 	ClientID: "social-sports-client-" + clientId,
	//     }, func(topic string, blob []byte) {

	//     }, func(topic string, blob []byte, err error) {

	//     })

	c := kafka.NewConsumer(kafka.ConsumerConfig{
		EndPoint: *u,
		ClientID: "social-sports-client-" + clientId,
		GroupID:  "matches",
	})

	// go func() {
	topic := "sport-raw-current-match"
	// topic = "sport-raw-concurrency"
	// topic = "finish_task_queue"
	c.Consume(topic, func(c kafka.ConsumerMessage) {
		// fmt.Printf("%s\n", c.Value)
		// m := map[string]string{}
		// json.Unmarshal(c.Value, &m)
		// fmt.Printf("%+v\n", m)
	})
	// }()
	// topic := "sport-raw-detail-match"
	// c.Consume(topic, func(c kafka.ConsumerMessage) {
	// 	fmt.Printf("%s\n", c.Value)
	// })

}
