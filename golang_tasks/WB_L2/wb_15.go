package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec" // удобная либа для использования терминала
	"strconv"
	"strings"
)

func main() {
	// Заполняем нужные команды
	cmd := map[string]bool{
		"cd":   true,
		"ps":   true,
		"pwd":  true,
		"echo": true,
		"kill": true,
	}

	currPath, _ := os.Getwd()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(currPath, "$ ")
		line, err := reader.ReadString('\n') // читаем до \n
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading:", err)
			break
		}

		// убираем лишние пробелы
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// делим запрос на части
		commands := strings.Fields(line)
		if len(commands) == 0 {
			continue
		}

		//провермяем есть ли команда в мапе
		if !cmd[commands[0]] {
			// Запуск внешней команды
			runExternal(commands[0], commands[1:])
			continue
		}

		switch commands[0] {
		case "cd":
			if len(commands) < 2 {
				home, _ := os.UserHomeDir()
				commands = append(commands, home)
			}
			newPath := cd(commands[1])
			if newPath != "" {
				currPath = newPath
			}
		case "pwd":
			fmt.Println(pwd())
		case "echo":
			echo(commands[1:])
		case "ps":
			ps()
		case "kill":
			if len(commands) < 2 {
				fmt.Println("kill: not enough arguments")
				continue
			}
			kill(commands[1])
		}
	}
}

func cd(path string) string {
	err := os.Chdir(path)
	if err != nil {
		fmt.Println("cd:", err)
		return ""
	}
	newPath, _ := os.Getwd()
	return newPath
}

func pwd() string {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("pwd:", err)
		return ""
	}
	return path
}

func echo(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func ps() {
	cmd := exec.Command("ps")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func kill(pidStr string) {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		fmt.Println("kill: invalid pid")
		return
	}
	// проверяем наличие процесса перед его закрытием
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("kill:", err)
		return
	}
	err = process.Kill()
	if err != nil {
		fmt.Println("kill:", err)
	}
}

func runExternal(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s: %v\n", name, err)
	}
}
