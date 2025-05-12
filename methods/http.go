package methods

func HTTPCHECK(choice int, endpoint string, portnumber int, jsonOutput bool, interactive bool) int {
	return doHTTPCheck("http", "HTTP", endpoint, portnumber, jsonOutput, interactive)
}
