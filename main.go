package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // Linux, macOS
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

type MonitorResult struct {
	Type       string `json:"type"`
	Endpoint   string `json:"endpoint"`
	Port       int    `json:"port,omitempty"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Latency    string `json:"latency,omitempty"`
	PacketLoss string `json:"packet_loss,omitempty"`
}

func main() {
	flagChoice := flag.Int("choice", 0, "Monitoring type: 1=HTTP, 2=HTTPS, etc.")
	flagEndpoint := flag.String("endpoint", "", "Target endpoint")
	flagPort := flag.Int("port", 0, "TCP port (only for choice 3)")
	flagJSON := flag.Bool("json", false, "Output result as JSON")
	flag.Parse()

	if *flagChoice > 0 && *flagEndpoint != "" && (*flagChoice != 3 || *flagPort > 0) {
		runMonitor(*flagChoice, *flagEndpoint, *flagPort, *flagJSON)
		return
	}

	for {
		var choice, portnumber int
		var endpoint string

		fmt.Printf("Select type of you want to monitor:\n1. HTTP\n2. HTTPS\n3. TCP\n4. ICMP\n")
		fmt.Scan(&choice)

		fmt.Printf("Enter your endpoint: ")
		fmt.Scan(&endpoint)

		if choice == 3 {
			fmt.Printf("Enter port number: ")
			fmt.Scan(&portnumber)
		}

		runMonitor(choice, endpoint, portnumber, *flagJSON)
	}
}

func runMonitor(choice int, endpoint string, portnumber int, jsonOutput bool) {

	clearScreen()

	result := MonitorResult{
		Endpoint: endpoint,
		Port:     portnumber,
	}

	if jsonOutput {
		defer func() {
			jsonData, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Println("JSON marshal error:", err)
			} else {
				fmt.Println(string(jsonData))
			}
		}()
	}

	switch choice {
	case 1:
		resp, err := http.Get("http://" + endpoint)
		result.Type = "HTTP"
		if err != nil {
			result.Status = "Unreachable"
			if !jsonOutput {
				fmt.Printf("the %s is not reachable at the moment\n", endpoint)
			}
			return
		}
		defer resp.Body.Close()
		statusCode := resp.StatusCode
		if statusCode != 200 {
			result.Status = "error"
			if !jsonOutput {
				fmt.Printf("something bad happened on %s\n", endpoint)
			}
		} else {
			result.Status = "OK"
			if !jsonOutput {
				fmt.Printf("HTTP check successful for %s (Status: %d)\n", endpoint, resp.StatusCode)
			}

		}

	case 2:
		resp, err := http.Get("https://" + endpoint)
		result.Type = "HTTPS"
		if err != nil {
			result.Status = "Unreachable"
			if !jsonOutput {
				fmt.Printf("the %s is not reachable at the moment\n", endpoint)
			}
			return
		}
		defer resp.Body.Close()
		statusCode := resp.StatusCode
		if statusCode != 200 {
			result.Status = "error"
			if !jsonOutput {
				fmt.Printf("something bad happened on %s\n", endpoint)
			}
		} else {
			result.Status = "OK"
			if !jsonOutput {
				fmt.Printf("HTTPS check successful for %s (Status: %d)\n", endpoint, resp.StatusCode)
			}
		}

	case 3:
		address := endpoint + ":" + strconv.Itoa(portnumber)
		conn, err := net.DialTimeout("tcp", address, time.Minute)
		if err != nil {
			result.Status = "Closed"
			if !jsonOutput {
				fmt.Printf("Port %d is closed on %s\n", portnumber, endpoint)
			}
			return
		}
		defer conn.Close()
		if conn != nil {
			result.Status = "Open"
			if !jsonOutput {
				fmt.Printf("Port %d is open on %s\n", portnumber, endpoint)
			}
			return
		}

	case 4:
		icmper, err := probing.NewPinger(endpoint)
		icmper.Count = 4
		icmper.Timeout = 5 * time.Second
		err = icmper.Run()
		if err != nil {
			result.Status = "unreachable"
			if !jsonOutput {
				fmt.Printf("The %s is not reachable via ICMP\n", endpoint)
			}
			break
		}

		if err != nil {
			result.Status = "Ping Setup Failed"
			if !jsonOutput {
				fmt.Printf("the %s is not reacble at the moment\n", endpoint)
			}
			return
		}
		plratio := strconv.FormatFloat(icmper.Statistics().PacketLoss, 'f', 2, 64)
		if icmper.Statistics().PacketLoss != 0 {
			result.Status = "Partial Loss"
			if !jsonOutput {
				fmt.Printf("Connectivity issues detected for %s — loss ratio: %s, average latency: %s ms\n", endpoint, plratio, icmper.Statistics().AvgRtt)
			}
		} else {
			result.Status = "OK"
			if !jsonOutput {
				fmt.Printf("Successfully reached %s with an average latency of %s ms\n", endpoint, icmper.Statistics().AvgRtt)
			}
		}

	default:
		result.Type = "unknown"
		result.Status = "invalid choice"
		if !jsonOutput {
			fmt.Println("Invalid choice.")
		}

		if !jsonOutput {
			fmt.Println("\n--- Press Enter to return ---")
			fmt.Scanln()
			fmt.Scanln()
		}
	}
}
