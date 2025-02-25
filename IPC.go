package main

import (
	"os"
	"fmt"
	"net"
	"strings"
	"errors"
	"path/filepath"
	"github.com/goccy/go-json"
	// "github.com/dlasky/gotk3-layershell/layershell"
	"github.com/gotk3/gotk3/glib"
)

type workspace struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Monitor         string `json:"monitor"`
	Windows         int    `json:"windows"`
	Hasfullscreen   bool   `json:"hasfullscreen"`
	Lastwindow      string `json:"lastwindow"`
	Lastwindowtitle string `json:"lastwindowtitle"`
}

type monitor struct {
	Id              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Make            string  `json:"make"`
	Model           string  `json:"model"`
	Serial          string  `json:"serial"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	RefreshRate     float64 `json:"refreshRate"`
	X               int     `json:"x"`
	Y               int     `json:"y"`

	ActiveWorkspace struct {
		Id   int    `json:"id"`
		Name string `json:"name"`

	} `json:"activeWorkspace"`

	Reserved   []int   `json:"reserved"`
	Scale      float64 `json:"scale"`
	Transform  int     `json:"transform"`
	Focused    bool    `json:"focused"`
	DpmsStatus bool    `json:"dpmsStatus"`
	Vrr        bool    `json:"vrr"`
}

type client struct {
	Address         string `json:"address"`
	Mapped          bool   `json:"mapped"`
	Hidden          bool   `json:"hidden"`
	At              []int  `json:"at"`
	Size            []int  `json:"size"`

	Workspace struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"workspace"`

	Floating        bool          `json:"floating"`
	Pseudo          bool          `json:"pseudo"`          // Добавлено
	Monitor         int           `json:"monitor"`
	Class           string        `json:"class"`
	Title           string        `json:"title"`
	InitialClass    string        `json:"initialClass"`
	InitialTitle    string        `json:"initialTitle"`
	Pid             int           `json:"pid"`
	Xwayland        bool          `json:"xwayland"`
	Pinned          bool          `json:"pinned"`
	Fullscreen      int           `json:"fullscreen"`      // Исправлено: int вместо bool
	FullscreenClient int          `json:"fullscreenClient"` // Добавлено
	Grouped         []interface{} `json:"grouped"`
	Tags            []interface{} `json:"tags"`            // Добавлено
	Swallowing      string        `json:"swallowing"`      // Исправлено: string вместо interface{}
	FocusHistoryID  int           `json:"focusHistoryID"`  // Добавлено
	InhibitingIdle  bool          `json:"inhibitingIdle"`  // Добавлено
}

var monitors                           []monitor
var clients                            []client
var activeClient                       *client
var lastWinAddr                        string

var hyprDir = filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "hypr")
var his = os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")

var unixSockAdress = fmt.Sprintf("%s/%s/.socket.sock", hyprDir, his)
var unixSock2Adress = &net.UnixAddr {
	Name: fmt.Sprintf("%s/%s/.socket2.sock", hyprDir, his),
	Net:  "unix",
}

var special = false

func hyprctl(cmd string) ([]byte, error) {
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

func listMonitors() error {
	response, err := hyprctl("j/monitors")
	if err != nil {
		return err
	} else {
		err = json.Unmarshal([]byte(response), &monitors)
	}
	return err
}

func listClients() error {
	response, err := hyprctl("j/clients")
	// fmt.Println(response)
	if err != nil {
		return err
	} else {
		err = json.Unmarshal([]byte(response), &clients)
		// fmt.Println("json reading...")
	}
	activeClient, _ = getActiveWindow()
	// fmt.Println(activeClient)
	return err
}

func getActiveWindow() (*client, error) {
	var activeWindow client
	response, err := hyprctl("j/activewindow")
	err = json.Unmarshal([]byte(response), &activeWindow)
	if err == nil {
		return &activeWindow, nil
	}
	return nil, err
}

func initHyprEvents() {
	unixConnect, _ := net.DialUnix("unix", nil, unixSock2Adress)
	defer unixConnect.Close()
	for {
		bufer := make([]byte, 10240)
		unixNumber, err := unixConnect.Read(bufer)
		if err != nil {fmt.Println(err)}
		hyprEvent := string(bufer[:unixNumber])
		events := splitEvents(hyprEvent)

		for _, event := range events {
			// fmt.Println(event)

			if strings.Contains(event, "configreloaded") {
				go addLayerRule()
			}

			if strings.Contains(event, "windowtitlev2>>") {
				go windowTitleHandler(event)
			}

			if strings.Contains(event, "openwindow>>") {
				go openwindowHandler(event)
			}

			if strings.Contains(event, "closewindow>>") {
				go closewindowHandler(event)
			}

			if strings.Contains(event, "activespecial>>") {
				go activatespecialHandler(event)
			}
		}
	}
}

func windowTitleHandler(event string) {
	data := eventHandler(event, 2)
	address := "0x" + strings.TrimSpace(data[0])
	go changeWindowTitle(address, data[1])
}

func activatespecialHandler(event string) {
	data := eventHandler(event, 2)

	if data[0] == "special:special" {
		special = true
	}
	if data[0] != "special:special" {
		special = false
	}
}

func openwindowHandler(event string) {
	data := eventHandler(event, 4)
	address := "0x" + strings.TrimSpace(data[0])
	windowClient, err := searchClientByAddress(address)
	if err != nil {
		fmt.Println(err)
	} else {
		glib.IdleAdd(func() {
			addApp(windowClient)
		})
	}
}

func closewindowHandler(event string) {
	data := eventHandler(event, 1)
	address := "0x" + strings.TrimSpace(data[0])
	glib.IdleAdd(func() {
		removeApp(address)
	})
}

func eventHandler(event string, n int) []string {
	parts := strings.SplitN(event, ">>", 2)
	dataClast := strings.TrimSpace(parts[1])
	dataParts := strings.SplitN(dataClast, ",", n)

	for i := range dataParts {
		dataParts[i] = strings.TrimSpace(dataParts[i])
	}

	return dataParts
}

func changeWindowTitle(address string, title string) {
	mu.Lock()
	defer mu.Unlock()

	for _, data := range addedApps {
		for _, appWindow := range data.Windows {
			if appWindow["Address"] == address {
				appWindow["Title"] = title
			}
		}
	}
}

func searchClientByAddress(address string) (client, error) {
	listClients()

	for _, ipcClient := range clients {
		if ipcClient.Address == address {
			return ipcClient, nil
		}
	}

	err := errors.New("Client non found by address: " + address)
	return client{}, err
}

func addLayerRule() {
	if config.Blur == "on" {
		hyprctl("keyword layerrule blur,hypr-dock")
		hyprctl("keyword layerrule ignorealpha 0.4,hypr-dock")
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