# Ссылка на видео с работой сервиса

- https://drive.google.com/file/d/1OFwBk60NHl3nlCQE7KOWl-XFqIIdCnd3/view

# Clen_webservice

Учебный проект от WBTech, реализующий сервис с использованием PostgreSQL, Kafka и внутреннего кеша на Go.

## Описание

Проект демонстрирует работу с Kafka и PostgreSQL, а также кеширование данных в памяти с TTL.  
В проекте отсутствуют интеграционные и юнит-тесты, а также корректный graceful shutdown — это связано с ограничениями по времени.

## Структура проекта

- `producer/producer1.go`, `producer2.go`, `producer3.go` — три продюсера Kafka с разными тестовыми данными для проверки работы сервиса.
- Внутренний кеш реализован на Go с TTL.
- Используется PostgreSQL для хранения заказов.
- Kafka служит для передачи сообщений о заказах.

## Запуск проекта

Для удобства проект можно собрать и запустить с помощью Docker Compose.

### Запуск через Docker Compose

1. Соберите и запустите все сервисы в фоне:

   ```bash
   docker compose up -d 
   ```

2. Проверьте, что все контейнеры запущены:

   ```bash
   docker compose ps
   ```

3. При необходимости можно поднять отдельный контейнер вручную:

   ```bash
   docker compose up -d go-app
   ```

## Тестирование

### Тест Kafka

1. Перейдите в контейнер с приложением:

   ```bash
   docker exec -it go-app /bin/bash
   ```

2. Запустите один из продюсеров для отправки тестовых сообщений в Kafka:

   ```bash
   go run producer/producer1.go
   go run producer/producer2.go
   go run producer/producer3.go
   ```

### Тест записи в базу данных

1. Перейдите в контейнер с PostgreSQL:

   ```bash
   docker exec -it postgres psql -U fucku -d wb
   ```

2. Выполните SQL-запросы для проверки записей в таблицах, например:

   ```sql
   SELECT * FROM orders;
   ```

## Важные замечания

- В проекте не реализованы: интеграционные и юнит тесты, и graceful shutdown из-за ограничений по времени.