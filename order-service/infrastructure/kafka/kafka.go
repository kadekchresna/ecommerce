package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/kadekchresna/ecommerce/order-service/helper/logger"
	"github.com/kadekchresna/ecommerce/order-service/infrastructure/messaging"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

var _ messaging.Producer = (*KafkaProducer)(nil)

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      brokers,
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: int(kafka.RequireAll),
		}),
	}
}

func (kp *KafkaProducer) Publish(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}

	err := kp.writer.WriteMessages(ctx, msg)
	if err != nil {
		err = fmt.Errorf("error publish message :: KafkaProducer-Publish() failed to publish message: %v", err)
		logger.LogWithContext(ctx).Error(err.Error())
		return err
	}

	logger.LogWithContext(ctx).Info(fmt.Sprintf("KafkaProducer-Publish() :: published message with key: %s", string(key)))
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}

type KafkaConsumer struct {
	reader *kafka.Reader
}

var _ messaging.Consumer = (*KafkaConsumer)(nil)

func NewKafkaConsumer(brokers []string, topic, groupID string) *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})

	return &KafkaConsumer{
		reader: r,
	}
}

func (kc *KafkaConsumer) ReadMessage(ctx context.Context) (messaging.Message, error) {
	msg, err := kc.reader.FetchMessage(ctx)
	if err != nil {
		return messaging.Message{}, err
	}

	return messaging.Message{
		Key:       msg.Key,
		Value:     msg.Value,
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	}, nil
}

func (kc *KafkaConsumer) Close() error {
	return kc.reader.Close()
}

func (kc *KafkaConsumer) Commit(ctx context.Context, msg messaging.Message) error {
	return kc.reader.CommitMessages(ctx, kafka.Message{
		Topic:     msg.Topic,
		Partition: msg.Partition,
		Offset:    msg.Offset,
	})
}
