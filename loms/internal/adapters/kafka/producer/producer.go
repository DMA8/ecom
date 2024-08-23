package producer

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
}

func New(brokers []string, opts ...Option) (*Producer, error) {
	config := prepareProducerSaramaConfig(opts...)

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	go func() {
		for err := range asyncProducer.Errors() {
			fmt.Println(err.Error())
		}
	}()

	go func() {
		for msg := range asyncProducer.Successes() {
			fmt.Println("Async success with key", msg.Key)
		}
	}()

	return &Producer{
		asyncProducer: asyncProducer,
	}, nil
}

func NewAsyncProducer(brokers []string, opts ...Option) (sarama.AsyncProducer, error) {
	config := prepareProducerSaramaConfig(opts...)

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, errors.Wrap(err, "error with async kafka-producer")
	}

	go func() {
		for err := range asyncProducer.Errors() {
			fmt.Println(err.Error())
		}
	}()

	go func() {
		for msg := range asyncProducer.Successes() {
			fmt.Println("Async success with key", msg.Key)
		}
	}()

	return asyncProducer, nil
}

func prepareProducerSaramaConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	c.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	c.Producer.RequiredAcks = sarama.WaitForAll
	// c.Producer.Idempotent = true
	c.Producer.Retry.Max = 3
	c.Producer.Retry.Backoff = 5 * time.Millisecond
	c.Producer.Flush.Messages = 10
	c.Producer.Flush.Frequency = time.Second
	c.Net.MaxOpenRequests = 1
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	c.Producer.Compression = sarama.CompressionGZIP

	for _, opt := range opts {
		opt.apply(c)
	}

	return c
}

type Option interface {
	apply(*sarama.Config) error
}

func (p *Producer) Send(ctx context.Context, topicName string, msgBytes []byte, id string) {

	msg, err := buildMessage(ctx, topicName, id, msgBytes)
	if err != nil {
		log.Println(err)
		return
	}

	select {
	case <-ctx.Done():
	case p.asyncProducer.Input() <- msg:
	}
}

func buildMessage(ctx context.Context, topic string, key string, message []byte, headersKV ...string) (*sarama.ProducerMessage, error) {
	if len(headersKV)%2 != 0 {
		return nil, errors.New("wrong number of headersKV")
	}
	headers := make([]sarama.RecordHeader, 0, len(headersKV)/2)

	traceID := ctx.Value(model.TraceIDKey)
	if val, ok := traceID.(string); ok {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte("traceID"),
			Value: []byte(val),
		},
		)
	}
	for i := 0; i < len(headersKV); i += 2 {
		headers = append(headers, sarama.RecordHeader{
			Key:   []byte(headersKV[i]),
			Value: []byte(headersKV[i+1]),
		})
	}

	return &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(message),
		Headers: headers,
	}, nil
}
