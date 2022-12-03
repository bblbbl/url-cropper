package messaging

import "urls/internal/repo"

const (
	urlTopic     = "urls"
	urlPartition = 0
)

type UrlProducer interface {
	PutUrlMessage(url repo.Url) error
	Close() error
}

type UrlConsumer interface {
	ListenMessages()
	Messages() chan repo.Url
}
