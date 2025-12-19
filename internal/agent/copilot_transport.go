package agent

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
)

// CopilotHeaderTransport is an http.RoundTripper that injects VSCode-compatible
// headers for Copilot request grouping. These headers tell GitHub Copilot to
// treat multiple agent turns within a single interaction as one billing unit (0.1 quota).
type CopilotHeaderTransport struct {
	Transport http.RoundTripper

	// interactionID is reused for all requests within one user prompt.
	// This groups multiple agent turns into a single billable interaction.
	interactionID string
	mu            sync.RWMutex
}

// NewCopilotHeaderTransport creates a new transport that wraps the given transport
// and injects Copilot-specific headers.
func NewCopilotHeaderTransport(transport http.RoundTripper) *CopilotHeaderTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}
	return &CopilotHeaderTransport{
		Transport:     transport,
		interactionID: uuid.NewString(),
	}
}

// NewInteraction generates a new interaction ID. Call this at the start of each
// user prompt to group all subsequent requests under the same billing unit.
func (c *CopilotHeaderTransport) NewInteraction() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.interactionID = uuid.NewString()
	return c.interactionID
}

// RoundTrip implements http.RoundTripper, injecting Copilot headers.
func (c *CopilotHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid mutating the original
	reqClone := req.Clone(req.Context())

	c.mu.RLock()
	interactionID := c.interactionID
	c.mu.RUnlock()

	// These headers tell Copilot to group requests as a single interaction
	// Matching VSCode Copilot Chat behavior for 0.1 quota consumption
	reqClone.Header.Set("x-interaction-id", interactionID)
	reqClone.Header.Set("x-interaction-type", "conversation-agent")
	reqClone.Header.Set("openai-intent", "conversation-agent")
	reqClone.Header.Set("x-initiator", "agent")
	reqClone.Header.Set("x-request-id", uuid.NewString())

	return c.Transport.RoundTrip(reqClone)
}
