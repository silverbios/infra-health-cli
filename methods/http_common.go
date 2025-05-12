package methods

import (
	"fmt"
	"infra-health-cli/misc"
	"infra-health-cli/output"
	"net/http"
	"strconv"
	"time"
)

func doHTTPCheck(scheme string, checkType string, endpoint string, portnumber int, jsonOutput bool, interactive bool) int {
	result := output.MonitorResult{
		Type:      checkType,
		Endpoint:  endpoint,
		Port:      portnumber,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	address := endpoint
	if portnumber > 0 {
		address = endpoint + ":" + strconv.Itoa(portnumber)
	}

	start := time.Now()

	resp, err := http.Get(scheme + "://" + address)

	result.Latency = misc.TrackLatency(start)

	if err != nil {
		result.Status = "Unreachable"
		if !jsonOutput {
			fmt.Printf("The %s is not reachable at the moment (%s)\n", endpoint, checkType)
		}
		output.Jsonify(result, jsonOutput)
		return 1
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode != 200 {
		result.Status = "error"
		if !jsonOutput {
			fmt.Printf("Something bad happened on %s (%s)\n", endpoint, checkType)
			fmt.Println(err)
		}
		output.Jsonify(result, jsonOutput)
		return 1
	}

	result.Status = "OK"
	if !jsonOutput {
		fmt.Printf("%s check successful for %s (Status: %d)\n", checkType, endpoint, resp.StatusCode)
	}
	output.Jsonify(result, jsonOutput)
	return 0
}
