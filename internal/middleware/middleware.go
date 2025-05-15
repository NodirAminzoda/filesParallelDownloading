package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"parallel_download_from_many_urls/internal/constants"
)

func SetRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Генерация уникального идентификатора для запроса
		requestID := uuid.New().String()

		// Установка идентификатора в заголовок запроса
		c.Set(constants.XRequestID, requestID)

		// Передача управления следующему обработчику
		c.Next()
	}
}
