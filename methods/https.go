package methods

func HTTPSCHECK(choice int, endpoint string, portnumber int, jsonOutput bool, interactive bool) int {
	return doHTTPCheck("https", "HTTPS", endpoint, portnumber, jsonOutput, interactive)
}
