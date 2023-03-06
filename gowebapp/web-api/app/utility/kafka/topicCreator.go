package kafka

import (
	log "web-api/app/utility/logger"

	"github.com/Shopify/sarama"
)

func CreateTopic(topic string) error {
	brokers := []string{
		BootstrapServersDefault,
		// "kafka1:9092",
		// "kafka2:9092",
	}
	detail := sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_4_0_0
	admin, err := sarama.NewClusterAdmin(brokers, sarama.NewConfig())
	if err != nil {
		log.ErrorLogger.Println(err)
		return err
	}
	err = admin.CreateTopic(topic, &detail, false)
	if err != nil {
		if err == sarama.ErrTopicAlreadyExists {
			log.InfoLogger.Println(err)
			return nil
		} else {
			log.ErrorLogger.Println("Error from kafka topic admin : ", err)
			return err
		}

	}
	return nil
}
