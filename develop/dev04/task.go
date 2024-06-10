package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
Написать функцию поиска всех множеств анаграмм по словарю.

Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.


Требования:
Входные данные для функции: ссылка на массив, каждый элемент которого - слово на русском языке в кодировке utf8
Выходные данные: ссылка на мапу множеств анаграмм
Ключ - первое встретившееся в словаре слово из множества. Значение - ссылка на массив, каждый элемент которого,
слово из множества.
Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

*/

func filter(tmp map[string][]string) {
	unique := make(map[string]bool)

	for key, val := range tmp {
		// удалить короткие наборы
		if len(val) < 2 {
			delete(tmp, key)
		}

		// удалить повторяющиеся слова
		for i := range val {
			if !unique[val[i]] {
				unique[val[i]] = true
			} else {
				val[i] = val[len(val)-1]
				tmp[key] = val[:len(val)-1]
			}
		}
	}
}

func getLettersInAlphabetOrder(word string) string {
	letters := strings.Split(word, "") // Разбиваем слово на отдельные буквы.
	sort.Strings(letters)              // Сортируем буквы в алфавитном порядке.
	return strings.Join(letters, "")   // Соединяем буквы обратно в строку.

}

func makeTemporaryMapOfAnagrams(dictionary []string) map[string][]string {
	tmp := make(map[string][]string) //Создаем временную карту для хранения анаграмм.

	//заполняем карту словами
	for _, val := range dictionary {
		loweredWord := strings.ToLower(val)
		letters := getLettersInAlphabetOrder(loweredWord)
		tmp[letters] = append(tmp[letters], loweredWord)
	}

	return tmp
}

func findAnagrams(dictionary []string) map[string][]string {
	tmp := makeTemporaryMapOfAnagrams(dictionary)
	filter(tmp)
	anagrams := make(map[string][]string, len(tmp))

	for _, val := range tmp {
		sort.Strings(val)
		anagrams[val[0]] = val
	}

	return anagrams
}

func main() {
	dictionary := []string{
		"Пятак",
		"Пятак",
		"пятка",
		"Тяпка",
		"слиток",
		"слиток",
		"столик",
		"листок",
		"Топот",
		"Потоп",
	}

	anagrams := findAnagrams(dictionary)

	for k, v := range anagrams {
		fmt.Printf("Key: %s\nValue: %v\n\n", k, v)
	}
}
