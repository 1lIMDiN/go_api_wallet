# Wallet Service
Это сервис для управления виртуальными кошельками, позволяющее выполнять операции пополнения и снятия средств, а также проверять баланс.

## Основные возможности
- Пополнение счета(Deposit)
- Снятие средств(Withdraw)
- Проверка текущего баланса
- Обработка высоких нагрузок (до 1000 RPS на один кошелек)
- Гарантия целостности данных при конкурентных операциях

## Запуск приложения
### Требования

- Установленный Docker и docker-compose
- Порт 8080 не занят другим приложением
- Порт 5432 не занят (для PostgreSQL)

### Инструкция по запуску

1. Склонируйте репозиторий:
    # git clone <repository-url>
    # cd wallet-service
2. Создайте файл .env на основе config.env:
    # cp config.env .env
3. Запустите сервисы:
    # docker-compose up --build

4. Приложение будет доступно по адресу: http://localhost:8080

### Получение баланса кошелька
    curl -X GET http://localhost:8080/api/v1/wallets/a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11

### Выполнение операции с кошельком
# Пополнение 
curl -X POST -H "Content-Type: application/json" -d '{
  "walletId": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "operationType": "DEPOSIT",
  "amount": 1000
}' http://localhost:8080/api/v1/wallet
# Снятие
curl -X POST -H "Content-Type: application/json" -d '{
  "walletId": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "operationType": "WITHDRAW",
  "amount": 500
}' http://localhost:8080/api/v1/wallet

## Тестирование
Для запуска тестов:

1. Убедитесь, что PostgreSQL запущен и доступен
2. Создайте тестовую БД:
    # createdb -U wallet_user wallet_test
3. Запустите тесты:
    # go test -v ./...