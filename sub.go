package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	msgCh := make(chan mqtt.Message)

	// importしたmqttのstructで定義されたMessageHandlerに関数を渡す
	var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		// msgChに対してmsgを送信
		msgCh <- msg
	}

	opts := mqtt.NewClientOptions()

	opts.AddBroker("tcp://localhost:1883")

	c := mqtt.NewClient(opts)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	if subscribeToken := c.Subscribe("go-mqtt/sample,", 0, f); subscribeToken.Wait() && subscribeToken.Error() != nil {
		log.Fatal(subscribeToken.Error())
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	for {
		select {
		case m := <-msgCh:
			fmt.Printf("topic: %v, payload: %v\n", m.Topic(), string(m.Payload()))
		case <-signalCh:
			fmt.Printf("Interrupt detected.\n")
			c.Disconnect(1000)
			return
		}
	}
}
