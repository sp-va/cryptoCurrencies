package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sp-va/cryptoCurrencies/internal/dto"
	"github.com/sp-va/cryptoCurrencies/internal/models"
)

func GetCurrencyPrice(coins []string) ([]*models.Currency, error) {
	baseUrl := "https://api.coingecko.com/api/v3/simple/price"

	joinedCoins := strings.Join(coins, ",")

	queryParams := url.Values{}
	queryParams.Add("ids", joinedCoins)
	queryParams.Add("vs_currencies", "usd")

	assembleUrl := fmt.Sprintf("%s?%s", baseUrl, queryParams.Encode())

	response, err := http.Get(assembleUrl)
	if err != nil {
		log.Printf("Ошибка при получении данных: %v", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Ошибка при чтении тела ответа: %v", err)
		return nil, err
	}

	var priceMap map[string]dto.PriceUSD
	err = json.Unmarshal(body, &priceMap)
	if err != nil {
		log.Printf("Ошибка при парсинге JSON: %v", err)
		return nil, err
	}

	var result []*models.Currency
	timestamp := uint32(time.Now().Unix())
	for _, coin := range coins {
		priceData, ok := priceMap[coin]
		if !ok {
			log.Printf("Нет данных для монеты: %s", coin)
			continue
		}

		curr := &models.Currency{
			Coin:      coin,
			Timestamp: timestamp,
			Price:     priceData.Usd,
		}
		result = append(result, curr)
	}

	return result, nil
}
