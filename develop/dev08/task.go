package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

/*
Взаимодействие с ОС


Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:


- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*

Так же требуется поддерживать функционал fork/exec-команд

*/

func readLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	currDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fmt.Printf("minishell:%s$", currDir)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// удаление \n
	input = input[:len(input)-1]
	return input, nil
}

// netcat - установить соединение с HOST:PORT через tcp/udp
func netcat(args []string) string {
	var (
		host string
		port string
	)

	protocol := "tcp"
	//Проверяем, что в аргументах командной строки указаны хост и порт. Если нет, возвращаем сообщение об ошибке.
	if len(args) < 2 {
		return "need to specify HOST and PORT\n"
	}

	//Если в аргументах командной строки указан параметр -u,
	//устанавливаем протокол в UDP и извлекаем хост и порт из следующих двух аргументов.
	//В противном случае, используем первые два аргумента как хост и порт.
	if len(args) >= 3 && args[0] == "-u" {
		protocol = "udp"
		host = args[1]
		port = args[2]
	} else {
		host = args[0]
		port = args[1]
	}

	addr := fmt.Sprintf("%s:%s", host, port)
	// Формируем строку адреса из хоста и порта, и устанавливаем соединение с помощью функции net.Dial.
	// В случае ошибки возвращаем сообщение о неудаче.
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		return fmt.Sprintf("netcat: connection failed: %v\n", err)
	}

	// Создаем каналы для сигналов и ошибок, и уведомляем о прерывании и завершении программы.
	sigs := make(chan os.Signal, 1)
	errors := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//Запускаем горутину, которая читает данные из стандартного ввода и отправляет их по соединению.
	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				errors <- err
			}

			_, err = conn.Write([]byte(input[:len(input)-1]))
			if err != nil {
				errors <- err
				return
			}
		}
	}()

	//В блоке select ожидаем сигналы или ошибки.
	//Если происходит ошибка, возвращаем ее сообщение.
	//Если получен сигнал, возвращаем сообщение о его получении.
	select {
	case err := <-errors:
		return fmt.Sprintf("%v\n", err)
	case s := <-sigs:
		return fmt.Sprintf("stopped by signal: %v\n", s)
	}

}

// cd -предназначена для изменения текущего рабочего каталога.
func cd(path []string) string {
	if len(path) > 1 {
		return "cd: too mane arguments\n"
	}

	if len(path) < 1 {
		home := os.Getenv("HOME")
		os.Chdir(home)
	} else {
		os.Chdir(path[0])
	}

	return ""
}

// ps - предназначена для получения списка всех активных процессов.
func ps() string {
	processes, err := process.Processes()
	if err != nil {
		return fmt.Sprintf("ps: error: %v\n", err)
	}
	out := strings.Builder{}
	out.WriteString("PID\tCMD\n")

	for _, p := range processes {
		name, _ := p.Name()

		out.WriteString(fmt.Sprintf("%d\t%s\n", p.Pid, name))
	}

	return out.String()
}

func kill(args []string) string {
	if len(args) < 1 {
		return "kill: usage: kill [-s sigspec  ] pid\n"
	}

	// Объявляем переменные для хранения PID процесса, указанного сигнала (если таковой имеется), флага, указывающего на наличие сигнала, и самого сигнала.
	var (
		pidToKill   string
		signalTitle string
		signalFlag  bool
		signal      syscall.Signal
	)

	for _, word := range args {
		if signalFlag {
			signalFlag = false
			signalTitle = word
			continue
		}

		if strings.Contains(word, "-s") {
			signalFlag = true
		} else {
			pidToKill = word
		}
	}

	// пустой -s сигнал
	// Если указан сигнал, но его значение не указано, возвращаем сообщение об ошибке.
	if signalFlag && len(signalTitle) < 1 {
		return fmt.Sprintf("kill: %s: invalid signal specification\n", signalTitle)
	}

	// default signal
	signal = syscall.SIGTERM

	switch signalTitle {
	case "SIGINT":
		signal = syscall.SIGINT
	case "SIGTERM":
		signal = syscall.SIGTERM
	case "SIGQUIT":
		signal = syscall.SIGQUIT
	case "SIGKILL":
		signal = syscall.SIGKILL
	case "SIGHUP":
		signal = syscall.SIGHUP
	}

	processes, err := process.Processes()
	if err != nil {
		return err.Error()
	}

	for _, p := range processes {
		pid := p.Pid
		if fmt.Sprintf("%d", pid) == pidToKill {
			err := p.SendSignal(signal)
			if err != nil {
				return fmt.Sprintf("kill: err: %v\n", err)
			}
			return fmt.Sprintf("kill: %s was %s\n", pidToKill, signal.String())
		}
	}

	return fmt.Sprintf("kill: there is no signal with pid: %s\n", pidToKill)
}

// applyCommand - принимает строку, содержащую команду, и выполняет соответствующую операцию в зависимости от этой команды.
func applyCommand(line string) string {
	if len(line) < 1 {
		return ""
	}

	// Разбиваем строку на слова и сохраняем их в переменной commands.
	// Если в строке больше одного слова, сохраняем остальные слова (аргументы) в переменной args.
	var args []string
	commands := strings.Fields(line)
	if len(commands) > 1 {
		args = commands[1:]
	}

	var outLine string

	switch commands[0] {
	case "cd":
		outLine = cd(args)
	case "pwd":
		outLine, _ = os.Getwd()
		outLine += "\n"
	case "echo":
		outLine = strings.Join(args, "")
		outLine += "\n"
	case "kill":
		outLine = kill(args)
	case "ps":
		outLine = ps()
	case "netcat":
		outLine = kill(args)
		// Если команда не соответствует ни одной из встроенных команд (cd, pwd, echo, kill, ps, netcat),
		//считаем, что это внешняя команда, и выполняем ее.
		//Результат выполнения команды сохраняется в outLine
	default:
		cmd := exec.Command(commands[0], args...)

		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		if isPipe {
			cmd.Stdout = &previousOut
		}

		err := cmd.Run()
		if err != nil {
			return fmt.Sprintf("%v\n", err)
		}
	}
	return outLine

}

var (
	previousOut bytes.Buffer
	isPipe      bool
)

// miniShell - реализует основной цикл оболочки. Она читает ввод пользователя,
//выполняет команды и выводит результаты, пока пользователь не введет команду "quit".

func miniShell() error {
	for {
		line, err := readLine()
		if err != nil {
			return err
		}

		// quit
		if line == "quit" {
			break
		}

		var outLine string

		cmds := strings.Split(line, "|")
		isPipe = true

		for _, cmd := range cmds {
			if cmd == cmds[len(cmds)-1] {
				isPipe = false
			}

			outLine = applyCommand(cmd)
		}

		fmt.Print(outLine)
	}

	return nil
}

func main() {
	err := miniShell()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
