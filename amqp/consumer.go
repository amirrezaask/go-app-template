package amqp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func MakeConsumerWithWorkers(
	rabbitMQURI string,
	queueName string,
	routingKey string,
	exchangeName string,
	workerCount int,
	deliveryHandler func(dv amqp091.Delivery) error,
	prefetch int,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var conn *amqp091.Connection
		var delivery <-chan amqp091.Delivery
		var amqpCloseNotifyC chan *amqp091.Error

		restartConnection := func() {
			var err error
			var ch *amqp091.Channel
			fmt.Printf("[Re]Starting %s consumer\n", queueName)
			cfg := amqp091.Config{
				Properties: amqp091.NewConnectionProperties(),
			}

			for range time.NewTicker(time.Second * 1).C {
				conn, err = amqp091.DialConfig(rabbitMQURI, cfg)
				if err != nil {
					slog.Error("cannot connect to rabbit", "err", err)
					continue
				}
				ch, err = conn.Channel()
				if err != nil {
					slog.Error("cannot create rabbit channel", "err", err)
					continue
				}
				amqpCloseNotifyC = make(chan *amqp091.Error, 100)
				conn.NotifyClose(amqpCloseNotifyC)

				delivery, err = ch.ConsumeWithContext(ctx, queueName, "", false, false, false, false, amqp091.Table{})
				if err != nil {
					slog.Error("cannot consume", "queueName", queueName, "err", err)
					continue
				}
				break
			}

		}

		restartConnection()

		workerChans := []chan amqp091.Delivery{}

		for i := 0; i < workerCount; i++ {
			thisWorkerChan := make(chan amqp091.Delivery, 20)
			go func(id int, c chan amqp091.Delivery) {
				for {
					dv, ok := <-c
					if !ok {
						fmt.Printf("Worker-%d is shutting down.\n", i)
						return
					}
					err := deliveryHandler(dv)
					if err != nil {
						slog.Error("cannot process delivery",
							"queueName", queueName,
							"err", err,
						)
						continue
					}

				}

			}(i, thisWorkerChan)
			workerChans = append(workerChans, thisWorkerChan)
		}

		go func() {
			var workerIndex int
			for {
				if workerIndex > workerCount-1 {
					workerIndex = 0
				}
				select {
				case <-amqpCloseNotifyC:
					fmt.Printf("notified of a close event on %s.\n", queueName)
					restartConnection()
				case evt, ok := <-delivery:
					if !ok {
						fmt.Printf("delivery channel is closed for %s.\n", queueName)
						restartConnection()
					}

					go func(evt amqp091.Delivery) {
						workerChans[workerIndex] <- evt

					}(evt)
					workerIndex++
				}
			}
		}()

		return nil

	}

}
