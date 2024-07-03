package amqp

import (
	"fmt"
	"log/slog"

	. "github.com/rabbitmq/amqp091-go"
)

type rabbitConnection struct {
	conn *amqp091.Connection
}

var (
	defaultRabbitProducer *rabbitConnection
)

func NewRabbitProducer(rabbitURI string) (*rabbitConnection, error) {
	cfg := amqp091.Config{
		Properties: amqp091.NewConnectionProperties(),
	}

	conn, err := amqp091.DialConfig(rabbitURI, cfg)
	if err != nil {
		slog.Error("cannot connect to rabbit", "err", err)
		return nil, err
	}

	return &rabbitConnection{conn: conn}, nil
}

func init() {
	if testing.Testing() {
		return
	}
	rp, err := NewRabbitProducer(config.RABBITMQ_URI)
	if err != nil {
		panic(err)
	}

	defaultRabbitProducer = rp
}

func (rp *rabbitConnection) PublishContext(ctx context.Context, key string, msg []byte) error {
	ch, err := rp.conn.Channel()
	if err != nil {
		slog.Error("cannot create channel from rabbit mq connection", "err", err)
		return err
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			slog.Error("cannot close channel from rabbit mq connection", "err", err)
		}
	}()

	slog.Debug("got rabbit channel")

	err = ch.PublishWithContext(ctx, config.RABBITMQ_EXCHANGE_NAME, key, false, false, amqp091.Publishing{
		Timestamp: time.Now(),
		Body:      msg,
	})
	if err != nil {
		return err
	}
	slog.Debug("sent message to rabbit", "message", string(msg), "key", key, "exchange", config.RABBITMQ_EXCHANGE_NAME)
	return nil
}
func (rp *rabbitConnection) ConsumeContext(ctx context.Context,
	queueName string,
	routingKey string,
	exchangeName string,
	consumerName string,
) (<-chan amqp091.Delivery, error) {
	ch, err := rp.conn.Channel()
	if err != nil {
		slog.Error("cannot create channel from rabbit mq connection", "err", err)
		return nil, err
	}
	defer func() {
		err := ch.Close()
		if err != nil {
			slog.Error("cannot close channel from rabbit mq connection", "err", err)
		}
	}()

	_, err = ch.QueueDeclare(queueName, true, false, false, false, amqp091.Table{})
	if err != nil {
		slog.Error("cannot declare rabbit queue", "err", err)
		return nil, err
	}

	err = ch.QueueBind(fmt.Sprintf("%s-binding", queueName),
		routingKey, exchangeName, false, amqp091.Table{})
	if err != nil {
		slog.Error("cannot bind rabbit queue", "err", err)
		return nil, err
	}
	delivery, err := ch.ConsumeWithContext(ctx, queueName, consumerName,
		false,
		false,
		false,
		false,
		amqp091.Table{})
	if err != nil {
		slog.Error("cannot consume rabbit queue", "err", err)
		return nil, err
	}

	return delivery, nil
}

// PublishContext publishes given message using given routing key into rabbitMQ. Note that it uses the default singleton connection object.
func PublishContext(ctx context.Context, key string, msg []byte) error {
	if testing.Testing() {
		return nil
	}
	return defaultRabbitProducer.PublishContext(ctx, key, msg)
}

func ConsumeContext(ctx context.Context,
	queueName string,
	routingKey string,
	exchangeName string,
	consumerName string) (<-chan amqp091.Delivery, error) {
	return defaultRabbitProducer.ConsumeContext(ctx, queueName, routingKey, exchangeName, consumerName)
}
