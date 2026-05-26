# Wallet API

REST-сервис для управления кошельками на Go + PostgreSQL

## Запуск

1. Создай `config.env` на основе примера:

```bash
cp config.env.example config.env
```

2. Заполни `config.env` своими данными

3. Запусти:

```bash
docker-compose up --build
```

Поднимает PostgreSQL, применяет миграции и запускает сервер на порту `9091`

## API

### Создать кошелёк

```
POST /api/v1/wallet/create
```

Ответ:
```json
{
    "wallet_id": "9f3580b6-bc4b-4ab7-87d1-9ccc2c331c0b",
    "balance": 0
}
```

---

### Пополнить или списать средства

```
POST /api/v1/wallet
Content-Type: application/json

{
    "walletId": "9f3580b6-bc4b-4ab7-87d1-9ccc2c331c0b",
    "operationType": "DEPOSIT",
    "amount": 100
}
```

`operationType`: `DEPOSIT` — пополнение, `WITHDRAW` — списание

Ответ:
```json
{
    "wallet_id": "9f3580b6-bc4b-4ab7-87d1-9ccc2c331c0b",
    "balance": 100
}
```

---

### Получить баланс

```
GET /api/v1/wallets/{WALLET_UUID}
```

Ответ:
```json
{
    "wallet_id": "9f3580b6-bc4b-4ab7-87d1-9ccc2c331c0b",
    "balance": 100
}
```

## Стек

- Golang
- PostgreSQL
- Docker / Docker Compose