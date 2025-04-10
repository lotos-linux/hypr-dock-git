Давай разберем концептуальный состав, функционал и примерные концепции реализации для cli утилиты. Утилита будет называться hypr-dock-ctl. Ее главная и единственная задача это отправлять команды в сокет собирая их из аргументов вызова. Например запрос вида hypr-dock-ctl layer get —json утилита превращает в строку j/layer get; запрос hypr-dock-ctl layer set <данные которые передаются следующими аргументами> превращает в строку “layer set данные которые передаются следующими аргументами”. Эти строки должны передаваться в сокет и получать ответ сервера. Если произошла ошибка код завершения берётся из ответа сервера, если ошибок нет, то код 0

# **Концепция утилиты hypr-dock-ctl**

## **1. Назначение и функционал**
Утилита `hypr-dock-ctl` служит универсальным клиентом для взаимодействия с демоном через IPC-сокет. Её ключевые функции:

1. **Трансляция команд**  
   Преобразует аргументы CLI в строку для сокета:
   ```
   hypr-dock-ctl <команда> <действие> [данные] [--json]
   → 
   j/<команда> <действие> [данные]  # с --json
   <команда> <действие> [данные]    # без --json
   ```

2. **Передача запросов**  
   Отправляет сформированную строку в Unix-domain сокет (`/tmp/hypr-dock.sock`).

3. **Обработка ответов**  
   - Выводит сырой ответ сервера (текст/JSON)
   - Возвращает код завершения:
     - `0` при успехе
     - Код ошибки из ответа (если есть)
     - `1` при ошибках соединения

---

## **2. Формат команд**

### **2.1 Структура вызова**
```
hypr-dock-ctl <command> <action> [data...] [--json]
```

| Часть      | Описание                          | Обязательность |
|------------|-----------------------------------|----------------|
| `command`  | Основная команда (например `layer`) | Да             |
| `action`   | Действие (`get`, `set`, `list`)   | Да             |
| `data`     | Аргументы через пробел            | Нет            |
| `--json`   | Флаг JSON-формата                 | Нет            |

### **2.2 Примеры**
```bash
# Текстовый запрос
hypr-dock-ctl layer get

# JSON-запрос
hypr-dock-ctl window list --json

# Команда с данными
hypr-dock-ctl config set theme dark

# Команда с составными данными
hypr-dock-ctl notification create "Hello" "Message text"
```

---

## **3. Обработка ответов**

### **3.1 Успешный ответ**
- **Текстовый режим**: выводит данные как есть
  ```
  workspace: main
  windows: 3
  ```
- **JSON-режим**: выводит сырой JSON
  ```json
  {"workspace":"main","windows":3}
  ```
- **Код завершения**: `0`

### **3.2 Ошибка**
- **Текстовый режим**:
  ```
  error: Invalid command [2]
  ```
- **JSON-режим**:
  ```json
  {"error":"Invalid command","code":2}
  ```
- **Код завершения**: извлекается из ответа (`[2]` или `"code":2`)

### **3.3 Ошибки соединения**
- Вывод: `error: Connection failed`
- Код: `1`

---

## **4. Концепция реализации (Go)**

### **4.1 Основные компоненты**
```go
package main

import (
  "fmt"
  "net"
  "os"
  "strings"
)

func main() {
  // 1. Парсинг аргументов
  args, useJSON := parseArgs(os.Args[1:])
  
  // 2. Формирование запроса
  request := buildRequest(args, useJSON)
  
  // 3. Отправка в сокет
  response, err := sendToSocket(request, "/tmp/hypr-dock.sock")
  if err != nil {
    fmt.Println("error: Connection failed")
    os.Exit(1)
  }
  
  // 4. Обработка ответа
  exitCode := handleResponse(response, useJSON)
  os.Exit(exitCode)
}
```

### **4.2 Ключевые функции**
1. **Парсинг аргументов**:
   ```go
   func parseArgs(args []string) (cmdParts []string, useJSON bool) {
     for i, arg := range args {
       if arg == "--json" {
         return append(args[:i], args[i+1:]...), true
       }
     }
     return args, false
   }
   ```

2. **Построение запроса**:
   ```go
   func buildRequest(parts []string, jsonFlag bool) string {
     req := strings.Join(parts, " ")
     if jsonFlag {
       return "j/" + req
     }
     return req
   }
   ```

3. **Отправка в сокет**:
   ```go
   func sendToSocket(request, socketPath string) (string, error) {
     conn, err := net.Dial("unix", socketPath)
     if err != nil {
       return "", err
     }
     defer conn.Close()
     
     _, err = conn.Write([]byte(request + "\n"))
     if err != nil {
       return "", err
     }
     
     buf := make([]byte, 4096)
     n, err := conn.Read(buf)
     return string(buf[:n]), err
   }
   ```

4. **Обработка ответа**:
   ```go
   func handleResponse(response string, isJSON bool) int {
     if isJSON {
       var resp struct{ Code int `json:"code"` }
       if json.Unmarshal([]byte(response), &resp) == nil && resp.Code != 0 {
         fmt.Println(response)
         return resp.Code
       }
     } else if strings.HasPrefix(response, "error: ") {
       parts := strings.Split(response, "[")
       if len(parts) > 1 {
         if code, err := strconv.Atoi(strings.TrimSuffix(parts[1], "]")); err == nil {
           fmt.Println(response)
           return code
         }
       }
     }
     
     fmt.Print(response)
     return 0
   }
   ```

---

## **5. Примеры работы**

### **Успешный запрос**
```bash
$ hypr-dock-ctl workspace get
workspace: main
active: true

$ echo $?
0
```

### **Ошибка**
```bash
$ hypr-dock-ctl invalid cmd
error: Unknown command [2]

$ echo $?
2
```

### **JSON-запрос**
```bash
$ hypr-dock-ctl window list --json
{"windows":["term","browser"]}

$ echo $?
0
```

---

## **6. Преимущества архитектуры**
1. **Минимализм**: делает ровно одну задачу (передача команд)
2. **Гибкость**: поддерживает любые команды и данные
3. **Совместимость**: 
   - Работает с существующими bash-скриптами
   - Поддерживает JSON для сложных клиентов
4. **Прозрачность**: выводит сырой ответ сервера без изменений

Такой дизайн соответствует Unix-философии "делать одну вещь и делать её хорошо".