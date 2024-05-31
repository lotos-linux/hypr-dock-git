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
	Address   string `json:"address"`
	Mapped    bool   `json:"mapped"`
	Hidden    bool   `json:"hidden"`
	At        []int  `json:"at"`
	Size      []int  `json:"size"`

	Workspace struct {
		Id   int    `json:"id"`
		Name string `json:"name"`

	} `json:"workspace"`

	Floating       bool          `json:"floating"`
	Monitor        int           `json:"monitor"`
	Class          string        `json:"class"`
	Title          string        `json:"title"`
	InitialClass   string        `json:"initialClass"`
	InitialTitle   string        `json:"initialTitle"`
	Pid            int           `json:"pid"`
	Xwayland       bool          `json:"xwayland"`
	Pinned         bool          `json:"pinned"`
	Fullscreen     bool          `json:"fullscreen"`
	FullscreenMode int           `json:"fullscreenMode"`
	FakeFullscreen bool          `json:"fakeFullscreen"`
	Grouped        []interface{} `json:"grouped"`
	Swallowing     interface{}   `json:"swallowing"`
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
	if err != nil {
		return err
	} else {
		err = json.Unmarshal([]byte(response), &clients)
	}
	activeClient, _ = getActiveWindow()
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
		// fmt.Println(hyprEvent) 

		if strings.Contains(hyprEvent, "configreloaded") {
			addLayerRule()
		}

		if strings.Contains(hyprEvent, "openwindow>>") {
			windowData := strings.TrimSpace(strings.Split(hyprEvent, "openwindow>>")[1])
			windowAddress := "0x" + strings.Split(windowData, ",")[0]
			windowClient, err := searchClientByAddress(windowAddress)
			if err != nil {
				fmt.Println(err)
			} else {
				go addApp(windowClient)
			}
		}

		if strings.Contains(hyprEvent, "closewindow>>") {
			windowData := strings.TrimSpace(strings.Split(hyprEvent, "closewindow>>")[1])
			windowAddress := "0x" + strings.Split(windowData, ",")[0]
			if err != nil {
				fmt.Println(err)
			} else {
				go removeApp(windowAddress)
			}
		}
	
		if strings.Contains(hyprEvent, "activespecial>>") {
			// fmt.Println(hyprEvent)
			specialData := strings.TrimSpace(strings.Split(hyprEvent, "activespecial>>")[1])
			specialDataArr := strings.Split(specialData, ",")
			if specialDataArr[0] == "special:special" {
				// fmt.Println("Open")
				special = true
			}
			if specialDataArr[0] != "special:special" {
				// fmt.Println("Close")
				special = false
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