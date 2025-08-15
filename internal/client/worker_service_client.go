package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"scorer/internal/models"
)

type QueServiceClient struct {
	host        string
	queueClient QueueClientInteface
}

type QueueClientInteface interface {
	Peek(ctx context.Context, queueName string) bool
	Close() error
}

type QueMessage struct {
	Message models.Request
	Score   int
}

func NewQueServiceClient(host string, queueClient QueueClientInteface) QueServiceClient {
	return QueServiceClient{host: host, queueClient: queueClient}
}

func (s QueServiceClient) EnqueueRequest(m models.Request, score int) error {
	//que logic here
	msg := QueMessage{Message: m, Score: score}
	jsonStruct, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = http.Post(fmt.Sprintf("%s/enqueue", s.host), "application/json", bytes.NewBuffer(jsonStruct))

	return err
}

func (s QueServiceClient) PeekQueue(queueName string) bool {
	if s.queueClient != nil {
		return s.queueClient.Peek(context.Background(), queueName)
	}
	return false
}
