package output

import (
	"encoding/json"
	"fmt"
	"time"
)

type MonitorResult struct {
	Type       string `json:"type"`
	Endpoint   string `json:"endpoint"`
	Port       int    `json:"port,omitempty"`
	Status     string `json:"status"`
	StatusCode int    `json:"status_code,omitempty"`
	Latency    string `json:"latency,omitempty"`
	PacketLoss string `json:"packet_loss,omitempty"`
	Timestamp  string `json:"timestamp"`
}

func Jsonify(result MonitorResult, jsonOutput bool) {

	if jsonOutput {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			fmt.Println("JSON marshal error:", err)
		} else {
			fmt.Println(string(jsonData))
		}
	}
}

func NewInvalidChoiceResult() MonitorResult {
	return MonitorResult{
		Type:      "unknown",
		Status:    "invalid choice",
		Timestamp: time.Now().Format(time.RFC3339),
	}
}
