package main

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	. "github.com/nandhithakamal/logrus_mqtt_hook"
	log "github.com/sirupsen/logrus"
	"time"
)


func main() {
	opts := paho.NewClientOptions().AddBroker("tcp://localhost:1883")
	mqttClient := paho.NewClient(opts)

	token := mqttClient.Connect()
	success := token.WaitTimeout(time.Minute)
	err := token.Error()
	if !success || err != nil {
		log.Errorf("Error while connecting, timedout: %v Error:%s ", !success, err.Error())
	}

	token = mqttClient.Publish("test", 0, false, "hello paho")
	success = token.WaitTimeout(time.Minute)
	err = token.Error()
	if !success || err != nil {
		log.Errorf("Error while publishing MQTT message, timedout: %v Error:%s ", !success, err.Error())
	}
	mqttHook := NewMqttHook("test", mqttClient, log.AllLevels, &log.JSONFormatter{})
	log.AddHook(mqttHook)

	log.Info("Info message")
	log.Error("Error message")

	log.StandardLogger().ReplaceHooks(map[log.Level][]log.Hook{})
	mqttClient.Disconnect(0)
}
