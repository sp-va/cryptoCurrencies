# API для получения стоимости криптовалют в USD

## Инструкция по локальному развертыванию
### 1. клонировать репозиторий
```bash
git clone github.com/sp-va/cryptoCurrencies
cd currencies
```

### 2. Создать в корне проекта файл .env со следующими переменными:
```ini
POSTGRES_USER = postgres
POSTGRES_PASSWORD = postgres
POSTGRES_HOST = db
POSTGRES_PORT = 5432
POSTGRES_DB = postgres
```

### 3. Собрать образы и запустить контейнеры:
```bash
docker compose up --build
```

### 4. Взаимодействие с API по адресу: http://localhost:8080

## Доступные эндпоинты:
| Адрес |Метод|Параметры|Тело запроса|Что делает|
| :---: | :---: | :---: | :---: | :---: |
| /api/v1/currency/add | POST | - | JSON:{"coin": "string"} | Добавление койна для его дальнейшего отслеживания |
| /api/v1/currency/remove | DELETE | coin (string, обязательный ) | - | Удаление койна из списка отслеживания (не удаляет при этом накопленные ранее по этому койну записи его цены) |
| /api/v1/currency/price | GET | coin (string, обязательный ), timestamp (int, обязательный  ) | - | Получить ближайшую к данному моменту времени timestamp стоимость монеты coin |
