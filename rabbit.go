package fsq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type RabbitConfig struct {
	RabbitHost string
	RabbitPort int
}

type RabbitQueue struct {
	cfg    *RabbitConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
	queue  string
	sender ISender
}

func NewRabbitQueue(cfg *RabbitConfig, sender ISender) (*RabbitQueue, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%v/", cfg.RabbitHost, cfg.RabbitPort))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq on host: %v and port: %v: %v", cfg.RabbitHost, cfg.RabbitPort, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queueName := "email_queue"
	_, err = ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitQueue{
		cfg:    cfg,
		conn:   conn,
		ch:     ch,
		queue:  queueName,
		sender: sender,
	}, nil
}

func (r *RabbitQueue) SendToQueue(to string, subject string, body string) error {
	email := Email{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	b, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("failed to marshal email: %v", err)
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	if err = r.ch.Publish("", r.queue, false, false, msg); err != nil {
		return fmt.Errorf("failed to publish message to queue: %v", err)
	}

	return nil
}

func (r *RabbitQueue) Consume(ctx context.Context) error {
	messages, err := r.ch.Consume(r.queue, "queue_sender", true, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-messages:
				if !ok {
					return
				}

				var email Email
				err := json.Unmarshal(msg.Body, &email)
				if err != nil {
					continue
				}

				if err = r.sender.SendMail(email.To, email.Subject, email.Body); err != nil {
					continue
				}
			}
		}
	}()

	return nil
}
