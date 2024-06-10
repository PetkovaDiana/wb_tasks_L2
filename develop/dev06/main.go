package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
Реализовать утилиту аналог консольной команды cut (man cut). Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/

type Args struct {
	f string
	d string
	s bool
}

// getArgs -  функция парсит флаги командной строки и возвращает структуру Args. Если флаг -f не указан, возвращается ошибка.
func getArgs() (*Args, error) {
	f := flag.String("f", "", "select only these fields(columns)")
	d := flag.String("d", "\t", "use a different delimiter")
	s := flag.Bool("s", false, "only print lines with the delimiter")

	flag.Parse()

	if *f == "" {
		return nil, errors.New("you need to specify fields with -f")
	}

	return &Args{
		f: *f,
		d: *d,
		s: *s,
	}, nil
}

func parseFields(fields string) ([]int, error) {
	var result []int

	for _, part := range strings.Split(fields, ",") {
		var field int
		_, err := fmt.Sscanf(part, "%d", &field)
		if err != nil {
			return nil, fmt.Errorf("invalid field value: %s", part)
		}
		result = append(result, field-1)
	}

	return result, nil
}

func processLine(line string, args *Args, fields []int) string {
	parts := strings.Split(line, args.d)
	if args.s && len(parts) < 2 {
		return ""
	}

	var result []string
	for _, field := range fields {
		if field < len(parts) {
			result = append(result, parts[field])
		}
	}

	return strings.Join(result, args.d)
}

func main() {
	args, err := getArgs()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fields, err := parseFields(args.f)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		result := processLine(line, args, fields)
		if result != "" || !args.s {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

}
