package utils

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Launch(command string) {
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

	log.Printf("env vars: %s; command: '%s'; args: %s\n", envVars, elements[cmdIdx], args)

	if err := cmd.Start(); err != nil {
		log.Println("Unable to launch command!", err.Error())
	}
}

func LoadTextFile(path string) ([]string, error) {
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

func TempDir() string {
	if os.Getenv("TMPDIR") != "" {
		return os.Getenv("TMPDIR")
	} else if os.Getenv("TEMP") != "" {
		return os.Getenv("TEMP")
	} else if os.Getenv("TMP") != "" {
		return os.Getenv("TMP")
	}
	return "/tmp"
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
