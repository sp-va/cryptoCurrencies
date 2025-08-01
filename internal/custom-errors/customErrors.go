package customErrors

import "errors"

var CurrencyAlreadyInserted = errors.New("Эта валюта уже добавлена в базу")
