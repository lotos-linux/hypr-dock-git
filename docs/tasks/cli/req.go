func ParseRequest(request string) (command string, action string, data string, json bool) {
    request = strings.TrimSpace(request)
    if request == "" {
        return "", "", "", false
    }

    if strings.HasPrefix(request, "j/") {
        json = true
        request = request[2:] 
        request = strings.TrimSpace(request) 
    }

    parts := strings.SplitN(request, " ", 3)
    
    if len(parts) > 0 {
        command = parts[0]
    }
    if len(parts) > 1 {
        action = parts[1]
    }
    if len(parts) > 2 {
        data = parts[2]
    }

    return command, action, data, json
}