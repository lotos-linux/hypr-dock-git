LOCAL_CONFIG_DIR = $(HOME)/.config/hypr-dock

PROJECT_BIN_DIR = bin
PROJECT_CONFIG_DIR = configs

EXECUTABLE = hypr-dock

install:
		sudo cp $(PROJECT_BIN_DIR)/$(EXECUTABLE) /usr/bin/

		mkdir -p $(LOCAL_CONFIG_DIR)
		cp -r $(PROJECT_CONFIG_DIR)/* $(LOCAL_CONFIG_DIR)/

		@echo "Installation completed."

uninstall:
		sudo rm -f /usr/bin/$(EXECUTABLE)

		rm -rf $(LOCAL_CONFIG_DIR)

		@echo "Installation removed."

get:
		go mod tidy

build:
		go build -v -o bin/hypr-dock ./main/.

run:
		go run .

exec:
		./bin/hypr-dock