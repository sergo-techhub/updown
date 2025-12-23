package updown

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// Test URLs/hosts
	testHTTPURL    = "https://example.com"
	testHTTPURLAlt = "https://google.fr"
	testHTTPURLUpd = "https://google.com"
	testICMPHost   = "8.8.8.8"
	testTCPHost    = "tcp://google.com:443"
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
		URL:   testHTTPURL,
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

	// Use unique alias to avoid conflicts with pre-existing checks
	uniqueAlias := fmt.Sprintf("Test Check %d", time.Now().UnixNano())

	// Create a test check with unique alias
	res, resp, err := client.Check.Add(CheckItem{
		URL:   testHTTPURL,
		Alias: uniqueAlias,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	token := res.Token
	defer deleteTestCheck(t, client, token)

	// Verify check was created by getting it directly
	check, _, err := client.Check.Get(token)
	require.NoError(t, err)
	require.Equal(t, uniqueAlias, check.Alias)

	// Cache miss + alias not found
	foundToken, err := client.Check.TokenForAlias("nonexistent-alias-12345")
	assert.Equal(t, "", foundToken)
	assert.Equal(t, ErrTokenNotFound, err)

	// Cache miss + match found after request
	foundToken, err = client.Check.TokenForAlias(uniqueAlias)
	assert.Nil(t, err)
	assert.Equal(t, token, foundToken)

	// Cache hit
	foundToken, err = client.Check.TokenForAlias(uniqueAlias)
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
	res, resp, err := client.Check.Add(CheckItem{URL: testHTTPURLAlt})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, testHTTPURLAlt, res.URL)

	// Update
	res, resp, err = client.Check.Update(res.Token, CheckItem{URL: testHTTPURLUpd})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, testHTTPURLUpd, res.URL)

	// Remove
	result, resp, err := client.Check.Remove(res.Token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestAddICMPCheck(t *testing.T) {
	client := newClient()

	// Test ICMP check
	res, resp, err := client.Check.Add(CheckItem{
		URL:  testICMPHost,
		Type: "icmp",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "icmp", res.Type)

	// Clean up
	_, _, _ = client.Check.Remove(res.Token)
}

func TestAddTCPCheck(t *testing.T) {
	client := newClient()

	// Test TCP check
	res, resp, err := client.Check.Add(CheckItem{
		URL:  testTCPHost,
		Type: "tcp",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "tcp", res.Type)

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
	from, to := now.AddDate(0, 0, -1).Format("2006-01-02"), now.Format("2006-01-02")
	metricRes, resp, err := client.Metric.List(token, "time", from, to)

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

func TestListRecipients(t *testing.T) {
	client := newClient()
	recipients, resp, err := client.Recipient.List()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// May have zero or more recipients
	assert.NotNil(t, recipients)
}

func TestAddRemoveRecipient(t *testing.T) {
	client := newClient()

	// Add email recipient
	res, resp, err := client.Recipient.Add(RecipientItem{
		Type:  RecipientTypeEmail,
		Value: "test@example.com",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, RecipientTypeEmail, res.Type)
	assert.NotEmpty(t, res.ID)

	// Remove recipient
	result, resp, err := client.Recipient.Remove(res.ID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestAddWebhookRecipient(t *testing.T) {
	client := newClient()

	// Add webhook recipient
	res, resp, err := client.Recipient.Add(RecipientItem{
		Type:  RecipientTypeWebhook,
		Value: "https://example.com/webhook",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, RecipientTypeWebhook, res.Type)

	// Clean up
	_, _, _ = client.Recipient.Remove(res.ID)
}

func TestListStatusPages(t *testing.T) {
	client := newClient()
	pages, resp, err := client.StatusPage.List()

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	// May have zero or more status pages
	assert.NotNil(t, pages)
}

func TestAddUpdateRemoveStatusPage(t *testing.T) {
	client := newClient()

	// Create a test check for the status page
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	// Add status page
	res, resp, err := client.StatusPage.Add(StatusPageItem{
		Name:       "Test Status Page",
		Visibility: "private",
		Checks:     []string{token},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "Test Status Page", res.Name)
	assert.Equal(t, "private", res.Visibility)
	assert.NotEmpty(t, res.Token)

	// Get status page
	page, resp, err := client.StatusPage.Get(res.Token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, res.Token, page.Token)
	assert.Equal(t, "Test Status Page", page.Name)

	// Update status page
	updated, resp, err := client.StatusPage.Update(res.Token, StatusPageItem{
		Name:       "Updated Status Page",
		Visibility: "private",
		Checks:     []string{token},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Updated Status Page", updated.Name)

	// Remove status page
	result, resp, err := client.StatusPage.Remove(res.Token)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, result)
}

func TestStatusPageProtected(t *testing.T) {
	client := newClient()

	// Create a test check
	token := createTestCheck(t, client)
	defer deleteTestCheck(t, client, token)

	// Add protected status page with custom access key
	res, resp, err := client.StatusPage.Add(StatusPageItem{
		Name:       "Protected Page",
		Visibility: "protected",
		AccessKey:  "test-access-key-123",
		Checks:     []string{token},
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "protected", res.Visibility)
	assert.Equal(t, "test-access-key-123", res.AccessKey)

	// Clean up
	_, _, _ = client.StatusPage.Remove(res.Token)
}
