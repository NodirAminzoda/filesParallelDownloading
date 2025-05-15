package service

import (
	"errors"
	"fmt"
	"os"
	"parallel_download_from_many_urls/internal/constants"
	"sync"

	"parallel_download_from_many_urls/internal/domain/models"
)

func (s *Service) DownloadFiles(urls models.URLs, requestID string) error {
	if len(urls.URL) == 0 {
		return errors.New("empty urls array")
	}

	var (
		mu         sync.Mutex
		failedURLs []string
		wg         sync.WaitGroup
		sem        = make(chan struct{}, 5) // лимит на 5 параллельных загрузок
	)

	for _, url := range urls.URL {
		wg.Add(1)
		sem <- struct{}{} // захватываем "слот"

		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }() // освобождаем "слот"

			err := s.adapter.DownloadFile(url, requestID)
			if err != nil {
				s.log.Warn().Str("url", url).Str("request_id", requestID).
					Err(err).Msg("Не удалось загрузить файл")

				mu.Lock()
				failedURLs = append(failedURLs, url)
				mu.Unlock()
				return
			}

			s.log.Info().Str("url", url).Str("request_id", requestID).
				Msg("Файл успешно загружен")
		}(url)
	}

	wg.Wait()

	if len(failedURLs) > 0 {
		return fmt.Errorf("ошибка загрузки файлов по следующим URL: %v", failedURLs)
	}

	return nil
}

func (s *Service) DownloadFile(url, output, requestID string) error {
	size, err := s.adapter.GetFileSize(url, requestID)
	if err != nil {
		return err
	}

	s.log.Warn().Str(constants.XRequestID, requestID).Msgf("Downloading file with size: %d bytes", size)

	numChunks := (size + constants.ChunkSize - 1) / constants.ChunkSize
	data := make([][]byte, numChunks)
	var wg sync.WaitGroup
	errChan := make(chan error, numChunks)

	for i := 0; i < numChunks; i++ {
		start := i * constants.ChunkSize
		end := start + constants.ChunkSize - 1
		if end >= size {
			end = size - 1
		}

		wg.Add(1)
		go s.adapter.LoadPartOfFile(url, start, end, &wg, data, i, errChan, requestID)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for chanErr := range errChan {
		if chanErr != nil {
			return chanErr
		}
	}

	out, err := os.Create(output)
	if err != nil {
		s.log.Warn().Str(constants.XRequestID, requestID).Msgf("can't create output file: %v", err)
		return err
	}
	defer func() {
		if errClose := out.Close(); errClose != nil {
			s.log.Warn().Str(constants.XRequestID, requestID).Msgf("can't close file: %v", errClose)
		}
	}()

	for _, chunk := range data {
		if _, errFileWrite := out.Write(chunk); errFileWrite != nil {
			return errFileWrite
		}
	}

	s.log.Info().Str(constants.XRequestID, requestID).Msg("File downloaded successfully")
	return nil
}
