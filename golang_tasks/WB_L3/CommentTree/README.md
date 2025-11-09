# Учебный проект: сервис древовидных комментариев

curl -X POST http://localhost:8080/comments \
-H "Content-Type: application/json" \
-d '{"content":"Первый комментарий"}'

curl http://localhost:8080/comments

curl http://localhost:8080/comments?parent=1

curl -X DELETE http://localhost:8080/comments/1


запустить  сервер go run main.go
открыть http://localhost:8080/static/ 
писать комменты на сайте и смотреть, благо есть юайка и терминалл для теста не понадобится




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