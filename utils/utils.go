package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"parallel_download_from_many_urls/internal/constants"
)

func GetRequestIDFromContext(c *gin.Context) string {
	// Пытаемся получить XRequestID из контекста
	XRequestID, ok := c.Get(constants.XRequestID)
	if !ok {
		// Если значение отсутствует, генерируем новый UUID
		return uuid.New().String()
	}

	// Приводим значение к строке
	requestID, ok := XRequestID.(string)
	if !ok {
		// Если значение не является строкой, генерируем новый UUID
		return uuid.New().String()
	}

	// Возвращаем корректный requestID
	return requestID
}
