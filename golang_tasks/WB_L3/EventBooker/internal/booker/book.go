package booker

import (
	"context"
	"errors"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

type Event struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Capacity  int       `json:"capacity"`
	Timeout   int       `json:"timeout"`
	CreatedAt time.Time `json:"created_at"`
}

type Booking struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	UserEmail string    `json:"user_email"`
	Places    int       `json:"places"`
	Confirmed bool      `json:"confirmed"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type EventWithDetails struct {
	Event
	BookedPlaces    int `json:"booked_places"`
	FreePlaces      int `json:"free_places"`
	ConfirmedPlaces int `json:"confirmed_places"`
}

var (
	ErrEventNotFound    = errors.New("event not found")
	ErrNotEnoughPlaces  = errors.New("not enough places")
	ErrBookingNotFound  = errors.New("booking not found")
	ErrAlreadyConfirmed = errors.New("booking already confirmed")
)

type Service struct {
	DB *dbpg.DB
}

func (s *Service) StartCleanupScheduler(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.CleanupExpiredBookings(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) CreateEvent(ctx context.Context, title string, date time.Time, capacity, timeout int) (*Event, error) {
	const query = `INSERT INTO events (title, date, capacity, timeout_seconds, created_at) 
	               VALUES ($1, $2, $3, $4, $5) RETURNING id, title, date, capacity, timeout_seconds, created_at`

	event := &Event{}
	err := s.DB.QueryRowContext(ctx, query, title, date, capacity, timeout, time.Now()).
		Scan(&event.ID, &event.Title, &event.Date, &event.Capacity, &event.Timeout, &event.CreatedAt)

	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) BookEvent(ctx context.Context, eventID int, userEmail string, places int) (*Booking, error) {
	tx, err := s.DB.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var capacity, timeout int
	err = tx.QueryRowContext(ctx, "SELECT capacity, timeout_seconds FROM events WHERE id = $1", eventID).
		Scan(&capacity, &timeout)
	if err != nil {
		return nil, ErrEventNotFound
	}

	var confirmedPlaces int
	tx.QueryRowContext(ctx, "SELECT COALESCE(SUM(places), 0) FROM bookings WHERE event_id = $1 AND confirmed = true", eventID).
		Scan(&confirmedPlaces)

	if confirmedPlaces+places > capacity {
		return nil, ErrNotEnoughPlaces
	}

	booking := &Booking{}
	now := time.Now()
	err = tx.QueryRowContext(ctx,
		"INSERT INTO bookings (event_id, user_email, places, confirmed, created_at, expires_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, event_id, user_email, places, confirmed, created_at, expires_at",
		eventID, userEmail, places, false, now, now.Add(time.Duration(timeout)*time.Second),
	).Scan(&booking.ID, &booking.EventID, &booking.UserEmail, &booking.Places, &booking.Confirmed, &booking.CreatedAt, &booking.ExpiresAt)

	if err != nil {
		return nil, err
	}

	tx.Commit()
	return booking, nil
}

func (s *Service) ConfirmBooking(ctx context.Context, eventID int, userEmail string) (*Booking, error) {
	tx, err := s.DB.Master.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var capacity int
	err = tx.QueryRowContext(ctx, "SELECT capacity FROM events WHERE id = $1", eventID).Scan(&capacity)
	if err != nil {
		return nil, ErrEventNotFound
	}

	booking := &Booking{}
	err = tx.QueryRowContext(ctx, `
		SELECT id, event_id, user_email, places, confirmed
		FROM bookings
		WHERE event_id = $1 AND user_email = $2
		ORDER BY id DESC LIMIT 1
	`, eventID, userEmail).Scan(&booking.ID, &booking.EventID, &booking.UserEmail, &booking.Places, &booking.Confirmed)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.Confirmed {
		return nil, ErrAlreadyConfirmed
	}

	var confirmedPlaces int
	_ = tx.QueryRowContext(ctx, `
		SELECT COALESCE(SUM(places), 0)
		FROM bookings
		WHERE event_id = $1 AND confirmed = true AND id != $2
	`, eventID, booking.ID).Scan(&confirmedPlaces)

	if confirmedPlaces+booking.Places > capacity {
		return nil, ErrNotEnoughPlaces
	}

	_, err = tx.ExecContext(ctx, "UPDATE bookings SET confirmed = true WHERE id = $1", booking.ID)
	if err != nil {
		return nil, err
	}

	booking.Confirmed = true
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) GetEvent(ctx context.Context, eventID int) (*EventWithDetails, error) {
	const query = `SELECT e.id, e.title, e.date, e.capacity, e.timeout_seconds, e.created_at,
	                      COALESCE(SUM(CASE WHEN b.confirmed = true THEN b.places ELSE 0 END), 0),
	                      COALESCE(SUM(b.places), 0)
	               FROM events e LEFT JOIN bookings b ON e.id = b.event_id 
	               WHERE e.id = $1 GROUP BY e.id`

	event := &EventWithDetails{}
	err := s.DB.QueryRowContext(ctx, query, eventID).
		Scan(&event.ID, &event.Title, &event.Date, &event.Capacity, &event.Timeout, &event.CreatedAt, &event.ConfirmedPlaces, &event.BookedPlaces)

	if err != nil {
		return nil, ErrEventNotFound
	}

	event.FreePlaces = event.Capacity - event.ConfirmedPlaces
	return event, nil
}

func (s *Service) GetAllEvents(ctx context.Context) ([]*EventWithDetails, error) {
	const query = `SELECT e.id, e.title, e.date, e.capacity, e.timeout_seconds, e.created_at,
	                      COALESCE(SUM(CASE WHEN b.confirmed = true THEN b.places ELSE 0 END), 0),
	                      COALESCE(SUM(b.places), 0)
	               FROM events e LEFT JOIN bookings b ON e.id = b.event_id 
	               GROUP BY e.id ORDER BY e.created_at DESC`

	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*EventWithDetails
	for rows.Next() {
		event := &EventWithDetails{}
		err := rows.Scan(&event.ID, &event.Title, &event.Date, &event.Capacity, &event.Timeout, &event.CreatedAt, &event.ConfirmedPlaces, &event.BookedPlaces)
		if err != nil {
			return nil, err
		}
		event.FreePlaces = event.Capacity - event.ConfirmedPlaces
		events = append(events, event)
	}
	return events, nil
}

func (s *Service) GetExpiredBookings(ctx context.Context) ([]*Booking, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT id, event_id, user_email, places FROM bookings WHERE confirmed = false AND expires_at < $1", time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*Booking
	for rows.Next() {
		booking := &Booking{}
		rows.Scan(&booking.ID, &booking.EventID, &booking.UserEmail, &booking.Places)
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func (s *Service) DeleteBooking(ctx context.Context, bookingID int) error {
	_, err := s.DB.ExecContext(ctx, "DELETE FROM bookings WHERE id = $1", bookingID)
	return err
}

func (s *Service) CleanupExpiredBookings(ctx context.Context) error {
	bookings, err := s.GetExpiredBookings(ctx)
	if err != nil {
		return err
	}

	for _, booking := range bookings {
		s.DeleteBooking(ctx, booking.ID)
	}
	return nil
}

func (s *Service) InitDatabase(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			date TIMESTAMP NOT NULL,
			capacity INTEGER NOT NULL,
			timeout_seconds INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		
		CREATE TABLE IF NOT EXISTS bookings (
			id SERIAL PRIMARY KEY,
			event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
			user_email VARCHAR(255) NOT NULL,
			places INTEGER NOT NULL,
			confirmed BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			expires_at TIMESTAMP NOT NULL
		);
		
		CREATE INDEX IF NOT EXISTS idx_bookings_expires_at ON bookings(expires_at) WHERE NOT confirmed;
	`)
	return err
}
