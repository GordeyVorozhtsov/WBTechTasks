package calendar

import (
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date")
)

type Event struct {
	ID      int       `json:"id"`
	UserID  int       `json:"user_id"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

type Calendar struct {
	events  []*Event
	nextID  int
}

func New() *Calendar {
	return &Calendar{
		events: make([]*Event, 0),
		nextID: 1,
	}
}

func (c *Calendar) CreateEvent(userID int, date time.Time, content string) (*Event, error) {
	if date.IsZero() {
		return nil, ErrInvalidDate
	}

	event := &Event{
		ID:      c.nextID,
		UserID:  userID,
		Date:    date,
		Content: content,
	}

	c.events = append(c.events, event)
	c.nextID++

	return event, nil
}

func (c *Calendar) UpdateEvent(id, userID int, date time.Time, content string) (*Event, error) {
	for _, event := range c.events {
		if event.ID == id && event.UserID == userID {
			if !date.IsZero() {
				event.Date = date
			}
			if content != "" {
				event.Content = content
			}
			return event, nil
		}
	}
	return nil, ErrEventNotFound
}

func (c *Calendar) DeleteEvent(id, userID int) error {
	for i, event := range c.events {
		if event.ID == id && event.UserID == userID {
			c.events = append(c.events[:i], c.events[i+1:]...)
			return nil
		}
	}
	return ErrEventNotFound
}

func (c *Calendar) GetEventsForDay(userID int, date time.Time) []*Event {
	var result []*Event
	for _, event := range c.events {
		if event.UserID == userID &&
			event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() &&
			event.Date.Day() == date.Day() {
			result = append(result, event)
		}
	}
	return result
}

func (c *Calendar) GetEventsForWeek(userID int, date time.Time) []*Event {
	var result []*Event
	year, week := date.ISOWeek()
	for _, event := range c.events {
		if event.UserID == userID {
			eventYear, eventWeek := event.Date.ISOWeek()
			if eventYear == year && eventWeek == week {
				result = append(result, event)
			}
		}
	}
	return result
}

func (c *Calendar) GetEventsForMonth(userID int, date time.Time) []*Event {
	var result []*Event
	for _, event := range c.events {
		if event.UserID == userID &&
			event.Date.Year() == date.Year() &&
			event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}
	return result
}