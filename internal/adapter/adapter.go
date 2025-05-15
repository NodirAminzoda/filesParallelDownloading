package adapter

import (
	"net/http"
	"parallel_download_from_many_urls/internal/config"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type IAdapter interface {
	DownloadFile(url, requestID string) error
	GetFileSize(url, requestID string) (int, error)
	LoadPartOfFile(url string, startPartOfBytes, endPartOfBytes int, wg *sync.WaitGroup, data [][]byte, index int, errChan chan error, requestID string)
}

type Adapter struct {
	log    *zerolog.Logger
	cfg    *config.Config
	client *http.Client
}

func New(log *zerolog.Logger, cfg *config.Config) *Adapter {
	client := &http.Client{
		Timeout: cfg.Adapter.TimeOut * time.Second,
	}
	return &Adapter{
		log:    log,
		cfg:    cfg,
		client: client,
	}
}
