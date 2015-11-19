package main

import (
	"crypto/tls"
	"encoding/json"
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/streadway/amqp"
)

func run(config Config) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  config.RethinkDBAddress,
		Database: config.RethinkDBDatabase,
		AuthKey:  config.RethinkDBAuthkey,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})
	if err != nil {
		log.Fatal(session, err)
	}

	rabbitchannel := make(chan []byte, 100)

	go func() {
		cfg := new(tls.Config)
		cfg.InsecureSkipVerify = true
		conn, err := amqp.DialTLS(config.RabbitMQURL, cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Fatal(err)
		}
		defer ch.Close()

		for {
			payload := <-rabbitchannel
			log.Println(string(payload))
			err := ch.Publish("urlwork", "todo", false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(payload),
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

	res, err := r.Table("urls").Changes().Run(session)
	if err != nil {
		log.Fatalln(err)
	}

	var value interface{}

	for res.Next(&value) {
		mapval := value.(map[string]interface{})
		if mapval["new_val"] != nil && mapval["old_val"] == nil {
			jsonbytes, err := json.Marshal(mapval["new_val"])
			if err != nil {
				log.Fatal(err)
			}
			rabbitchannel <- jsonbytes
		}
	}
}
