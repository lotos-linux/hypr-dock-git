package ipc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/goccy/go-json"

	"hypr-dock/enternal/pkg/cfg"
)

var monitors []Monitor
var clients []Client
var config cfg.Config

var hyprDir = filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "hypr")
var his = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

var unixSockAdress = fmt.Sprintf("%s/%s/.socket.sock", hyprDir, his)
var unixSock2Adress = &net.UnixAddr{
	Name: fmt.Sprintf("%s/%s/.socket2.sock", hyprDir, his),
	Net:  "unix",
}

type EventCallback struct {
	Event   string
	Handler func(string)
}

var EventCallBacks []EventCallback
var mu sync.Mutex

func NewEventHandler(event string, callBack func(string)) {
	mu.Lock()
	defer mu.Unlock()

	EventCallBacks = append(EventCallBacks, EventCallback{
		Event:   event,
		Handler: callBack,
	})
}

func Hyprctl(cmd string) ([]byte, error) {
	conn, err := net.Dial("unix", unixSockAdress)
	if err != nil {
		return nil, err
	}

	message := []byte(cmd)
	_, err = conn.Write(message)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 102400)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return response[:n], nil
}

func ListMonitors() error {
	response, err := Hyprctl("j/monitors")
	if err != nil {
		return err
	} else {
		err = json.Unmarshal([]byte(response), &monitors)
	}
	return err
}

func ListClients() ([]Client, error) {
	response, err := Hyprctl("j/clients")
	// fmt.Println(response)
	if err != nil {
		return nil, err
	} else {
		err = json.Unmarshal([]byte(response), &clients)
		// fmt.Println("json reading...")
	}
	return clients, err
}

func GetActiveWindow() (*Client, error) {
	var activeWindow Client
	response, err := Hyprctl("j/activewindow")
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal([]byte(response), &activeWindow)
	if err == nil {
		return &activeWindow, nil
	}
	return nil, err
}

func InitHyprEvents() {
	unixConnect, _ := net.DialUnix("unix", nil, unixSock2Adress)
	defer unixConnect.Close()
	for {
		bufer := make([]byte, 10240)
		unixNumber, err := unixConnect.Read(bufer)
		if err != nil {
			fmt.Println(err)
		}
		hyprEvent := string(bufer[:unixNumber])
		events := splitEvents(hyprEvent)

		for _, event := range events {
			// fmt.Println(event)

			for _, eventCallBack := range EventCallBacks {
				if strings.Contains(event, eventCallBack.Event) {
					eventCallBack.Handler(event)
				}
			}

			// if strings.Contains(event, "configreloaded") {
			// 	go addLayerRule()
			// }

			// if strings.Contains(event, "windowtitlev2>>") {
			// 	go windowTitleHandler(event)
			// }

			// if strings.Contains(event, "openwindow>>") {
			// 	go openwindowHandler(event)
			// }

			// if strings.Contains(event, "closewindow>>") {
			// 	go closewindowHandler(event)
			// }

			// if strings.Contains(event, "activespecial>>") {
			// 	go activatespecialHandler(event)
			// }
		}
	}
}

func SearchClientByAddress(address string) (Client, error) {
	ListClients()

	for _, ipcClient := range clients {
		if ipcClient.Address == address {
			return ipcClient, nil
		}
	}

	err := errors.New("Client non found by address: " + address)
	return Client{}, err
}

func AddLayerRule() {
	if config.Blur == "on" {
		Hyprctl("keyword layerrule blur,hypr-dock")
		Hyprctl("keyword layerrule ignorealpha 0.4,hypr-dock")
	}
}

func splitEvents(multiLineEvent string) []string {
	events := strings.Split(multiLineEvent, "\n")

	var filteredEvents []string
	for _, event := range events {
		event = strings.TrimSpace(event)
		if event != "" {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents
}

func SetConfig(inpConfig cfg.Config) {
	config = inpConfig
}
