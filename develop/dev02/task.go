package main

import (
	"errors"
	"fmt"
	"os"
)

/*
 Задача на распаковку

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var errCorrectness = errors.New("the string must not start with a digit and contain two digits in a row")

// isNumber - returns true if rune is between '1' - '9'
func isNumber(c rune) bool {
	if c >= '1' && c <= '9' {
		return true
	}
	return false
}

// isSlash - returns true if c is slash
func isSlash(c rune) bool {
	return c == '\\'
}

// isCorrect - checking for string correctness
func isCorrect(r []rune) bool {
	if isNumber(r[0]) {
		return false
	}

	var (
		twoDigit bool
		slash    bool
	)

	// checking for two numbers in a row
	for i := range r {
		if isSlash(r[i]) {
			slash = true
			continue
		} else if isNumber(r[i]) {
			if twoDigit && !slash {
				return false
			}
			twoDigit = true
			continue
		}
		slash = false
		twoDigit = false
	}
	return true
}

// repeatRune - repeats rune c count times
func repeatRune(c rune, count int) []rune {
	res := make([]rune, 0, count)
	for i := 0; i < count; i++ {
		res = append(res, c)
	}
	return res
}

/*
unpackString - unpacking string like:

	"a4bc2d5e" => "aaaabccddddde"
	"abcd" => "abcd"
	"45" => "" (incorrect string)
	"" => ""

	qwe\4\5 => qwe45
	qwe\45 => qwe44444
	qwe\\5 => qwe\\\\\
*/
func unpackString(str string) (string, error) {
	if len(str) < 1 {
		return str, nil
	}

	runes := []rune(str)

	if !isCorrect(runes) {
		return "", errCorrectness
	}

	unpackedStr := make([]rune, 0, len(runes))

	var slash bool

	for i, val := range runes {
		if isSlash(val) {
			if slash {
				unpackedStr = append(unpackedStr, val)
				slash = false
				continue
			}
			slash = true
		} else if isNumber(val) {
			if slash {
				unpackedStr = append(unpackedStr, val)
				slash = false
			} else {
				repeatedRune := repeatRune(runes[i-1], int(val-'0')-1)
				unpackedStr = append(unpackedStr, repeatedRune...)
			}
		} else {
			unpackedStr = append(unpackedStr, val)
		}
	}

	return string(unpackedStr), nil
}

func main() {
	strs := []string{
		"a4bc2d5e",
		"abcd",
		"45",
		"",
		`qwe\4\5`,
		`qwe\45`,
		`qwe\\5`,
	}

	for i := range strs {
		unpacked, err := unpackString(strs[i])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		fmt.Println(unpacked)
	}
}
