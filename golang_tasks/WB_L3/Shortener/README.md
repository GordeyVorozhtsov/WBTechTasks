Учебный проект: сервис сокращатель ссылок

Запуск Redis:

1. Скачать легкий образ Redis:
   docker pull redis\:alpine

2. Запустить контейнер с Redis:
   docker run --name redis-test -p 6379:6379 -d redis\:alpine

### Запуск сервера Go:

go run main.go

Сервер будет слушать на localhost:8080.

## Работа с короткими ссылками:

### Создание короткой ссылки:
curl -X POST http://localhost:8080/shorten -H "Content-Type: application/json" -d '{"url":"https://example.com/long/path"}' 

Ответ типа:
{"short\_url":"/s/<ссылка>"}

### Переход по короткой ссылке:
curl -i http://localhost:8080/s/<ссылка> 

## Работа с Redis:

### Просмотр всех ключей:
docker exec -it redis-test redis-cli KEYS '\*'

### Просмотр значения по ключу:
docker exec -it redis-test redis-cli GET <ключ>

## Наполнение аналитики:

### Выполнить несколько переходов по коротким ссылкам:
curl -i http://localhost:8080/s/<ссылка> 

### Просмотр аналитики переходов:
curl http://localhost:8080/analytics/<ссылка>

