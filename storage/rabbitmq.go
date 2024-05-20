package storage

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

type RabbitConsumerConfig struct {
	URI          string
	Exchange     string
	ExchangeType string
	QueueName    string
	BindingKey   string
	ConsumerTag  string
	AutoAck      bool
}

func NewRabbitMQConsumer(cfg RabbitConsumerConfig, messageHandler func(deliveries <-chan amqp.Delivery, done chan error)) (*RabbitConsumer, error) {
	c := &RabbitConsumer{
		conn:    nil,
		channel: nil,
		tag:     cfg.ConsumerTag,
		done:    make(chan error),
	}

	var err error

	config := amqp.Config{Properties: amqp.NewConnectionProperties()}
	config.Properties.SetClientConnectionName("sample-consumer")
	c.conn, err = amqp.DialConfig(cfg.URI, config)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		err := <-c.conn.NotifyClose(make(chan *amqp.Error))
		slog.Error("error in closing rabbitMQ consumer", "err", err)
	}()

	// Log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	// Log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		cfg.Exchange,     // name of the exchange
		cfg.ExchangeType, // type
		true,             // durable
		false,            // delete when complete
		false,            // internal
		false,            // noWait
		nil,              // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	queue, err := c.channel.QueueDeclare(
		cfg.Exchange, // name of the queue
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,     // name of the queue
		cfg.BindingKey, // bindingKey
		cfg.Exchange,   // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	deliveries, err := c.channel.Consume(
		queue.Name,  // name
		c.tag,       // consumerTag,
		cfg.AutoAck, // autoAck
		false,       // exclusive
		false,       // noLocal
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	go messageHandler(deliveries, c.done)

	return c, nil

}

//TODO: Better to write this part based on our needs, it's better to talk with rest of the team to find which way we should write this part. this is the example from rabbitMQ library itself: https://github.com/rabbitmq/amqp091-go/blob/main/_examples/producer/producer.go

type RabbitProducer struct {
}

type RabbitProducerConfig struct {
}

func NewRabbitMQProducer() {}
