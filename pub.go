package main

import (
	"fmt"
	"log"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// optsにClientOptionsインスタンスのpointerを格納
	opts := mqtt.NewClientOptions()
	//　BrokerServerのlistに追加
	opts.AddBroker("tcp://localhost:1883")
	// clientクラスのインスタンスを作成
	c := mqtt.NewClient(opts)
	// BrokerへのconnectionにErrorがないか判定
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Mqtt error: %s", token.Error())
	}

	// clientからosのシステムを利用してパケットをbrokerにstoreする
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("this is msg #%d!", i)
		token := c.Publish("go-mqtt/sample", 0, false, text)
		token.Wait()
	}

	c.Disconnect(250)
	http.ListenAndServe(":8080", nil)
	fmt.Println("Complete publish")
}
