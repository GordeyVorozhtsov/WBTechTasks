package server

import (
	"net/http"
	"calendar/internal/calendar"
	"github.com/gin-gonic/gin"
)

type Handlers struct {
	calendar *calendar.Calendar
}

func NewHandlers(cal *calendar.Calendar) *Handlers {
	return &Handlers{calendar: cal}
}

func (h *Handlers) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid request: " + err.Error()})
		return
	}

	date, err := req.ParseDate()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	event, err := h.calendar.CreateEvent(req.UserID, date, req.Content)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: event})
}

func (h *Handlers) UpdateEvent(c *gin.Context) {
	var req UpdateEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid request: " + err.Error()})
		return
	}

	date, err := req.ParseDate()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	event, err := h.calendar.UpdateEvent(req.ID, req.UserID, date, req.Content)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: event})
}

func (h *Handlers) DeleteEvent(c *gin.Context) {
	var req DeleteEventRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid request: " + err.Error()})
		return
	}

	err := h.calendar.DeleteEvent(req.ID, req.UserID)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, Response{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Result: "Event deleted successfully"})
}

func (h *Handlers) EventsForDay(c *gin.Context) {
	h.eventsForPeriod(c, "day")
}

func (h *Handlers) EventsForWeek(c *gin.Context) {
	h.eventsForPeriod(c, "week")
}

func (h *Handlers) EventsForMonth(c *gin.Context) {
	h.eventsForPeriod(c, "month")
}

func (h *Handlers) eventsForPeriod(c *gin.Context, period string) {
	var query EventsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid query parameters: " + err.Error()})
		return
	}

	date, err := query.ParseDate()
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	var events []*calendar.Event
	switch period {
	case "day":
		events = h.calendar.GetEventsForDay(query.UserID, date)
	case "week":
		events = h.calendar.GetEventsForWeek(query.UserID, date)
	case "month":
		events = h.calendar.GetEventsForMonth(query.UserID, date)
	default:
		c.JSON(http.StatusBadRequest, Response{Error: "Invalid period"})
		return
	}

	c.JSON(http.StatusOK, Response{Result: events})
}