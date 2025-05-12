package main

import (
	"flag"
	"fmt"
	"infra-health-cli/methods"
	"infra-health-cli/output"
	"os"
	"os/exec"
	"runtime"
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
func main() {
	flagChoice := flag.Int("choice", 0, "Monitoring type: 1=HTTP, 2=HTTPS, 3=TELNET, 4=ICMP")
	flagEndpoint := flag.String("endpoint", "", "Target endpoint")
	flagPort := flag.Int("port", 0, "TCP port (only for choice 3)")
	flagJSON := flag.Bool("json", false, "Output result as JSON")
	flag.Parse()

	if *flagChoice > 0 && *flagEndpoint != "" && (*flagChoice != 3 || *flagPort > 0) {
		runMonitor(*flagChoice, *flagEndpoint, *flagPort, *flagJSON, false)
		return
	}

	for {
		var choice, portnumber int
		var endpoint string

		fmt.Printf("Select type of you want to monitor:\n1. HTTP\n2. HTTPS\n3. TELNET\n4. ICMP\n")
		fmt.Scan(&choice)

		fmt.Printf("Enter your endpoint: ")
		fmt.Scan(&endpoint)

		if choice == 3 {
			fmt.Printf("Enter port number: ")
			fmt.Scan(&portnumber)
		}

		runMonitor(choice, endpoint, portnumber, *flagJSON, true)

		fmt.Println("\n--- Press Enter to continue ---")
		fmt.Scanln()
		clearScreen()
	}
}

func runMonitor(choice int, endpoint string, portnumber int, jsonOutput bool, interactive bool) {

	switch choice {

	case 1:
		methods.HTTPCHECK(choice, endpoint, portnumber, jsonOutput, interactive)

	case 2:
		methods.HTTPSCHECK(choice, endpoint, portnumber, jsonOutput, interactive)

	case 3:
		methods.TELNET(choice, endpoint, portnumber, jsonOutput, interactive)

	case 4:
		methods.ICMPER(endpoint, jsonOutput, interactive)

	default:
		result := output.NewInvalidChoiceResult()
		if !jsonOutput {
			fmt.Println("Invalid choice.")
		}
		output.Jsonify(result, jsonOutput)
		return
	}
}
