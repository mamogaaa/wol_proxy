package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	configFilePath string
)

type Config struct {
	ListenPort    string `yaml:"listen_port"`
	MacAddress    string `yaml:"mac_address"`
	ServerAddress string `yaml:"server_address"`
	WolPort       int    `yaml:"wol_port"`
	CheckInterval int    `yaml:"check_interval"`
	RetryAttempts int    `yaml:"retry_attempts"`
}

var config Config

// wakeOnLan sends a magic packet to wake up the server via WOL
func wakeOnLan(macAddr string) error {
	mac, err := net.ParseMAC(macAddr)
	if err != nil {
		return err
	}

	// Create the magic packet
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:], mac)
	}

	// Send the magic packet to the network's broadcast address
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("255.255.255.255:%d", config.WolPort))
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	return err
}

// isServerUp checks if the target server is up by sending a test request
func isServerUp() bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("http://%s", config.ServerAddress))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// proxyHandler handles the incoming HTTP requests and forwards them to the server
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	if !isServerUp() {
		log.Println("Server is down, sending Wake-on-LAN packet...")
		if err := wakeOnLan(config.MacAddress); err != nil {
			http.Error(w, "Failed to send WOL packet", http.StatusInternalServerError)
			return
		}

		// Wait for the server to become available
		log.Println("Waiting for the server to wake up...")
		for attempts := 0; attempts < config.RetryAttempts; attempts++ {
			if isServerUp() {
				break
			}
			time.Sleep(time.Duration(config.CheckInterval) * time.Second)
		}
		if !isServerUp() {
			http.Error(w, "Server is not responding", http.StatusGatewayTimeout)
			return
		}
		log.Println("Server is back up!")
	}

	// Forward the request to the server
	log.Println("Forwarding request to server...")
	resp, err := http.Get(fmt.Sprintf("http://%s%s", config.ServerAddress, r.URL.Path))
	if err != nil {
		http.Error(w, "Failed to reach server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy the response back to the client
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// loadConfig loads the YAML config file
func loadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Parse command-line arguments
	flag.StringVar(&configFilePath, "config", "/etc/wake_on_lan_proxy/config.yaml", "Path to config file")
	flag.Parse()

	// Load the configuration
	if err := loadConfig(configFilePath); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	http.HandleFunc("/", proxyHandler)
	log.Printf("Starting proxy on %s\n", config.ListenPort)
	log.Printf("Proxying requests to %s\n", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ListenPort, nil))
}
