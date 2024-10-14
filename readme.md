# Wake-on-LAN Proxy Service

This project provides an HTTP proxy that sends a Wake-on-LAN (WOL) packet to a server if the server is down, waits for it to come online, and then forwards the request. The proxy service is configured through a YAML file and managed as a systemd service.

## Features

- Sends a WOL magic packet to wake up a server.
- Waits for the server to become available and retries if necessary.
- Forwards HTTP requests to the server after it becomes available.
- Configurable via YAML file located in `/etc/wake_on_lan_proxy/config.yaml`.
- Managed as a systemd service, enabling easy startup and monitoring.

## Installation

### Prerequisites

- **Go**: Ensure that Go is installed. You can download and install Go from [here](https://golang.org/dl/).
- **Systemd**: Systemd is required to manage the service.

### Steps to Install

1. Clone the repository and navigate to the project directory:

   ```bash
   git clone https://github.com/yourusername/wake_on_lan_proxy.git
   cd wake_on_lan_proxy
   ```

2. Build the project and install the proxy service and configuration:

   ```bash
   sudo make install
   ```

   This will:
   - Build the Go binary.
   - Copy the binary to `/usr/local/bin/wol_proxy`.
   - Copy the default configuration file to `/etc/wake_on_lan_proxy/config.yaml`.
   - Install the systemd service file to `/etc/systemd/system/wol_proxy.service`.
   - Start and enable the service.

3. Check the service status to ensure it’s running:

   ```bash
   sudo systemctl status wol_proxy
   ```

### Configuration

The configuration file is located at `/etc/wake_on_lan_proxy/config.yaml`. You can edit this file to update the settings for your proxy:

```yaml
listen_port: ":8080"                 # Port to listen on
mac_address: "AA:BB:CC:DD:EE:FF"     # MAC address of the server to wake
server_address: "192.168.1.100:80"   # Server address to proxy
wol_port: 9                          # Port for sending WOL packets (usually 9)
check_interval: 10                   # Time (in seconds) between server availability checks
retry_attempts: 10                   # Number of retry attempts to check server availability after WOL
```

After editing the configuration file, restart the service for changes to take effect:

```bash
sudo systemctl restart wol_proxy
```

### Uninstallation

To uninstall the proxy and systemd service, run:

```bash
sudo make uninstall
```

This will:
- Stop and disable the systemd service.
- Remove the binary, configuration file, and service file.

## Systemd Service

The proxy service is managed by `systemd` for easy startup and monitoring. The service file is located at `/etc/systemd/system/wol_proxy.service`.

### Useful Commands

- **Start the service**:
  
  ```bash
  sudo systemctl start wol_proxy
  ```

- **Stop the service**:
  
  ```bash
  sudo systemctl stop wol_proxy
  ```

- **Check the service status**:
  
  ```bash
  sudo systemctl status wol_proxy
  ```

- **View service logs**:
  
  ```bash
  journalctl -u wol_proxy
  ```

- **Restart the service**:
  
  ```bash
  sudo systemctl restart wol_proxy
  ```

- **Enable the service to start on boot**:

  ```bash
  sudo systemctl enable wol_proxy
  ```

- **Disable the service**:

  ```bash
  sudo systemctl disable wol_proxy
  ```

## Development

If you want to modify or extend this project, you can rebuild it by running:

```bash
make build
```

This will generate a new binary `wol_proxy` in the project directory.

To clean up build files:

```bash
make clean
```

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

## Contributions

Contributions are welcome! Feel free to open issues or submit pull requests to improve this project.

---

Let me know if you need any further clarification or changes!