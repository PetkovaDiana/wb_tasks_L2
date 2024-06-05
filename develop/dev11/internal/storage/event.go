package storage

import (
	"errors"
	"fmt"
	"main.go/dto"
	"sync"
	"time"
)

// EventStorage - это структура, которая хранит данные о событиях
type EventStorage struct {
	sync.RWMutex
	db map[string]dto.Event
}

// Создаем новую бд
func NewEventStorage() *EventStorage {
	return &EventStorage{
		db: make(map[string]dto.Event),
	}
}

func (e *EventStorage) CreateEvent(event *dto.Event) error {
	id := fmt.Sprintf("%d:%d", event.UserID, event.EventID)

	e.Lock()
	if _, ok := e.db[id]; ok {
		e.Unlock()
		return errors.New("event with such id already exist")
	}
	e.db[id] = *event
	e.Unlock()

	return nil
}

func (e *EventStorage) UpdateEvent(userID, eventID int, newEvent *dto.Event) error {
	combinedID := fmt.Sprintf("%d:%d", userID, eventID)

	e.Lock()
	if _, ok := e.db[combinedID]; ok {
		e.Unlock()
		return fmt.Errorf("there is no event with id: %s", combinedID)
	}

	e.db[combinedID] = *newEvent
	e.Unlock()

	return nil
}

func (e *EventStorage) DeleteEvent(userID, eventID int) {
	id := fmt.Sprintf("%d:%d", userID, eventID)

	e.Lock()
	delete(e.db, id)
	e.Unlock()
}

func (e *EventStorage) GetEventsForWeek(date time.Time, userID int) ([]dto.Event, error) {
	var eventsForWeek []dto.Event

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

// GetEventsForWeek - returns all events for current week

func (e *EventStorage) GetEventsForDay(date time.Time, userID int) ([]dto.Event, error) {
	var eventsForDay []dto.Event

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

func (e *EventStorage) GetEventsForMonth(date time.Time, userID int) ([]dto.Event, error) {
	var GetEventsForMonth []dto.Event

	y, m, _ := date.Date()

	e.RLock()

	for _, event := range e.db {
		eventY, eventM, _ := event.Date.Date()
		if y == eventY && int(m) == int(eventM) && userID == event.UserID {
			GetEventsForMonth = append(GetEventsForMonth, event)
		}
	}

	e.RUnlock()

	return GetEventsForMonth, nil
}
