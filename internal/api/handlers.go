package api

import (
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	customErrors "github.com/sp-va/cryptoCurrencies/internal/custom-errors"
	"github.com/sp-va/cryptoCurrencies/internal/dto"
	"github.com/sp-va/cryptoCurrencies/internal/service"
)

type CurrencyHandler struct {
	service *service.CurrencyService
}

func NewCurrencyHandler(service *service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{service: service}
}

func (h *CurrencyHandler) InsertCurrencyHandler(c *gin.Context) {
	var currency dto.AddCurrency
	if err := c.ShouldBindJSON(&currency); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(400, gin.H{"error": "Пустое тело запроса"})
			return
		}
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.service.InsertCurrency(c.Request.Context(), currency)
	if err != nil {
		log.Printf("При добавлении записи возоникла ошибка: %v", err)
		switch err {
		case customErrors.CurrencyAlreadyInserted:
			c.JSON(409, gin.H{"error": err.Error()})
		default:
			c.JSON(500, gin.H{"error": "Ошибка при добавлении валюты"})
		}
		return
	}

	c.JSON(200, gin.H{"status": "Валюта добавлена"})
}

func (h *CurrencyHandler) DeleteCurrencyFromTrackHandler(c *gin.Context) {
	coin := c.Query("coin")
	ra, err := h.service.DeleteCurrencyFromTracking(c.Request.Context(), coin)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при удалении валюты"})
		return
	}

	if ra == 0 {
		c.JSON(404, gin.H{"error": "Не было найдено валюты с таким имене"})
	} else if ra > 0 {
		c.JSON(204, gin.H{})
	}
}

func (h *CurrencyHandler) GetCoinValue(c *gin.Context) {
	coin := c.Query("coin")
	timestamp, _ := strconv.ParseUint(c.Query("timestamp"), 10, 32)

	price, err := h.service.GetCoinValue(c.Request.Context(), coin, uint32(timestamp))

	if err != nil {
		c.JSON(404, gin.H{"error": "не удалось получить цену данной монеты"})
		return
	} else {
		c.JSON(200, gin.H{"result": price.String()})
	}

}

// func (h *CurrencyHandler) GetCurrencyHandler(c *gin.Context) {
// 	coin := c.Query("coin")
// 	timestampStr := c.Query("timestamp")
// 	// конвертация timestampStr в uint32, обработка ошибки

// 	// пример конвертации
// 	var timestamp uint32
// 	// ... парсинг timestampStr ...

// 	currency, err := h.service.GetCurrency(c.Request.Context(), coin, timestamp)
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "Failed to get currency"})
// 		return
// 	}

// 	c.JSON(200, currency)
// }
