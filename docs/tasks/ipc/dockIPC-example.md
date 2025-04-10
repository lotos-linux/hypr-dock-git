# **Модуль dockIPC**

Реализация модуля для работы с IPC через Unix domain sockets:

```go
package dockIPC

import (
	"errors"
	"net"
	"os"
	"syscall"
)

// StartServer - запускает IPC сервер
// fileName: путь к сокету (например "/tmp/hypr-dock.sock")
// handler: функция обработки входящих команд
func StartServer(fileName string, handler func(string) ([]byte, error)) error {
	// Удаляем старый сокет если существует
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}

	// Создаем Unix domain socket
	listener, err := net.Listen("unix", fileName)
	if err != nil {
		return err
	}
	defer listener.Close()

	// Устанавлием права на сокет
	if err := os.Chmod(fileName, 0666); err != nil {
		return err
	}

	// Обработка входящих соединений
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			continue
		}

		go handleConnection(conn, handler)
	}
}

// handleConnection обрабатывает одно соединение
func handleConnection(conn net.Conn, handler func(string) ([]byte, error)) {
	defer conn.Close()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return
	}

	command := string(buf[:n])
	response, err := handler(command)
	if err != nil {
		// Форматируем ошибку в стандартный формат
		response = []byte("error: " + err.Error() + "\n")
	}

	conn.Write(response)
}

// Send отправляет команду в IPC сокет
// command: строка команды (например "j/layer get")
func Send(fileName string, command string) ([]byte, error) {
	conn, err := net.Dial("unix", fileName)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(command + "\n"))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

// StopServer останавливает сервер (вспомогательная функция)
func StopServer(fileName string) error {
	return syscall.Unlink(fileName)
}
```

## **Использование модуля**

### **1. Серверная часть**
```go
package main

import (
	"fmt"
	"dockIPC"
)

func main() {
	// Обработчик команд
	handler := func(command string) ([]byte, error) {
		switch command {
		case "j/layer get":
			return []byte(`{"layers":["bottom","top"]}`), nil
		case "layer get":
			return []byte("bottom\ntop"), nil
		default:
			return nil, fmt.Errorf("unknown command")
		}
	}

	// Запуск сервера
	err := dockIPC.StartServer("/tmp/hypr-dock.sock", handler)
	if err != nil {
		panic(err)
	}
}
```

### **2. Клиентская часть**
```go
package main

import (
	"fmt"
	"dockIPC"
)

func main() {
	// Отправка команды
	response, err := dockIPC.Send("/tmp/hypr-dock.sock", "j/layer get")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", string(response))
}
```

## **Особенности реализации**

1. **Безопасность**:
   - Удаление старого сокета перед созданием
   - Установка прав 0666 на сокет
   - Корректная обработка закрытия соединений

2. **Производительность**:
   - Каждое соединение обрабатывается в отдельной goroutine
   - Буферизированное чтение (4096 байт)

3. **Гибкость**:
   - Обработчик команд может возвращать любые бинарные данные
   - Поддержка текстовых и JSON-команд

4. **Вспомогательные функции**:
   - `StopServer()` для корректного завершения
   - Автоматическое форматирование ошибок

Модуль готов к интеграции в проект `hypr-dock-ctl` и может быть расширен при необходимости.