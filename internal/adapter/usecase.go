package adapter

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
	"parallel_download_from_many_urls/internal/constants"
	"path/filepath"
	"strconv"
	"sync"
)

func (a *Adapter) DownloadFile(url, requestID string) error {
	if url == "" {
		return errors.New("empty URL provided")
	}

	dir := "downloads"
	if err := os.MkdirAll(dir, 0755); err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при создании директории: %v", err)
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	fileID := uuid.New().String() // создание уникальное название для каждого загружаемого файла

	outputPath := filepath.Join(dir, fmt.Sprintf("%s.zip", fileID))

	resp, err := a.client.Get(url)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при выполнении GET-запроса: %v", err)
		return fmt.Errorf("failed to download from %s: %w", url, err)
	}
	defer func() {
		if errResponseBodyClose := resp.Body.Close(); errResponseBodyClose != nil {
			a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при закрытии тела ответа: %v", errResponseBodyClose)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("Получен неуспешный HTTP-статус: %d %s", resp.StatusCode, resp.Status)
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при создании файла: %v", err)
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}

	defer func() {
		if errFileClose := outFile.Close(); errFileClose != nil {
			a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при закрытии файла: %v", errFileClose)
		}
	}()

	if _, err = io.Copy(outFile, resp.Body); err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("Ошибка при сохранении содержимого: %v", err)
		return fmt.Errorf("failed to write response to file: %w", err)
	}

	a.log.Info().Str(constants.XRequestID, requestID).Msgf("Файл успешно загружен: %s", outputPath)
	return nil
}

func (a *Adapter) GetFileSize(url, requestID string) (int, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't create new request for getting file size, err:%v", err)
		return 0, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't send request, error: %v", err)
		return 0, err
	}

	defer func() {
		if errCloseResponseBody := resp.Body.Close(); errCloseResponseBody != nil {
			a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't close response body, error: %v", errCloseResponseBody)
		}
	}()

	sizeStr := resp.Header.Get("Content-Length")

	if sizeStr == "" {
		return 0, errors.New("Content-Length not found")
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't convertation string to int, size: %v, error: %v", sizeStr, err)
		return 0, err
	}

	return size, nil
}

func (a *Adapter) LoadPartOfFile(url string, startPartOfBytes, endPartOfBytes int, wg *sync.WaitGroup, data [][]byte, index int, errChan chan error, requestID string) {
	defer wg.Done()

	// подготовка запроса для получении данных
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't create new request for getting file size, err:%v", err)
		return
	}
	rangeHeader := fmt.Sprintf("bytes=%d-%d", startPartOfBytes, endPartOfBytes)

	// добавление заголовка который указывает диапозон получения данных
	req.Header.Set("Range", rangeHeader)

	resp, err := a.client.Do(req)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't send request, error: %v", err)
		errChan <- err
		return
	}

	defer func() {
		if errCloseResponseBody := resp.Body.Close(); errCloseResponseBody != nil {
			a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't close response body, error: %v", errCloseResponseBody)
			errChan <- errCloseResponseBody
		}
	}()

	chunk, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error().Str(constants.XRequestID, requestID).Msgf("can't read response body, error: %v", err)
		errChan <- err
		return
	}

	// добавляем данные в массив для получение правильного получение данных после получения
	data[index] = chunk
}
