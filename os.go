package main

import (
	"os"
	"fmt"
	"slices"
	"syscall"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"github.com/gotk3/gotk3/gtk"
)

type desktopData struct {
	Name 	string
	Icon	string
	Exec	string
}

var desktopDirs = getAppDirs()

func getDesktopData(className string) (desktopData, error) {
	allData, err := loadTextFile(searchDesktopFile(className))
	if err != nil {
		return desktopData{}, err
	}

	data := desktopData{
		Name: getDesktopOption(allData, "Name"),
		Icon: getDesktopOption(allData, "Icon"),
		Exec: getDesktopOption(allData, "Exec"),
	}

	return data, nil
}

func searchDesktopFile(className string) string{
	for _, appDir := range desktopDirs {
		desktopFile := className + ".desktop"
		_, err := os.Stat(filepath.Join(appDir, desktopFile))
		if err == nil {
			return filepath.Join(appDir, desktopFile)
		}
	
		// If file non found
		files, _ := os.ReadDir(appDir)
		for _, file := range files {
			fileName := file.Name()
	
			// "krita" > "org.kde.krita.desktop" / "lutris" > "net.lutris.Lutris.desktop"
			if strings.Count(fileName, ".") > 1 && strings.Contains(fileName, className) {
				return filepath.Join(appDir, fileName)
			}
			// "VirtualBox Manager" > "virtualbox.desktop"
			if fileName == strings.Split(strings.ToLower(className), " ")[0] + ".desktop" {
				return filepath.Join(appDir, fileName)
			}
		}
	}

	return ""
}

func getDesktopOption(allData []string, option string) string {
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

	msg := fmt.Sprintf("env vars: %s; command: '%s'; args: %s\n", envVars, elements[cmdIdx], args)
	fmt.Println(msg)

	if err := cmd.Start(); err != nil {
		fmt.Println("Unable to launch command!", err.Error())
	}
}

func getAppDirs() []string {
	var dirs []string
	xdgDataDirs := ""

	home := os.Getenv("HOME")
	xdgDataHome := os.Getenv("XDG_DATA_HOME")
	if os.Getenv("XDG_DATA_DIRS") != "" {
		xdgDataDirs = os.Getenv("XDG_DATA_DIRS")
	} else {
		xdgDataDirs = "/usr/local/share/:/usr/share/"
	}
	if xdgDataHome != "" {
		dirs = append(dirs, filepath.Join(xdgDataHome, "applications"))
	} else if home != "" {
		dirs = append(dirs, filepath.Join(home, ".local/share/applications"))
	}
	for _, d := range strings.Split(xdgDataDirs, ":") {
		dirs = append(dirs, filepath.Join(d, "applications"))
	}
	flatpakDirs := []string{filepath.Join(home, ".local/share/flatpak/exports/share/applications"),
		"/var/lib/flatpak/exports/share/applications"}

	for _, d := range flatpakDirs {
		if !slices.Contains(dirs, d) {
			dirs = append(dirs, d)
		}
	}
	return dirs
}

func tempDir() string {
	if os.Getenv("TMPDIR") != "" {
		return os.Getenv("TMPDIR")
	} else if os.Getenv("TEMP") != "" {
		return os.Getenv("TEMP")
	} else if os.Getenv("TMP") != "" {
		return os.Getenv("TMP")
	}
	return "/tmp"
}

func signalHandler() {
	signalChanel := make(chan os.Signal, 1)
    signal.Notify(signalChanel, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for {
			signalU := <-signalChanel
			switch signalU {
			case syscall.SIGTERM:
				fmt.Println("Exit... (SIGTERM)")
				gtk.MainQuit()
			case syscall.SIGUSR1:
				fmt.Println("Exit... (SIGUSR1)")
				gtk.MainQuit()
			default:
				fmt.Println("Unknow signal")
			}
		}
	}()
}