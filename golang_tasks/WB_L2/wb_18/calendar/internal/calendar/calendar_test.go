package calendar

import (
	"testing"
	"time"
)

func TestCalendar_CreateEvent(t *testing.T) {
	cal := New()
	now := time.Now()

	event, err := cal.CreateEvent(1, now, "Test event")
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	if event.ID != 1 {
		t.Errorf("Expected event ID 1, got %d", event.ID)
	}
	if event.UserID != 1 {
		t.Errorf("Expected user ID 1, got %d", event.UserID)
	}
	if event.Content != "Test event" {
		t.Errorf("Expected content 'Test event', got '%s'", event.Content)
	}
}

func TestCalendar_UpdateEvent(t *testing.T) {
	cal := New()
	now := time.Now()

	event, _ := cal.CreateEvent(1, now, "Test event")
	
	updatedEvent, err := cal.UpdateEvent(event.ID, 1, now.Add(24*time.Hour), "Updated event")
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	if updatedEvent.Content != "Updated event" {
		t.Errorf("Expected updated content 'Updated event', got '%s'", updatedEvent.Content)
	}
}

func TestCalendar_DeleteEvent(t *testing.T) {
	cal := New()
	now := time.Now()

	event, _ := cal.CreateEvent(1, now, "Test event")
	
	err := cal.DeleteEvent(event.ID, 1)
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	// Try to get deleted event
	events := cal.GetEventsForDay(1, now)
	if len(events) != 0 {
		t.Errorf("Expected 0 events after deletion, got %d", len(events))
	}
}

func TestCalendar_GetEventsForPeriod(t *testing.T) {
	cal := New()
	now := time.Now()

	cal.CreateEvent(1, now, "Today event")
	cal.CreateEvent(1, now.Add(24*time.Hour), "Tomorrow event")
	cal.CreateEvent(1, now.Add(7*24*time.Hour), "Next week event")

	dayEvents := cal.GetEventsForDay(1, now)
	if len(dayEvents) != 1 {
		t.Errorf("Expected 1 event for day, got %d", len(dayEvents))
	}

	weekEvents := cal.GetEventsForWeek(1, now)
	if len(weekEvents) < 2 {
		t.Errorf("Expected at least 2 events for week, got %d", len(weekEvents))
	}

	monthEvents := cal.GetEventsForMonth(1, now)
	if len(monthEvents) != 3 {
		t.Errorf("Expected 3 events for month, got %d", len(monthEvents))
	}
}