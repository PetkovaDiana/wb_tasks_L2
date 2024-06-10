package storage

import (
	"errors"
	"fmt"
	"main.go/internal/model"
	"sync"
	"time"
)

// EventStorage - database structure
type EventStorage struct {
	sync.RWMutex
	db map[string]model.Event
}

// NewEventStorage - creates new database instance
func NewEventStorage() *EventStorage {
	return &EventStorage{
		db: make(map[string]model.Event),
	}
}

// CreateEvent - creating new event in data store
func (e *EventStorage) CreateEvent(event *model.Event) error {
	id := fmt.Sprintf("%d%d", event.UserID, event.EventID)

	e.Lock()

	if _, ok := e.db[id]; ok {
		e.Unlock()
		return errors.New("event with such id already exist")
	}
	e.db[id] = *event

	e.Unlock()

	return nil
}

// UpdateEvent - updating event in data store
func (e *EventStorage) UpdateEvent(userID, eventID int, newEvent *model.Event) error {
	combinedID := fmt.Sprintf("%d%d", userID, eventID)

	e.Lock()

	if _, ok := e.db[combinedID]; !ok {
		e.Unlock()
		return fmt.Errorf("there is no event with id: %s", combinedID)
	}

	e.db[combinedID] = *newEvent

	e.Unlock()

	return nil
}

// DeleteEvent - deleting event from data store
func (e *EventStorage) DeleteEvent(userID, eventID int) {

	id := fmt.Sprintf("%d%d", userID, eventID)

	e.Lock()

	delete(e.db, id)

	e.Unlock()
}

// GetEventsForWeek - returns all events for current week
func (e *EventStorage) GetEventsForWeek(date time.Time, userID int) ([]model.Event, error) {
	var eventsForWeek []model.Event

	currYear, currWeek := date.ISOWeek()

	e.RLock()
	for _, event := range e.db {
		eventYear, eventWeek := event.Date.ISOWeek()
		time.Now().ISOWeek()
		if eventYear == currYear && eventWeek == currWeek && userID == event.UserID {
			eventsForWeek = append(eventsForWeek, event)
		}
	}
	e.RUnlock()

	return eventsForWeek, nil
}

// GetEventsForDay - returns all events for current day
func (e *EventStorage) GetEventsForDay(date time.Time, userID int) ([]model.Event, error) {
	var eventsForDay []model.Event

	y, m, d := date.Date()

	e.RLock()

	for _, event := range e.db {
		eventY, eventM, eventD := event.Date.Date()
		if y == eventY && int(eventM) == int(m) && d == eventD && userID == event.UserID {
			eventsForDay = append(eventsForDay, event)
		}
	}

	e.RUnlock()

	return eventsForDay, nil
}

// GetEventsForMonth - returns all events for current month
func (e *EventStorage) GetEventsForMonth(date time.Time, userID int) ([]model.Event, error) {
	var eventsForMonth []model.Event

	y, m, _ := date.Date()

	e.RLock()

	for _, event := range e.db {
		eventY, eventM, _ := event.Date.Date()
		if y == eventY && int(eventM) == int(m) && userID == event.UserID {
			eventsForMonth = append(eventsForMonth, event)
		}
	}

	e.RUnlock()

	return eventsForMonth, nil
}
