package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

var brokers = []string{
	"kafka-broker-1:9091",
	"kafka-broker-2:9092",
	"kafka-broker-3:9093",
}
const topicName = "order_status_changed"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Наш обработчик реализующий интерфейс sarama.ConsumerGroupHandler
	consumerGroupHandler := NewConsumerGroupHandler()

	// Создаем коньюмер группу
	consumerGroup, err := NewConsumerGroup(
		brokers,
		"consumer-group-example",
		[]string{topicName},
		consumerGroupHandler,
	)
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		// запускаем вычитку сообщений
		consumerGroup.Run(ctx)
	}()

	<-consumerGroupHandler.Ready() // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	var (
		consumptionIsPaused = false
		keepRunning         = true
	)
	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(consumerGroup, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()

	if err = consumerGroup.Close(); err != nil {
		log.Fatalf("Error closing consumer group: %v", err)
	}
}

func toggleConsumptionFlow(cg sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		cg.ResumeAll()
		log.Println("Resuming consumption")
	} else {
		cg.PauseAll()
		log.Println("Pausing consumption")
	}

	*isPaused = !*isPaused
}

var _ sarama.ConsumerGroupHandler = (*ConsumerGroupHandler)(nil)

type ConsumerGroupHandler struct {
	ready chan bool
}

func NewConsumerGroupHandler() *ConsumerGroupHandler {
	return &ConsumerGroupHandler{
		ready: make(chan bool),
	}
}

func (h *ConsumerGroupHandler) Ready() <-chan bool {
	return h.ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %v, timestamp = %v, topic = %s",
				message.Value,
				message.Timestamp,
				message.Topic,
			)

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func NewConsumerGroup(brokers []string, groupID string, topics []string, consumerGroupHandler sarama.ConsumerGroupHandler, opts ...Option) (*consumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	config.Consumer.Group.ResetInvalidOffsets = true
	// Сердцебиение консьюмера
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	// Таймаут сессии
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	// Таймаут ребалансировки
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	//
	config.Consumer.Return.Errors = true

	const BalanceStrategy = "roundrobin"
	switch BalanceStrategy {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", BalanceStrategy)
	}

	// Применяем свои конфигурации
	for _, opt := range opts {
		opt.apply(config)
	}

	/*
	  Setup a new Sarama consumer group
	*/
	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &consumerGroup{
		ConsumerGroup: cg,
		handler:       consumerGroupHandler,
		topics:        topics,
	}, nil
}

type consumerGroup struct {
	sarama.ConsumerGroup
	handler sarama.ConsumerGroupHandler
	topics  []string
}

func (c *consumerGroup) Run(ctx context.Context) {
	// Обязательно в случае config.Consumer.Return.Errors = true
	go func() {
		for err := range c.ConsumerGroup.Errors() {
			log.Printf("Error from consumer: %v", err)
		}
	}()

	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.ConsumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
			log.Printf("Error from consume: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
	}
}

type Option interface {
	apply(*sarama.Config) error
}

type optionFn func(*sarama.Config) error

func (fn optionFn) apply(c *sarama.Config) error {
	return fn(c)
}
