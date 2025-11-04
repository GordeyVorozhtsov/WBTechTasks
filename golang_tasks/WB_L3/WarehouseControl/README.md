# Учебный проект: сервис склада с историей и ролями

## Запуск проекта

1. Поднять PostgreSQL через Docker Compose:

```
docker-compose up -d
```

2. Запустить Go-сервер:

```
go run main.go
```

3. Открыть браузер на [http://localhost:8080](http://localhost:8080)

На сайте есть селектор ролей (`admin`, `manager`, `viewer`) для имитации входа:

* `admin` — полный доступ
* `manager` — CRUD кроме удаления
* `viewer` — только просмотр(без возможности редактрования и удаления)

---

## Примеры тестовых запросов через curl

### Получение JWT по роли

```bash
curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"role":"admin"}' | jq
```

Ответ:

```bash
{ "token": "...", "role": "admin" }
```

Сохраняем токен:

```bash
export TOKEN="..."
```

---

### CRUD операции

**Создать товар:**

```bash
curl -s -X POST http://localhost:8080/api/items \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"A-1001","name":"Ноутбук Lenovo","quantity":15,"location":"склад-1","note":"новая партия"}' | jq
```

**Список товаров:**

```bash
curl -s -X GET http://localhost:8080/api/items \
  -H "Authorization: Bearer $TOKEN" | jq
```

**Обновить товар (id=1):**

```bash
curl -s -X PUT http://localhost:8080/api/items/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"sku":"A-1001","name":"Ноутбук Lenovo X","quantity":10,"location":"склад-1","note":"частично продано"}' | jq
```

**Удалить товар (только admin):**

```bash
curl -i -X DELETE http://localhost:8080/api/items/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

### Просмотр истории изменений

```bash
curl -s -X GET http://localhost:8080/api/items/1/history \
  -H "Authorization: Bearer $TOKEN" | jq
```

Ответ — массив аудита с полями:

* `operation` — CREATE / UPDATE / DELETE
* `changed_by` — роль, которая сделала изменение
* `changed_at` — время изменения
* `old_data` — данные до изменения
* `new_data` — данные после изменения

---

## Ограничения по ролям

* `admin` — полный доступ
* `manager` — CRUD кроме удаления
* `viewer` — только просмотр
