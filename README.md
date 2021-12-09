# Тестовое задание avitoTech

<!-- ToC start -->
# Содержание

1. [Описание задачи](#Описание-задачи)
1. [Реализация](#Реализация)
1. [Endpoints](#Endpoints)
1. [Запуск](#Запуск)
1. [Тестирование](#Тестирование)
1. [Примеры](#Примеры)
<!-- ToC end -->

# Описание задачи

Разработать микросервис для работы с балансом пользователей (баланс, зачисление/списание/перевод средств). 
Сервис должен предоставлять HTTP API и принимать/отдавать запросы/ответы в формате JSON.
Дополнительно реализовать методы конвертации баланса и получение списка транзакций.
Полное описание в [TASK](TASK.md).
# Реализация

- Следование дизайну REST API.
- Подход "Чистой Архитектуры" и техника внедрения зависимости.
- Работа с фреймворком [gin-gonic/gin](https://github.com/gin-gonic/gin).
- Работа с СУБД Postgres с использованием библиотеки [sqlx](https://github.com/jmoiron/sqlx) и написанием SQL запросов.
- Конфигурация приложения - библиотека [viper](https://github.com/spf13/viper).
- Запуск из Docker.
- Unit - тестирование уровней бизнес-логики и взаимодействия с БД с помощью моков - библиотеки [testify](https://github.com/stretchr/testify), [mock](https://github.com/golang/mock).

**Структура проекта:**
```
.
├── pkg
│   ├── handler     // обработчики запросов
│   ├── service     // бизнес-логика
│   └── repository  // взаимодействие с БД
├── cmd             // точка входа в приложение
├── schema          // SQL файлы с миграциями
├── configs         // файлы конфигурации
```

# Endpoints

- GET /balance/ - получение баланса пользователя
    - Тело запроса:
        - user_id - уникальный идентификатор пользователя.
  - Параметры запроса:
      - currency - валюта баланса.
- GET /transaction/ - получение транзакций пользователя
    - Тело запроса:
        - user_id - уникальный идентификатор пользователя.
    - Параметры запроса:
        - sort - сортировка списка транзакций.
- POST /top-up/ - пополнение баланса пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя,
        - amount - сумма пополнения в RUB.
- POST /debit/ - списание из баланса пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя,
        - amount - сумма списания в RUB.
- POST /transfer/ - перевод средств на баланс другого пользователя
    - Тело запроса:
        - user_id - идентификатор пользователя, с баланса которого списываются средства,
        - to_id - идентификатор пользователя, на баланс которого начисляются средства,
        - amount - сумма перевода в RUB.
# Запуск

```
make build
make run
```

Если приложение запускается впервые, необходимо применить миграции к базе данных:

```
make migrate-up
```

# Тестирование

Локальный запуск тестов:
```
make run-test
```

# Примеры

Запросы сгенерированы из Postman для cURL.

### 1. GET  /balance для _user_id=1_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/balance' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 1000
}
```

### 2. GET /balance для _user_id=1_ и _currency=USD_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/balance?currency=USD' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 13.542863492536123
}
```

### 3. GET /transaction для _user_id=1_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    },
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    },
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    }
]
```

### 4. GET /transaction для _user_id=1, sort=date_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction?sort=date' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    },
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    }
]
```

### 5. GET /transaction для _user_id=1, sort=amount_

**Запрос:**
```
$ curl --location --request GET 'localhost:8000/transaction?sort=amount' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1
}'
```
**Тело ответа:**
```
[
    {
        "transaction_id": 4,
        "user_id": 1,
        "amount": 10000,
        "operation": "Top-up by bank_card 10000.000000RUB",
        "date": "2021-12-06T13:05:53Z"
    },
    {
        "transaction_id": 7,
        "user_id": 1,
        "amount": 9000,
        "operation": "Debit by purchase 9000.000000RUB",
        "date": "2021-12-06T15:50:15Z"
    },
    {
        "transaction_id": 3,
        "user_id": 1,
        "amount": 100,
        "operation": "Top-up by bank_card 100.000000RUB",
        "date": "2021-12-06T13:05:42Z"
    },
    {
        "transaction_id": 5,
        "user_id": 1,
        "amount": 100,
        "operation": "Debit by transfer 100.000000RUB",
        "date": "2021-12-06T13:06:02Z"
    }
]
```

### 6. POST /top-up для _user_id=1, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/top-up' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 1000
}
```

### 7. POST /debit для _user_id=1, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/debit' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 1,
    "balance": 0
}
```

### 8. POST /transfer для _user_id=1, to_id=2, amount=1000_

**Запрос:**
```
$ curl --location --request POST 'localhost:8000/transfer' \
--header 'Content-Type: application/json' \
--data-raw '{
    "user_id":1,
    "to_id":2,
    "amount":1000
}'
```
**Тело ответа:**
```
{
    "user_id": 2,
    "balance": 1000
}
```
