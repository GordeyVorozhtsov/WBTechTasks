package server

import "time"

type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type CreateEventRequest struct {
	UserID  int    `json:"user_id" form:"user_id" binding:"required"`
	Date    string `json:"date" form:"date" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}

type UpdateEventRequest struct {
	ID      int    `json:"id" form:"id" binding:"required"`
	UserID  int    `json:"user_id" form:"user_id" binding:"required"`
	Date    string `json:"date" form:"date"`
	Content string `json:"content" form:"content"`
}

type DeleteEventRequest struct {
	ID     int `json:"id" form:"id" binding:"required"`
	UserID int `json:"user_id" form:"user_id" binding:"required"`
}

type EventsQuery struct {
	UserID int    `form:"user_id" binding:"required"`
	Date   string `form:"date" binding:"required"`
}

func (r *CreateEventRequest) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.Date)
}

func (r *UpdateEventRequest) ParseDate() (time.Time, error) {
	if r.Date == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", r.Date)
}

func (q *EventsQuery) ParseDate() (time.Time, error) {
	return time.Parse("2006-01-02", q.Date)
}