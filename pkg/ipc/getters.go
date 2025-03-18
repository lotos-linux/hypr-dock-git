package ipc

import (
	"fmt"
	"log"

	"github.com/goccy/go-json"
)

func GetMonitors() ([]Monitor, error) {
	var monitors []Monitor
	response, err := Hyprctl("j/monitors")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(response), &monitors)
	return monitors, err
}

func GetClients() ([]Client, error) {
	var clients []Client
	response, err := Hyprctl("j/clients")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(response), &clients)
	return clients, err
}

func GetActiveWindow() (*Client, error) {
	var activeWindow Client
	response, err := Hyprctl("j/activewindow")
	if err != nil {
		log.Printf("Failed to get active window: %v", err)
		return nil, err
	}
	err = json.Unmarshal([]byte(response), &activeWindow)
	if err != nil {
		log.Printf("Failed to unmarshal active window: %v", err)
		return nil, err
	}
	return &activeWindow, nil
}

func GetOption(option string, v interface{}) error {
	cmd := fmt.Sprintf("j/getoption %s", option)
	response, err := Hyprctl(cmd)
	if err != nil {
		log.Printf("Failed to execute Hyprctl command for option '%s': %v", option, err)
		return err
	}
	err = json.Unmarshal(response, v)
	if err != nil {
		log.Printf("Failed to unmarshal JSON response for option '%s': %v", option, err)
		return err
	}
	return nil
}
