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
	DeviceID() string
	Name() string
	Type() string
}

type TempRHDevice struct {
	DeviceId   string  `json:"device_id"`
	DeviceName string  `json:"name,omitempty"`
	DeviceType string  `json:"device_type"`
	Temp       float32 `json:"temp,omitempty"`
	Rh         float32 `json:"rh,omitempty"`
}

func (t TempRHDevice) DeviceID() string {
	return t.DeviceId
}

func (t TempRHDevice) Name() string {
	return t.DeviceName
}

func (t TempRHDevice) Type() string {
	return t.DeviceType
}

func NewTempRHDevice(id string, name string) Device {
	return &TempRHDevice{
		DeviceId:   id,
		DeviceName: name,
		DeviceType: "TempRH",
	}
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
			DeviceId:   "043e5af81c",
			DeviceName: "Greenhouse",
			DeviceType: "TempRH",
			Temp:       76.3 + rand.Float32()*2.5,
			Rh:         52.9 + rand.Float32()*1.3,
		},
	}
	return msg
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}
