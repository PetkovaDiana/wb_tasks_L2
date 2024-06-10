package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"main.go/internal/model"
	"main.go/internal/storage"
	"net/http"
	"strconv"
	"time"
)

// Store - interface for working with data store
type Store interface {
	CreateEvent(event *model.Event) error
	UpdateEvent(userID, eventID int, newEvent *model.Event) error
	DeleteEvent(userID, eventID int)
	GetEventsForWeek(date time.Time, userID int) ([]model.Event, error)
	GetEventsForDay(date time.Time, userID int) ([]model.Event, error)
	GetEventsForMonth(date time.Time, userID int) ([]model.Event, error)
}

// ResultResponse - result response struct
type ResultResponse struct {
	Result []model.Event `json:"result"`
}

// ErrorResponse - error response struct
type ErrorResponse struct {
	Err string `json:"error"`
}

// Handler - http handler struct
type Handler struct {
	eventService Store
}

// NewHandler - creates new handler instance
func NewHandler() *Handler {
	return &Handler{
		eventService: storage.NewEventStorage(),
	}
}

// Register registers all routes to mux
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/create_event", h.CreateEvent)
	mux.HandleFunc("/update_event", h.UpdateEvent)
	mux.HandleFunc("/delete_event", h.DeleteEvent)
	mux.HandleFunc("/events_for_day", h.GetEventsForDay)
	mux.HandleFunc("/events_for_week", h.GetEventsForWeek)
	mux.HandleFunc("/events_for_month", h.GetEventsForMonth)
}

// CreateEvent - gets request data and passes to the service for creating
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := h.decodeJSON(r)
	if err != nil {
		h.errorResponse(w, fmt.Errorf("error while decoding input value: %v", err), http.StatusBadRequest)
		return
	}

	err = h.eventService.CreateEvent(event)
	if err != nil {
		h.errorResponse(w, fmt.Errorf("error while : %v", err), http.StatusServiceUnavailable)
		return
	}

	h.resultResponse(w, []model.Event{*event})
}

// DeleteEvent - gets request data and passes to the service for deleting
func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	event, err := h.decodeJSON(r)
	if err != nil {
		h.errorResponse(w, fmt.Errorf("error while decoding input value: %v", err), http.StatusBadRequest)
		return
	}

	h.eventService.DeleteEvent(event.UserID, event.EventID)

	h.resultResponse(w, []model.Event{*event})
}

// UpdateEvent - gets request data and passes to the service for updating
func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	event, err := h.decodeJSON(r)
	if err != nil {
		h.errorResponse(w, fmt.Errorf("error while decoding input value: %v", err), http.StatusBadRequest)
		return
	}

	err = h.eventService.UpdateEvent(event.UserID, event.EventID, event)
	if err != nil {
		h.errorResponse(w, err, http.StatusServiceUnavailable)
		return
	}

	h.resultResponse(w, []model.Event{*event})
}

// GetEventsForDay - gets request for event for day and  returns slice of events
func (h *Handler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	eventDate, err := h.ParseDate(date)
	if err != nil {
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetEventsForDay(eventDate, uID)
	if err != nil {
		h.errorResponse(w, err, http.StatusServiceUnavailable)
		return
	}

	h.resultResponse(w, events)
}

// GetEventsForWeek - gets request for event for week and  returns slice of events
func (h *Handler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	eventDate, err := h.ParseDate(date)
	if err != nil {
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetEventsForWeek(eventDate, uID)
	if err != nil {
		h.errorResponse(w, err, http.StatusServiceUnavailable)
		return
	}

	h.resultResponse(w, events)
}

// GetEventsForMonth - gets request for event for month and  returns slice of events
func (h *Handler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")

	uID, err := strconv.Atoi(userID)
	if err != nil || uID < 1 {
		if uID < 1 {
			err = errors.New("userID should be positive")
		}
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	eventDate, err := h.ParseDate(date)
	if err != nil {
		h.errorResponse(w, err, http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetEventsForMonth(eventDate, uID)
	if err != nil {
		h.errorResponse(w, err, http.StatusServiceUnavailable)
		return
	}

	h.resultResponse(w, events)
}

// decodeJSON - decode json format from request and returns event
func (h *Handler) decodeJSON(r *http.Request) (*model.Event, error) {
	var event model.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return nil, err
	}

	if event.UserID < 1 || event.EventID < 1 {
		return nil, errors.New("eventID or userID should pe positive")
	}

	return &event, nil
}

// ResultResponse - positive response
func (h *Handler) resultResponse(w http.ResponseWriter, events []model.Event) {
	w.Header().Set("Content-Type", "application/json")

	result, _ := json.MarshalIndent(&ResultResponse{Result: events}, " ", "")
	_, err := w.Write(result)
	if err != nil {
		h.errorResponse(w, err, http.StatusInternalServerError)
		return
	}
}

// ErrorResponse - response with error status
func (h *Handler) errorResponse(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")

	jsonErr, _ := json.MarshalIndent(&ErrorResponse{Err: err.Error()}, " ", " ")
	http.Error(w, string(jsonErr), status)
}

// ParseDate - parsing date from string
func (h *Handler) ParseDate(date string) (time.Time, error) {
	var (
		eventDate time.Time
		err       error
	)

	eventDate, err = time.Parse("2006-01-02T15:04", date)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02", date)
		if err != nil {
			eventDate, err = time.Parse("2006-01-02T15:04:00Z", date)
			if err != nil {
				return time.Time{}, fmt.Errorf("date format: e.g. 2022-05-10T14:10 error: %v", err)
			}
		}
	}

	return eventDate, nil
}
