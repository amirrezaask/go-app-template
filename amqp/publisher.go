package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"{{template}}/config"

	"github.com/rabbitmq/amqp091-go"
)

type PublishingPayload struct {
	exchange string
	key      string
	message  []byte
}

type AMQPPublisher struct {
	publishingChannel chan PublishingPayload
}

var (
	Publisher *AMQPPublisher
)

func (a *AMQPPublisher) Publish(exchange string, key string, message []byte) {
	go func() {
		a.publishingChannel <- PublishingPayload{
			exchange: exchange,
			key:      key,
			message:  message,
		}
	}()
}

func NewAMQPPublisher(rabbitMQURI string) *AMQPPublisher {
	var a AMQPPublisher
	var amqpCloseNotifyC chan *amqp091.Error
	var conn *amqp091.Connection

	a.publishingChannel = make(chan PublishingPayload, 100)
	restartConnection := func() {
		fmt.Printf("[Re]Starting publisher connection\n")
		var err error
		cfg := amqp091.Config{
			Properties: amqp091.NewConnectionProperties(),
		}
		conn, err := amqp091.DialConfig(rabbitMQURI, cfg)
		if err != nil {
			slog.Error("cannot connect to rabbit", "err", err)
			for range time.NewTicker(time.Second * 1).C {
				conn, err = amqp091.DialConfig(rabbitMQURI, cfg)
				if err != nil {
					slog.Error("cannot connect to rabbit", "err", err)
					continue
				}
				break
			}

		}
		amqpCloseNotifyC = make(chan *amqp091.Error, 100)
		conn.NotifyClose(amqpCloseNotifyC)
	}

	go func() {
		for {
			select {
			case <-amqpCloseNotifyC:
				fmt.Printf("notified of close connection for amqp publishing.")
				restartConnection()

			case p := <-a.publishingChannel:
				go func() {
					ch, err := conn.Channel()
					if err != nil {
						slog.Error("cannot get amqp channel in handling publishing channel", "err", err)
						time.Sleep(time.Second * 1)
						go a.Publish(p.exchange, p.key, p.message)
						restartConnection()
					}
					defer ch.Close()
					err = ch.PublishWithContext(context.Background(), p.exchange, p.key, false, false, amqp091.Publishing{
						Timestamp: time.Now(),
						Body:      p.message,
					})
					if err != nil {
						slog.Error("error in publishing into amqp", "err", err)
					}
				}()

			}

		}
	}()

	return &a
}

func InitPublisher() {
	Publisher = NewAMQPPublisher(config.RABBITMQ_URI)
}
