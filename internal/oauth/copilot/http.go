package copilot

import "strings"

const (
	userAgent           = "GitHubCopilotChat/0.32.4"
	editorVersion       = "vscode/1.105.1"
	editorPluginVersion = "copilot-chat/0.32.4"
	integrationID       = "vscode-chat"
)

func Headers() map[string]string {
	return map[string]string{
		"User-Agent":             userAgent,
		"Editor-Version":         editorVersion,
		"Editor-Plugin-Version":  editorPluginVersion,
		"Copilot-Integration-Id": integrationID,
	}
}

// ParseEndpointFromToken extracts the API endpoint from a Copilot token.
// The token format is: "tid=xxx;exp=xxx;proxy-ep=proxy.individual.githubcopilot.com;..."
// Returns the HTTPS API endpoint URL, or empty string if not found.
func ParseEndpointFromToken(token string) string {
	// Token fields are separated by semicolons
	parts := strings.Split(token, ";")
	for _, part := range parts {
		if strings.HasPrefix(part, "proxy-ep=") {
			proxyEP := strings.TrimPrefix(part, "proxy-ep=")
			// Convert proxy endpoint to API endpoint
			// proxy.individual.githubcopilot.com -> api.individual.githubcopilot.com
			if strings.HasPrefix(proxyEP, "proxy.") {
				apiEP := "api." + strings.TrimPrefix(proxyEP, "proxy.")
				return "https://" + apiEP
			}
			return "https://" + proxyEP
		}
	}
	return ""
}
