package service

import (
	"github.com/rs/zerolog"
	"parallel_download_from_many_urls/internal/adapter"
	"parallel_download_from_many_urls/internal/config"
	"parallel_download_from_many_urls/internal/domain/models"
)

type IService interface {
	DownloadFiles(urls models.URLs, requestID string) error
	DownloadFile(url, output, requestID string) error
}

type Service struct {
	cfg     *config.Config
	log     *zerolog.Logger
	adapter adapter.IAdapter
}

func New(cfg *config.Config, log *zerolog.Logger, adapter adapter.IAdapter) *Service {
	return &Service{
		cfg:     cfg,
		log:     log,
		adapter: adapter,
	}
}
