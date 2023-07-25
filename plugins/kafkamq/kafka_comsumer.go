package kafkamq

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go-live-broadcast-downloader/plugins/log"
	"strings"
)

type KafcInterface interface {
	Process(kafkaTopic string, bs []byte)
}

type KafConsumer struct {
	Addr    []string
	Topics  []string
	GroupId string
	Process KafcInterface
}

func (c *KafConsumer) Consumer() error {
	kfkConsumer, err := doInitConsumer(strings.Join(c.Addr, ","), c.GroupId)
	defer kfkConsumer.Close()
	if err != nil {
		return err
	}
	err = kfkConsumer.SubscribeTopics(c.Topics, nil)
	if err != nil {
		return err
	}

	for {
		msg, err0 := kfkConsumer.ReadMessage(-1)
		if err0 != nil {
			log.Error("ReadMessage").Msgf("%v", err)
			return err0
		}
		//log.Debug("Consumer").Msgf("Message on %s: %s", msg.TopicPartition, string(msg.Value))
		if c.Process != nil {
			c.Process.Process(*msg.TopicPartition.Topic, msg.Value)
		}
	}
}

type KafkaConfig struct {
	Topic            string
	GroupId          string
	BootstrapServers string
}

func doInitConsumer(bootstrapServers, groupId string) (*kafka.Consumer, error) {
	var kafkaConf = &kafka.ConfigMap{
		"api.version.request":       "true",
		"auto.offset.reset":         "latest",
		"heartbeat.interval.ms":     3000,
		"session.timeout.ms":        30000,
		"max.poll.interval.ms":      120000,
		"fetch.max.bytes":           1024000,
		"max.partition.fetch.bytes": 256000,
	}
	_ = kafkaConf.SetKey("bootstrap.servers", bootstrapServers)
	_ = kafkaConf.SetKey("group.id", groupId)
	_ = kafkaConf.SetKey("security.protocol", "plaintext")

	consumer, err := kafka.NewConsumer(kafkaConf)
	if err != nil {
		log.Error("doInitConsumer").Msgf("%v", err)
		return nil, err
	}
	log.Info("doInitConsumer").Msg("kafka 消费者启动成功")
	return consumer, nil
}
