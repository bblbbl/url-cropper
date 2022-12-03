package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"urls/internal/repo"
	"urls/pkg/etc"
)

type KafkaUrlProducer struct {
	conn *kafka.Conn
}

func NewKafkaUrlProducer(ctx context.Context) (*KafkaUrlProducer, error) {
	cnf := etc.GetConfig()
	conn, err := kafka.DialLeader(
		ctx,
		"tcp",
		fmt.Sprintf("%s:%d", cnf.Kafka.Host, cnf.Kafka.Port),
		urlTopic,
		urlPartition,
	)

	if err != nil {
		return nil, err
	}

	return &KafkaUrlProducer{conn: conn}, nil
}

func (p *KafkaUrlProducer) PutUrlMessage(url repo.Url) (err error) {
	data, err := json.Marshal(url)

	_, err = p.conn.Write(data)

	return err
}

func (p *KafkaUrlProducer) Close() error {
	return p.conn.Close()
}
