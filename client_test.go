package updown

import (
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClient() *Client {
	apiKey := os.Getenv("UPDOWN_API_KEY")
	if apiKey == "" {
		panic("API key is not set. Set UPDOWN_API_KEY environment variable.")
	}
	return NewClient(apiKey, nil)
}

// createTestCheck creates a check for testing and returns its token
func createTestCheck(t *testing.T, client *Client) string {
	res, resp, err := client.Check.Add(CheckItem{
		URL:   "https://example.com",
		Alias: "Test Check",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	return res.Token
}

// deleteTestCheck removes a test check
func deleteTestCheck(t *testing.T, client *Client, token string) {
	_, _, _ = client.Check.Remove(token)
}

func TestTokenForAlias(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	// Cache miss + alias not found
	foundToken, err := client.Check.TokenForAlias("nonexistent-alias-12345")
	assert.Equal(t, "", foundToken)
	assert.Equal(t, ErrTokenNotFound, err)

	// Cache miss + match found after request
	foundToken, err = client.Check.TokenForAlias("Test Check")
	assert.Nil(t, err)
	assert.Equal(t, token, foundToken)

	// Cache hit
	foundToken, err = client.Check.TokenForAlias("Test Check")
	assert.Nil(t, err)
	assert.Equal(t, token, foundToken)
}

func TestList(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	checks, resp, err := client.Check.List()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, len(checks) > 0, "Should have at least one check")

	// Verify checks have expected fields
	for _, check := range checks {
		assert.NotEmpty(t, check.Token)
		assert.NotEmpty(t, check.URL)
	}
}

func TestGet(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	check, resp, err := client.Check.Get(token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Test Check", check.Alias)

	// Test with invalid token
	_, resp, err = client.Check.Get("aaaaaa")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}

func TestListDowntimes(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	// New check won't have downtimes, but API should respond OK
	downs, resp, err := client.Downtime.List(token, 1)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 0, len(downs)) // New check has no downtimes
}

func TestAddUpdateRemoveCheck(t *testing.T) {
	client := newClient()

	// Add
	res, resp, err := client.Check.Add(CheckItem{URL: "https://google.fr"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "https://google.fr", res.URL)

	// Update
	res, resp, err = client.Check.Update(res.Token, CheckItem{URL: "https://google.com"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "https://google.com", res.URL)

	// Remove
	result, resp, err := client.Check.Remove(res.Token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestAddCheckWithType(t *testing.T) {
	client := newClient()

	// Test ICMP check
	res, resp, err := client.Check.Add(CheckItem{
		URL:  "8.8.8.8",
		Type: "icmp",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "icmp", res.Type)

	// Clean up
	_, _, _ = client.Check.Remove(res.Token)
}

func TestListMetrics(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	// Wait a moment for the check to be processed
	time.Sleep(2 * time.Second)

	now := time.Now()
	timeFormat := "2006-01-02 15:04:05 -0700"
	from, to := now.AddDate(0, 0, -1).Format(timeFormat), now.Format(timeFormat)
	metricRes, resp, err := client.Metric.List(token, "host", from, to)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// New check may not have metrics yet, just verify API works
	_ = metricRes
}

func TestListNodes(t *testing.T) {
	client := newClient()
	nodeRes, resp, err := client.Node.List()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, len(nodeRes) > 0, "Should have at least one node")
}

func TestListIPv4(t *testing.T) {
	client := newClient()
	IPs, resp, err := client.Node.ListIPv4()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, len(IPs) > 0, "Should have at least one IPv4 address")

	for _, ip := range IPs {
		parsed := net.ParseIP(ip)
		assert.NotNil(t, parsed, "Should be valid IP: %s", ip)
		assert.True(t, isIPv4(parsed), "Should be IPv4: %s", ip)
	}
}

func TestListIPv6(t *testing.T) {
	client := newClient()
	IPs, resp, err := client.Node.ListIPv6()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, len(IPs) > 0, "Should have at least one IPv6 address")

	for _, ip := range IPs {
		parsed := net.ParseIP(ip)
		assert.NotNil(t, parsed, "Should be valid IP: %s", ip)
		assert.True(t, isIPv6(parsed), "Should be IPv6: %s", ip)
	}
}

func isIPv4(ip net.IP) bool {
	return ip.To4() != nil
}

func isIPv6(ip net.IP) bool {
	return ip.To4() == nil && ip.To16() != nil
}
