package sheduledtasks

import (
	"context"
	"log"
	"time"

	"github.com/sp-va/cryptoCurrencies/internal/repository"
	"github.com/sp-va/cryptoCurrencies/internal/service"
)

func StartCurrencyFetcher(ctx context.Context, repo *repository.PostgresCurrencyRepo) {
	ticker := time.NewTicker(30 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				var repo repository.CurrencyRepository = repo
				coins, err := repo.GetCoinsToTrack()
				if err != nil {
					log.Println("ошибка при получении отслеживаемых койнов:", err)
					continue
				}

				if len(coins) == 0 {
					log.Println("нет койнов")
					continue
				}

				data, err := service.GetCurrencyPrice(coins)
				if err != nil {
					log.Printf("ошибка при получении данных для  %v: %v\n", data, err)
					continue
				}

				for _, coin := range data {

					err = repo.InsertCoinData(coin)
					if err != nil {
						log.Printf("ошибка при сохранении монеты %v: %v\n", coin, err)
					}
				}

			case <-ctx.Done():
				ticker.Stop()
				log.Println("фоновый обработчик остановлен")
				return
			}
		}
	}()
}
