create event
curl -X POST http://localhost:8080/create_event -H "Content-Type: application/json" -d '{"user_id": 1, "date": "2025-01-09", "content": "Meeting"}'

get event
curl "http://localhost:8080/events_for_day?user_id=1&date=2024-01-09"

delete event
curl -X POST http://localhost:8080/delete_event -H "Content-Type: application/json" -d '{"id": 1, "user_id": 1}'