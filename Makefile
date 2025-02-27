LOCAL_BIN_DIR = $(HOME)/.local/bin
LOCAL_CONFIG_DIR = $(HOME)/.config/hypr-dock

PROJECT_BIN_DIR = bin
PROJECT_CONFIG_DIR = config

EXECUTABLE = hypr-dock

install:
    @echo "Installing binary to $(LOCAL_BIN_DIR)..."
    mkdir -p $(LOCAL_BIN_DIR)
    cp $(PROJECT_BIN_DIR)/$(EXECUTABLE) $(LOCAL_BIN_DIR)/

    @echo "Copying configuration to $(LOCAL_CONFIG_DIR)..."
    mkdir -p $(LOCAL_CONFIG_DIR)
    cp -r $(PROJECT_CONFIG_DIR)/* $(LOCAL_CONFIG_DIR)/

    @echo "Local installation completed."

uninstall:
    @echo "Removing binary from $(LOCAL_BIN_DIR)..."
    rm -f $(LOCAL_BIN_DIR)/$(EXECUTABLE)

    @echo "Removing configuration from $(LOCAL_CONFIG_DIR)..."
    rm -rf $(LOCAL_CONFIG_DIR)

    @echo "Local installation removed."

get:
		go mod tidy

build:
		go build -v -o bin/hypr-dock ./main/.

run:
		go run .

exec:
		./bin/hypr-dock