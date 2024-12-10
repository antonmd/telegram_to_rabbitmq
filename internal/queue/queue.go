package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

// MessagePayload defines the structure of messages we send to RabbitMQ
type MessagePayload struct {
	UpdateID int    `json:"update_id"`
	ChatID   int64  `json:"chat_id"`
	Text     string `json:"text"`
}

// QueueClient now manages a single channel and a return handler.
type QueueClient struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	returnChan chan amqp.Return
	// Known queues map: to avoid redeclaring the same queue multiple times.
	knownQueues map[int64]bool
}

// NewQueueClient establishes the connection and channel, sets up return handling.
// CHANGED: We now create a returnChan and start a goroutine to handle returned messages.
func NewQueueClient(amqpURL string) (*QueueClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	qc := &QueueClient{
		conn:        conn,
		channel:     ch,
		returnChan:  make(chan amqp.Return),
		knownQueues: make(map[int64]bool),
	}

	// CHANGED: Notify the channel about returns and handle them asynchronously.
	qc.channel.NotifyReturn(qc.returnChan)
	go qc.handleReturns()

	return qc, nil
}

// handleReturns listens for returned messages. If a message is returned, it likely means the queue doesn't exist.
// CHANGED: New method to handle returned messages.
func (qc *QueueClient) handleReturns() {
	for ret := range qc.returnChan {
		log.Printf("Message returned by broker: RoutingKey=%s, Reason=%s", ret.RoutingKey, ret.ReplyText)

		// The routing key should match our queue name. Let's extract the chatID from the queue name.
		// We are using the convention: "chat_<chatID>"
		if strings.HasPrefix(ret.RoutingKey, "chat_") {
			chatIDStr := strings.TrimPrefix(ret.RoutingKey, "chat_")
			chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
			if err != nil {
				log.Printf("Failed to parse chatID from returned message: %v", err)
				// We cannot re-publish if we don't know the chat ID.
				continue
			}

			// Declare the queue now since it didn't exist previously.
			err = qc.declareQueueForChat(chatID)
			if err != nil {
				log.Printf("Failed to declare queue for returned message: %v", err)
				// We can't re-publish without a queue.
				continue
			}

			// Now re-publish the returned message.
			// ret.Body contains the original message body.
			// The content type and headers are also available.
			err = qc.channel.Publish(
				"",
				ret.RoutingKey,
				true,  // mandatory still true, in case we fail again
				false, // immediate
				amqp.Publishing{
					ContentType: ret.ContentType,
					Body:        ret.Body,
				},
			)
			if err != nil {
				log.Printf("Failed to re-publish message after queue declaration: %v", err)
			} else {
				log.Printf("Successfully re-published message to queue %s", ret.RoutingKey)
			}
		} else {
			log.Printf("Returned message had unexpected routing key: %s", ret.RoutingKey)
		}
	}
}

// declareQueueForChat is called when we know a chatID doesn't have a queue yet.
// CHANGED: New helper method to declare a queue after a returned message.
func (qc *QueueClient) declareQueueForChat(chatID int64) error {
	if qc.knownQueues[chatID] {
		return nil // already known and declared
	}

	queueName := "chat_" + strconv.FormatInt(chatID, 10)
	q, err := qc.channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue for chat %d: %w", chatID, err)
	}

	qc.knownQueues[chatID] = true
	log.Printf("Declared new queue: %s for chat %d", q.Name, chatID)
	return nil
}

// PublishMessage tries to publish a message to a queue named after the chat ID with mandatory = true.
// CHANGED: We no longer declare the queue beforehand. If the queue does not exist, the message returns and handleReturns() declares it.
func (qc *QueueClient) PublishMessage(msg MessagePayload) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message payload: %w", err)
	}

	queueName := "chat_" + strconv.FormatInt(msg.ChatID, 10)

	// Use mandatory = true so that if the queue doesn't exist, we get a returned message.
	err = qc.channel.Publish(
		"",
		queueName,
		true,  // CHANGED: mandatory = true
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	// If the queue exists, message is routed and no return will happen.
	// If the queue doesn't exist, we'll get the message back in handleReturns().
	return nil
}

func (qc *QueueClient) Close() {
	if qc.channel != nil {
		qc.channel.Close()
	}
	if qc.conn != nil {
		qc.conn.Close()
	}
}
