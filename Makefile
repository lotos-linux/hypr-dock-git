LOCAL_CONFIG_DIR = $(HOME)/.config/hypr-dock

PROJECT_BIN_DIR = bin
PROJECT_CONFIG_DIR = configs

EXECUTABLE = hypr-dock

install:
		sudo cp $(PROJECT_BIN_DIR)/$(EXECUTABLE) /usr/bin/

		mkdir -p $(LOCAL_CONFIG_DIR)
		cp -r $(PROJECT_CONFIG_DIR)/* $(LOCAL_CONFIG_DIR)/

		@echo -e "\033[32mInstallation completed."

uninstall:
		sudo rm -f /usr/bin/$(EXECUTABLE)

		rm -rf $(LOCAL_CONFIG_DIR)

		@echo -e "\033[32mInstallation removed."

update:
		sudo rm -f /usr/bin/$(EXECUTABLE)
		sudo cp $(PROJECT_BIN_DIR)/$(EXECUTABLE) /usr/bin/

		@echo -e "\033[32mUpdating comleted."

get:
		go mod tidy

build:
		go build -v -o bin/hypr-dock ./main/.

exec:
		./bin/hypr-dock -dev
