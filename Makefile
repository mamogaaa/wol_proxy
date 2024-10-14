APP_NAME=wol_proxy
INSTALL_DIR=/usr/local/bin
CONFIG_DIR=/etc/wol_proxy
SERVICE_DIR=/etc/systemd/system

# Build the Go application
build:
	go build -o $(APP_NAME) .

# Install the binary, config, and systemd service
install: build
	# Create the directories for configuration and logs
	mkdir -p $(CONFIG_DIR)
	mkdir -p /var/log/$(APP_NAME)
	# Install the binary
	install -m 0755 $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
	# Install the configuration file
	install -m 0644 config.yaml $(CONFIG_DIR)/config.yaml
	# Install the systemd service file
	install -m 0644 $(APP_NAME).service $(SERVICE_DIR)/$(APP_NAME).service
	# Reload systemd to recognize the new service
	systemctl daemon-reload
	# Enable and start the service
	systemctl enable $(APP_NAME).service
	systemctl start $(APP_NAME).service

# Uninstall the application and systemd service
uninstall:
	# Stop the service
	systemctl stop $(APP_NAME).service
	# Disable the service
	systemctl disable $(APP_NAME).service
	# Remove the binary, configuration, and service file
	rm -f $(INSTALL_DIR)/$(APP_NAME)
	rm -rf $(CONFIG_DIR)
	rm -f $(SERVICE_DIR)/$(APP_NAME).service
	# Reload systemd to remove the service
	systemctl daemon-reload

# Clean build files
clean:
	rm -f $(APP_NAME)

.PHONY: build install uninstall clean
