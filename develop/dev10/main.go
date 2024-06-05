package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/*
Утилита telnet

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

// args - структура для аргументов
type args struct {
	host    string
	port    string
	timeout time.Duration
}

// getArgs - возвращает парсинг аргументов
func getArgs() (*args, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("you need to specify HOST and PORT")
	}

	var (
		timeout time.Duration
		host    string
		port    string
	)

	if strings.Contains(os.Args[1], "--timeout=") {
		modif := os.Args[1][len(os.Args[1])-1]
		if modif != 's' {
			return nil, errors.New("you need to specify time unit: e.g.: 10s")
		}

		index := strings.Index(os.Args[1], "=")
		num, err := strconv.Atoi(os.Args[1][index+1 : len(os.Args[1])-1])
		if err != nil || num < 1 {
			return nil, err
		}

		timeout = time.Duration(num) * time.Second
		host = os.Args[2]
		port = os.Args[3]
	} else {
		host = os.Args[1]
		port = os.Args[2]
		timeout = time.Second * 10
	}

	return &args{
		host:    host,
		port:    port,
		timeout: timeout,
	}, nil
}

// readSocket - чтение из conn и вывод stdout
func readSocket(conn net.Conn, errChan chan<- error) {
	input := make([]byte, 1024)

	go func() {
		for {
			n, err := conn.Read(input)
			if err != nil {
				select {
				case errChan <- fmt.Errorf("remoute server stopped: %v", err):
				default:
				}
				return
			}
			fmt.Println(string(input[:n]))
		}
	}()

}

// writeSocket- чтение из stdin и запись в conn
func writeSocket(conn net.Conn, errChan chan<- error) {
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			text, err := reader.ReadBytes('\n')
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
			text = text[:len(text)-1]

			_, err = conn.Write(text)
			if err != nil {
				errChan <- err
				return
			}
		}
	}()
}

// telnet - подключается с помощью args и считывает данные из сокета и выводит в stdin
func telnet(args *args) error {
	address := fmt.Sprintf("%s:%s", args.host, args.port)

	//fmt.Println("Connecting to", address, "...")

	conn, err := net.DialTimeout("tcp", address, args.timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Println("Connecting to", address)

	sigs := make(chan os.Signal, 1)

	//syscall.SIGIN - сигнал прерывания
	//syscall.SIGTERM - сигнал завершения
	//syscall.SIGQUIT - сигнал завершения с дампом памяти
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	errChan := make(chan error)

	readSocket(conn, errChan)
	writeSocket(conn, errChan)

	select {
	case s := <-sigs: // если приходит сигнал
		fmt.Println("\nConnection stopped by signal:", s)
	case e := <-errChan: // если приходит ошибка из горутин
		fmt.Println("Connection stopped by", e)
	}
	return nil
}

func main() {

	arg, err := getArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	err = telnet(arg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
