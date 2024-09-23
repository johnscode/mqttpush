package main

import (
	"encoding/json"
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

type Message struct {
	Time   time.Time `json:"time"`
	Device Device    `json:"device"`
}

type Device interface {
	ID() string
	Name() string
}

type TempRHDevice struct {
	Id         string  `json:"id"`
	DeviceName string  `json:"name,omitempty"`
	Temp       float32 `json:"temp,omitempty"`
	Rh         float32 `json:"rh,omitempty"`
}

func (t TempRHDevice) ID() string {
	return t.Id
}

func (t TempRHDevice) Name() string {
	return t.DeviceName
}

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
		payload, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			panic(err)
		}
		token := client.Publish(topic, 0, false, payload)
		token.Wait()
		fmt.Printf("Published json message: %s\n", payload)
		time.Sleep(1 * time.Second)
	}
}

func generateRandomMessage() Message {
	msg := Message{
		Time: time.Now(),
		Device: TempRHDevice{
			Id:         "043e5af81c",
			DeviceName: "Greenhouse",
			Temp:       76.3 + rand.Float32()*2.5,
			Rh:         52.9 + rand.Float32()*1.3,
		},
	}
	return msg
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
