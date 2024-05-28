package main

import (
	"os"
	"fmt"
	"os/exec"
	"strings"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/gdk"
)

type clientData struct {
	Name 	string
	Icon	string
	Exec	string
}

const appDir = "/usr/share/applications/"

func getClientData(className string) (clientData, error) {
	allData, err := loadTextFile(searchDesktopFile(className))
	if err != nil {
		return clientData{}, err
	}

	data := clientData{
		Name: getClientOption(allData, "Name"),
		Icon: getClientOption(allData, "Icon"),
		Exec: getClientOption(allData, "Exec"),
	}

	return data, nil
}

func searchDesktopFile(className string) string{
	desktopFile := className + ".desktop"
	_, err := os.Stat(appDir + desktopFile)
	if err == nil {
		return appDir + desktopFile
	}

	// If file non found
	files, _ := os.ReadDir(appDir)
	for _, file := range files {
		fileName := file.Name()

		// "krita" > "org.kde.krita.desktop" / "lutris" > "net.lutris.Lutris.desktop"
		if strings.Count(fileName, ".") > 1 && strings.Contains(fileName, className) {
			return appDir + fileName
		}
		// "VirtualBox Manager" > "virtualbox.desktop"
		if fileName == strings.Split(strings.ToLower(className), " ")[0] + ".desktop" {
			return appDir + fileName
		}
	}

	return ""
}

func getClientOption(allData []string, option string) string {
	for lineIndex := range len(allData) {
		line := allData[lineIndex]
		if strings.HasPrefix(line, option + "=") {
			optionValue := strings.Split(line, "=")[1]
			return optionValue
		}
	}
	return ""
}

func loadTextFile(path string) ([]string, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(bytes), "\n")
	var output []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			output = append(output, line)
		}

	}
	return output, nil
}

func createImage(source string, size int) *gtk.Image {
	iconTheme, err := gtk.IconThemeGetDefault()
	if err != nil {
		fmt.Println("Unable to icon theme:", err)
	}

	// Create image in file
	if strings.Contains(source, "/") {
		pixbuf, err := gdk.PixbufNewFromFileAtSize(
			source, size, size)
		if err != nil {
			fmt.Println(err)
			pixbuf, _ = iconTheme.LoadIcon(
				"steam", size, gtk.ICON_LOOKUP_FORCE_SIZE)
			
		}

		image, err := gtk.ImageNewFromPixbuf(pixbuf)
		if err != nil {
			fmt.Println(err)
		}
		return image
	}

	// Create image in icon name
	pixbuf, err := iconTheme.LoadIcon(
		source, config.IconSize, gtk.ICON_LOOKUP_FORCE_SIZE)
	if err != nil {
		fmt.Println(err)
	}

	image, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		fmt.Println(err)
	}

	return image
}

func launch(command string) {
	if strings.Contains(command, "\"") {
		command = strings.ReplaceAll(command, "\"", "")
	}

	badArg := strings.Index(command, "%")
	if badArg != -1 {
		command = command[:badArg-1]
	}

	elements := strings.Split(command, " ")

	// find prepended env variables, if any
	envVarsNum := strings.Count(command, "=")
	var envVars []string

	cmdIdx := -1

	if envVarsNum > 0 {
		for idx, item := range elements {
			if strings.Contains(item, "=") {
				envVars = append(envVars, item)
			} else if !strings.HasPrefix(item, "-") && cmdIdx == -1 {
				cmdIdx = idx
			}
		}
	}
	if cmdIdx == -1 {
		cmdIdx = 0
	}
	var args []string
	for _, arg := range elements[1+cmdIdx:] {
		if !strings.Contains(arg, "=") {
			args = append(args, arg)
		}
	}

	cmd := exec.Command(elements[cmdIdx], elements[1+cmdIdx:]...)

	// set env variables
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}

	// msg := fmt.Sprintf("env vars: %s; command: '%s'; args: %s\n", envVars, elements[cmdIdx], args)
	// fmt.Println(msg)

	if err := cmd.Start(); err != nil {
		fmt.Println("Unable to launch command!", err.Error())
	}
}
