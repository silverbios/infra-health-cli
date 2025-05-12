package methods

import (
	"fmt"
	"strconv"
	"time"

	"infra-health-cli/misc"
	"infra-health-cli/output"

	probing "github.com/prometheus-community/pro-bing"
)

func ICMPER(endpoint string, jsonOutput bool, interactive bool) int {

	result := output.MonitorResult{
		Type:      "ICMP",
		Endpoint:  endpoint,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	start := time.Now()

	icmper, err := probing.NewPinger(endpoint)

	result.Latency = misc.TrackLatency(start)

	if err != nil || icmper == nil {
		result.Status = "Ping Setup Failed"
		if !jsonOutput {
			fmt.Printf("Failed to initialize ping for %s\n", endpoint)
		}
		return 2
	}
	icmper.Count = 4
	icmper.Timeout = 5 * time.Second
	//icmper.SetPrivileged(true)
	err = icmper.Run()
	if err != nil {
		fmt.Println("Ping error:", err)
		result.Status = "unreachable"
		if !jsonOutput {
			fmt.Printf("The %s is not reachable via ICMP\n", endpoint)
			fmt.Println(err)
		}
		output.Jsonify(result, jsonOutput)
		return 1
	}
	plratio := strconv.FormatFloat(icmper.Statistics().PacketLoss, 'f', 2, 64)
	if icmper.Statistics().PacketLoss != 0.0 {
		result.Status = "Partial Loss"
		if !jsonOutput {
			fmt.Printf("Connectivity issues detected for %s — loss ratio: %s, average latency: %s ms\n", endpoint, plratio, icmper.Statistics().AvgRtt)
		}
		return 1

	} else {
		result.Status = "OK"
		if !jsonOutput {
			fmt.Printf("Successfully reached %s with an average latency of %s ms\n", endpoint, icmper.Statistics().AvgRtt)
		}
		output.Jsonify(result, jsonOutput)
		return 0
	}
}
