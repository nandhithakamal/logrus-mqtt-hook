package logrus_mqtt_hook

import (
	paho "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
	"time"
)



type MqttHook struct {
	client    paho.Client
	levels    []log.Level
	topic     string
	qos       int
	retain    bool
	Formatter log.Formatter
}

const TimeoutInterval = time.Minute

func NewMqttHook(topic string, mqttClient paho.Client, levels []log.Level, logFormatter log.Formatter) *MqttHook {
	//log.Info("Creating new mqtt hook...")
	if logFormatter == nil {
		//log.Info("logformatter is nil")
		logFormatter = &log.JSONFormatter{}
	}
	//log.Info("Creating mqtt client...")
	mqttClient = mqttClient
	//log.Info("Created...")
	return &MqttHook{
		client:    mqttClient,
		topic:     topic,
		qos:       0,
		retain:    false,
		levels:    levels,
		Formatter: logFormatter,
	}
}


func (h *MqttHook) Fire(entry *log.Entry) error {
	//log.Info("Firing...")
	logEntry, err := h.Formatter.Format(entry)
	if err != nil {
		return err
	}

	token := h.client.Publish(h.topic, 0, false, logEntry)
	success := token.WaitTimeout(TimeoutInterval)
	err = token.Error()
	if !success || err != nil {
		log.Errorf("Error while publishing MQTT message, timedout: %v Error:%s ", !success, err.Error())
	}
	return err
}

func (h *MqttHook) Levels() []log.Level {
	//log.Info("Finding levels...")
	return h.levels
}
