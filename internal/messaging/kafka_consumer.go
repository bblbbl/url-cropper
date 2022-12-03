package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"net/url"
	"time"
	"urls/internal/repo"
	"urls/pkg/etc"
	"urls/pkg/utils"
)

type KafkaUrlConsumer struct {
	reader  *kafka.Reader
	ctx     context.Context
	urlChan chan repo.Url
}

func NewKafkaUrlConsumer(ctx context.Context) KafkaUrlConsumer {
	cnf := etc.GetConfig()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{fmt.Sprintf("%s:%d", cnf.Kafka.Host, cnf.Kafka.Port)},
		GroupID:        "url-writer",
		Topic:          urlTopic,
		MinBytes:       512,
		MaxBytes:       512,
		CommitInterval: time.Second,
	})

	return KafkaUrlConsumer{reader: reader, urlChan: make(chan repo.Url), ctx: ctx}
}

func (c *KafkaUrlConsumer) ListenMessages() {
	for {
		select {
		case <-c.ctx.Done():
			close(c.urlChan)
			_ = c.reader.Close()
			break
		default:
			contextTimeout, cancelFunc := context.WithTimeout(c.ctx, 5*time.Second)
			message, err := c.reader.ReadMessage(contextTimeout)
			cancelFunc()

			if err != nil || len(message.Value) == 0 {
				logIfErr(err)
				continue
			}

			unescaped, err := url.QueryUnescape(utils.B2S(message.Value))
			if err != nil {
				etc.GetLogger().Error("failed decode url")
			}

			var u repo.Url
			err = json.Unmarshal([]byte(unescaped), &u)
			if err != nil {
				etc.GetLogger().Info("failed unmarshal url message")
				continue
			}

			c.urlChan <- u
		}
	}

}

func (c *KafkaUrlConsumer) Messages() chan repo.Url {
	return c.urlChan
}

func logIfErr(err error) {
	if err != nil && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
		etc.GetLogger().Info("failed to read message")
	}
}
