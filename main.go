package main

import (
	"fmt"
	"math/rand"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	broker   = "tcp://localhost:1883"
	clientID = "go-mqtt-client"
	topic    = "iot-messages"
)

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to MQTT Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v", err)
}

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		message := generateRandomMessage()
		token := client.Publish(topic, 0, false, message)
		token.Wait()
		fmt.Printf("Published message: %s\n", message)
		time.Sleep(1 * time.Second)
	}
}

func generateRandomMessage() string {
	messages := []string{
		"Hello, World!",
		"Greetings from Go!",
		"MQTT is awesome!",
		"Random message incoming!",
		"Go is fun!",
	}
	return messages[rand.Intn(len(messages))]
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
