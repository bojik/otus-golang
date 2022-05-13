package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/bojik/otus-golang/hw12_13_14_15_calendar/internal/logger"
	backoff "github.com/cenkalti/backoff/v4"
	"github.com/streadway/amqp"
)

type Queue struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	queue        amqp.Queue
	amqpURI      string
	exchangeName string
	queueName    string
	exchangeType string
	done         chan error
	l            logger.Logger
}

func New(amqp, exchangeName, exchangeType, queueName string, l logger.Logger) *Queue {
	return &Queue{
		l:            l,
		amqpURI:      amqp,
		queueName:    queueName,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
	}
}

func (q *Queue) Connect() error {
	var err error

	q.conn, err = amqp.Dial(q.amqpURI)
	if err != nil {
		return fmt.Errorf("amqp dial: %w", err)
	}

	q.channel, err = q.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}

	go func() {
		q.l.Error(fmt.Sprintf("closing: %s", <-q.conn.NotifyClose(make(chan *amqp.Error))))
		q.done <- errors.New("channel Closed")
	}()

	if err = q.channel.ExchangeDeclare(
		q.exchangeName,
		q.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}

	q.queue, err = q.channel.QueueDeclare(
		q.queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if err = q.channel.QueueBind(
		q.queue.Name,   // name of the queue
		"",             // bindingKey
		q.exchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return fmt.Errorf("queue bind: %w", err)
	}

	return nil
}

func (q *Queue) Publish(ctx context.Context, d time.Duration, a *app.App) error {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			q.l.Debug("deleting old events")
			if err := a.DeleteOldEvents(ctx); err != nil {
				return err
			}
			q.l.Debug("finding events to sent")
			events, err := a.FindToSend(ctx)
			if err != nil {
				return err
			}
			q.l.Debug(fmt.Sprintf("found %d events", len(events)))
			for _, event := range events {
				if err := q.publish(event); err != nil {
					return err
				}
				if err := a.MarkAsSent(ctx, event.ID); err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (q *Queue) Consume(ctx context.Context, threads int) error {
	msgs, err := q.announceQueue()
	if err != nil {
		return fmt.Errorf("queue Consume: %w", err)
	}
	for {
		for i := 0; i < threads; i++ {
			go q.worker(ctx, msgs)
		}

		if <-q.done != nil {
			msgs, err = q.reConnect(ctx)
			if err != nil {
				return fmt.Errorf("reconnecting Error: %w", err)
			}
		}
		q.l.Info("reconnected... possibly")
	}
}

func (q *Queue) announceQueue() (<-chan amqp.Delivery, error) {
	msgs, err := q.channel.Consume(
		q.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	return msgs, err
}

func (q *Queue) worker(ctx context.Context, msgs <-chan amqp.Delivery) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-q.done:
			return
		case msg := <-msgs:
			event := &app.Event{}
			if err := json.Unmarshal(msg.Body, event); err != nil {
				q.l.Error(
					"unmarshal error: "+err.Error(),
					logger.NewStringParam("json", string(msg.Body)),
				)
				if err := msg.Reject(false); err != nil {
					q.l.Error(
						"reject error: "+err.Error(),
						logger.NewStringParam("json", string(msg.Body)),
					)
				}
				continue
			}
			q.l.Info("receiving message: "+event.ID, logger.NewStringParam("title", event.Title))
		}
	}
}

func (q *Queue) Close() error {
	if err := q.channel.Close(); err != nil {
		return fmt.Errorf("amqp channel closing: %w", err)
	}
	if err := q.conn.Close(); err != nil {
		return fmt.Errorf("amqp closing: %w", err)
	}
	return nil
}

func (q *Queue) publish(event *app.Event) error {
	encoded, err := json.Marshal(event)
	if err != nil {
		return err
	}
	q.l.Info(fmt.Sprintf("publishing: %s", event.ID))
	if err := q.channel.Publish(
		q.exchangeName, // publish to an exchange
		"",             // routing to 0 or more queues
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/javascript",
			ContentEncoding: "",
			Body:            encoded,
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("exchange Publish: %w", err)
	}
	return nil
}

func (q *Queue) reConnect(ctx context.Context) (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-time.After(d):
			if err := q.Connect(); err != nil {
				q.l.Error("could not connect in reconnect call: " + err.Error())
				continue
			}
			msgs, err := q.announceQueue()
			if err != nil {
				fmt.Printf("Couldn't connect: %+v", err)
				continue
			}

			return msgs, nil
		case <-ctx.Done():
			return nil, nil
		}
	}
}
