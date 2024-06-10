package main

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestFindAnagrams(t *testing.T) {
	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
		"ирак":   {"ирак", "каир", "раки"},
	}

	dictionary := []string{
		"Пятак", "Пятак", "пятка", "Тяпка",
		"слиток", "слиток", "столик", "листок",
		"Топот", "Потоп", "Каир", "Ирак", "раки",
	}

	myResult := findAnagrams(dictionary)

	normalizeMap := func(m map[string][]string) map[string][]string {
		nm := make(map[string][]string)
		for k, v := range m {
			sort.Strings(v)
			nm[strings.ToLower(k)] = v
		}
		return nm
	}

	expected = normalizeMap(expected)
	myResult = normalizeMap(myResult)

	if !reflect.DeepEqual(expected, myResult) {
		t.Errorf("Expected %v, but got %v", expected, myResult)
	}
}
