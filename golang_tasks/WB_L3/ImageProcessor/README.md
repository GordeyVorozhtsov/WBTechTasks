# Учебный проект: сервис обработки изображений


## Запуск Kafka в контейнере

```bash
docker rm -f kafka
docker run -d --name kafka \
  -p 9092:9092 \
  -e KAFKA_KRAFT_MODE=true \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_OFFSETS_TOPIC_NUM_PARTITIONS=1 \
  -e KAFKA_LOG_DIRS=/tmp/kraft-combined-logs \
  -e CLUSTER_ID=$(uuidgen) \
  apache/kafka:4.1.0
```


## Создание Kafka-топика

```bash
docker exec -it kafka /opt/kafka/bin/kafka-topics.sh \
 --create --topic image_tasks \
 --bootstrap-server localhost:9092 \
 --partitions 1 \
 --replication-factor 1
```


## Загрузка изображения

Отправляем POST-запрос с изображением
Проверяется наличие оригинала в `storage/original` и обработанного файла в `storage/processed`.

```bash
curl -X POST -F "file=@img.jpg" http://localhost:8080/upload
```

Ожидаемый ответ:
```json
{"id":"<id>","status":"queued"}
```


## Проверка наличия изображения

```bash
curl http://localhost:8080/image/<id> --output result.jpg
```

Ожидается ответ **200 OK** и сохранение файла.


## Удаление изображения

```bash
curl -X DELETE http://localhost:8080/image/<id>
```


## Проверка через браузер

Можно открыть интерфейс в браузере по адресу:  
[http://localhost:8080](http://localhost:8080)
