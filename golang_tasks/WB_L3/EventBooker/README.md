# Учебный проект: сервис бронирования

# Запустить PostgreSQL:
docker run -d --name EBoo \
  -e POSTGRES_DB=eventbooker \
  -e POSTGRES_USER=eventuser \
  -e POSTGRES_PASSWORD=eventpass \
  -p 5433:5432 \
  postgres:13

# Запустить сервер 
* go run main.go

## Проверка через браузер

Можно открыть интерфейс в браузере по адресу:  
[http://localhost:8080](http://localhost:8080)