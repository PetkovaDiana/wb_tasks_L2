package dev11

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main.go/internal/handler"
	"main.go/internal/model"
	"net/http/httptest"
	"os"
	"sort"
	"testing"
)

func TestService(t *testing.T) {
	testData := readData("test_data/test_data")
	api := handler.NewHandler()

	// fill data
	for _, event := range testData {
		jsonData := jsonEncode(event)

		body := bytes.NewBuffer(jsonData)

		r := httptest.NewRequest("POST", "http://127.0.0.1:8080/create_user", body)
		w := httptest.NewRecorder()

		var responseEvent handler.ResultResponse

		api.CreateEvent(w, r)

		err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
			os.Exit(2)
		}

		if !compare(event, responseEvent.Result[0]) {
			t.Errorf("events are not equal\nShould: %v\nGot: %v", event, responseEvent)
			os.Exit(2)
		}
	}

	t.Run("update data", func(t *testing.T) {
		for _, event := range testData {
			event.Title += " updated title"
			jsonData := jsonEncode(event)

			body := bytes.NewBuffer(jsonData)

			r := httptest.NewRequest("POST", "http://127.0.0.1:8080/update_event", body)
			w := httptest.NewRecorder()

			var responseEvent handler.ResultResponse

			api.UpdateEvent(w, r)

			err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
				os.Exit(2)
			}

			if event.Title != responseEvent.Result[0].Title {
				t.Errorf("events are not equal\nShould: %v\nGot: %v", event, responseEvent)
				os.Exit(2)
			}
		}
	})

	t.Run("event for day", func(t *testing.T) {
		event := testData[0]
		jsonData := jsonEncode(event)

		body := bytes.NewBuffer(jsonData)

		addr := fmt.Sprintf("http://127.0.0.1:8080/events_for_day?user_id=%d&date=%s",
			event.UserID, event.Date.Format("2006-01-02"))

		r := httptest.NewRequest("POST", addr, body)
		w := httptest.NewRecorder()

		var responseEvent handler.ResultResponse

		api.GetEventsForDay(w, r)

		err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
			os.Exit(2)
		}

		sort.Slice(responseEvent.Result, func(i, j int) bool {
			return responseEvent.Result[i].EventID < responseEvent.Result[j].EventID
		})

		for i, events := range testData[:3] {
			if !compare(events, responseEvent.Result[i]) {
				t.Errorf("events are not equal\nShould: %v\nGot: %v", testData[:3], responseEvent)
				os.Exit(2)
			}
		}
	})

	t.Run("event for week", func(t *testing.T) {
		event := testData[6]
		jsonData := jsonEncode(event)

		body := bytes.NewBuffer(jsonData)

		addr := fmt.Sprintf("http://127.0.0.1:8080/events_for_week?user_id=%d&date=%s",
			event.UserID, event.Date.Format("2006-01-02"))

		r := httptest.NewRequest("POST", addr, body)
		w := httptest.NewRecorder()

		var responseEvent handler.ResultResponse

		api.GetEventsForWeek(w, r)

		err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
			os.Exit(2)
		}

		sort.Slice(responseEvent.Result, func(i, j int) bool {
			return responseEvent.Result[i].EventID < responseEvent.Result[j].EventID
		})

		for i, events := range testData[6:9] {
			if !compare(events, responseEvent.Result[i]) {
				t.Errorf("events are not equal\nShould: %v\nGot: %v", testData[:3], responseEvent)
				os.Exit(2)
			}
		}
	})

	t.Run("event for month", func(t *testing.T) {
		event := testData[9]
		jsonData := jsonEncode(event)

		body := bytes.NewBuffer(jsonData)

		addr := fmt.Sprintf("http://127.0.0.1:8080/events_for_month?user_id=%d&date=%s",
			event.UserID, event.Date.Format("2006-01-02"))

		r := httptest.NewRequest("POST", addr, body)
		w := httptest.NewRecorder()

		var responseEvent handler.ResultResponse

		api.GetEventsForMonth(w, r)

		err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
			os.Exit(2)
		}

		sort.Slice(responseEvent.Result, func(i, j int) bool {
			return responseEvent.Result[i].EventID < responseEvent.Result[j].EventID
		})

		for i, events := range testData[9:] {
			if !compare(events, responseEvent.Result[i]) {
				t.Errorf("events are not equal\nShould: %v\nGot: %v", testData[:3], responseEvent)
				os.Exit(2)
			}
		}
	})

	t.Run("remove data", func(t *testing.T) {
		for _, event := range testData {
			jsonData := jsonEncode(event)

			body := bytes.NewBuffer(jsonData)

			r := httptest.NewRequest("POST", "http://127.0.0.1:8080/delete_event", body)
			w := httptest.NewRecorder()

			var responseEvent handler.ResultResponse

			api.DeleteEvent(w, r)

			err := json.Unmarshal(w.Body.Bytes(), &responseEvent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
				os.Exit(2)
			}

			if !compare(event, responseEvent.Result[0]) {
				t.Errorf("events are not equal\nShould: %v\nGot: %v", event, responseEvent)
				os.Exit(2)
			}
		}
	})
}

func TestServiceWithBadData(t *testing.T) {
	t.Run("negative IDs", func(t *testing.T) {
		testData := readData("test_data/bad_data_1")
		api := handler.NewHandler()

		// fill data
		for _, event := range testData {
			jsonData := jsonEncode(event)

			body := bytes.NewBuffer(jsonData)

			r := httptest.NewRequest("POST", "http://127.0.0.1:8080/create_user", body)
			w := httptest.NewRecorder()

			var errEvent handler.ErrorResponse

			api.CreateEvent(w, r)

			err := json.Unmarshal(w.Body.Bytes(), &errEvent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
				os.Exit(2)
			}

			if errEvent.Err == "" {
				t.Errorf("events are not equal\nShould: error while : eventID or userID should pe positive\nGot: %v", errEvent)
				os.Exit(2)
			}
		}
	})

	t.Run("identical IDs", func(t *testing.T) {
		testData := readData("test_data/bad_data_2")
		api := handler.NewHandler()

		// fill data
		for _, event := range testData {
			jsonData := jsonEncode(event)

			body := bytes.NewBuffer(jsonData)

			r := httptest.NewRequest("POST", "http://127.0.0.1:8080/create_user", body)
			w := httptest.NewRecorder()

			var errEvent handler.ErrorResponse

			api.CreateEvent(w, r)

			err := json.Unmarshal(w.Body.Bytes(), &errEvent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
				os.Exit(2)
			}

			if event.EventID == 2 && errEvent.Err == "" {
				t.Errorf("events are not equal\nShould: error while : event with such id already exist\nGot: %v", errEvent)
				os.Exit(2)
			}
		}
	})

}

// compare two events
func compare(event1, event2 model.Event) bool {
	if event1.UserID != event2.UserID {
		return false
	}

	if event1.EventID != event2.EventID {
		return false
	}

	return true
}

// jsonEncode - returns json format of event
func jsonEncode(event model.Event) []byte {
	jsonData, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot marshal testing data: %v", err)
		os.Exit(2)
	}

	return jsonData
}

// readData - reads data from files and returns event
func readData(filenames string) []model.Event {
	var events []model.Event

	file, err := ioutil.ReadFile(filenames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot read testing data: %v", err)
		os.Exit(2)
	}

	err = json.Unmarshal(file, &events)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot parse data: %v", err)
		os.Exit(2)
	}

	return events
}
