package methods

import (
	"fmt"
	"infra-health-cli/misc"
	"infra-health-cli/output"
	"net"
	"strconv"
	"time"
)

func TELNET(choice int, endpoint string, portnumber int, jsonOutput bool, interactive bool) int {

	result := output.MonitorResult{
		Type:      "TELNET",
		Endpoint:  endpoint,
		Port:      portnumber,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var address string

	if portnumber > 0 {
		address = endpoint + ":" + strconv.Itoa(portnumber)
	} else {
		address = endpoint
	}

	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, time.Minute)

	result.Latency = misc.TrackLatency(start)

	if err != nil {
		result.Status = "Closed"
		if !jsonOutput {
			fmt.Printf("Port %d is closed on %s\n", portnumber, endpoint)
			fmt.Println(err)
		}
		output.Jsonify(result, jsonOutput)
		return 1
	}
	defer conn.Close()
	if conn != nil {
		result.Status = "Open"
		if !jsonOutput {
			fmt.Printf("Port %d is open on %s\n", portnumber, endpoint)
		}
		output.Jsonify(result, jsonOutput)
		return 0
	}
	return 2
}
