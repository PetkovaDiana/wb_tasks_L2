package main

import (
	"testing"
	"time"
)

func TestGetCurrentTime(t *testing.T) {
	passingCase := time.Now()

	failureCases := []time.Time{
		time.Now().Add(time.Hour * 2),
		time.Now().Add(time.Minute * 2),
		time.Now().Add(time.Second * 2),
	}

	currTime, err := getCurrentTime()
	if err != nil {
		t.Error(err)
	}

	// testing pass case
	if currTime.Before(passingCase) {
		t.Errorf("%v != %v", passingCase, currTime)
	}

	// testing failure cases
	for _, val := range failureCases {
		if val.Before(currTime) {
			t.Errorf("%v == %v\nShould: %v\nGot: %v\n", currTime, val, currTime, val)
		}
	}
}
