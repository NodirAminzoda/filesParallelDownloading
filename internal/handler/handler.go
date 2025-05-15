package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"parallel_download_from_many_urls/internal/constants"
	"parallel_download_from_many_urls/internal/domain/models"
	"parallel_download_from_many_urls/utils"
)

func (h *Handler) DownloadFiles(c *gin.Context) {
	requestID := utils.GetRequestIDFromContext(c)

	var urlArr models.URLs
	if err := c.ShouldBindJSON(&urlArr); err != nil {
		h.log.Error().Str(constants.XRequestID, requestID).Msgf("can't parse incoming data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Str(constants.XRequestID, requestID).Msgf("Request body: %+v", urlArr)

	if err := h.service.DownloadFiles(urlArr, requestID); err != nil {
		h.log.Error().Str(constants.XRequestID, requestID).Msgf("can't download some files: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Str(constants.XRequestID, requestID).Msg("Files downloaded successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Данные успешно загружены"})
}

func (h *Handler) DownloadFile(c *gin.Context) {
	requestID := utils.GetRequestIDFromContext(c)

	var input models.FileURL // лучше назвать структуру осмысленно
	if err := c.ShouldBindJSON(&input); err != nil {
		h.log.Error().Str(constants.XRequestID, requestID).Msgf("can't parse incoming data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Str(constants.XRequestID, requestID).Msgf("Request body: %v", input)

	err := h.service.DownloadFile(input.Url, constants.OutputFilePath, requestID)
	if err != nil {
		h.log.Error().Str(constants.XRequestID, requestID).Msgf("can't download file, error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Str(constants.XRequestID, requestID).Msg("File successfully downloaded")

	c.JSON(http.StatusOK, gin.H{"message": "file downloaded successfully"})
}
