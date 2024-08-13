package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Excoriate/daggerverse/module-template/internal/dagger"
)

// headerSplitCount is the number of parts expected when splitting a header string.
const headerSplitCount = 2

// WithHTTPCurl sets up the HTTP client and configuration for the ModuleTemplate
// You can configure various options like headers, authentication, etc.
//
// Example:
//
//	module := &ModuleTemplate{}
//	module = module.WithHTTP("https://api.example.com",
//	                         []string{"Content-Type=application/json"},
//	                         time.Second*30,
//	                         "basic", "username:password")
//
// Parameters:
//   - baseURL: The base URL for the HTTP requests.
//   - headers: A slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
//   - timeout: The timeout duration for the HTTP client.
//   - authType: The type of authentication (e.g., "basic", "bearer").
//   - authCredentials: The credentials for the specified authentication type.
//
// Returns:
//   - A pointer to the updated ModuleTemplate.
func (m *ModuleTemplate) WithHTTPCurl(
	// baseURL is the base URL for the HTTP requests.
	baseURL string,
	// headers is a slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
	// +optional
	headers []string,
	// timeout is the timeout duration for the HTTP client.
	// +optional
	timeout string,
	// authType is the type of authentication (e.g., "basic", "bearer").
	// +optional
	authType string,
	// authCredentials is the credentials for the specified authentication type.
	// +optional
	authCredentials string,
) *ModuleTemplate {
	headerMap := parseHeaders(headers)
	timeoutDuration, _ := parseTimeout(timeout)

	curlCmd := buildCurlCommand(baseURL, headerMap, timeoutDuration, authType, authCredentials)

	m.Ctr = m.Ctr.WithExec([]string{"sh", "-c", curlCmd})

	return m
}

// DoHTTPRequest sets up the HTTP client and configuration for the ModuleTemplate
// You can configure various options like headers, authentication, etc.
//
// Example:
//
//	module := &ModuleTemplate{}
//	module = module.DoHTTPRequest(ctx, "GET", "https://api.example.com",
//	                              []string{"Content-Type=application/json"},
//	                              time.Second*30,
//	                              "basic", "username:password", "")
//
// Parameters:
//   - ctx: The context for the HTTP request.
//   - method: The HTTP method (e.g., "GET", "POST").
//   - baseURL: The base URL for the HTTP requests.
//   - headers: A slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
//   - timeout: The timeout duration for the HTTP client.
//   - authType: The type of authentication (e.g., "basic", "bearer").
//   - authCredentials: The credentials for the specified authentication type.
//   - body: The request body for POST requests (can be empty string for GET requests).
//
// Returns:
//   - A pointer to the updated ModuleTemplate.
func (m *ModuleTemplate) DoHTTPRequest(
	ctx context.Context,
	// method is the HTTP method (e.g., "GET", "POST").
	method string,
	// baseURL is the base URL for the HTTP requests.
	baseURL string,
	// headers is a slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
	// +optional
	headers []string,
	// timeout is the timeout duration for the HTTP client.
	// +optional
	timeout string,
	// authType is the type of authentication (e.g., "basic", "bearer").
	// +optional
	authType string,
	// authCredentials is the credentials for the specified authentication type.
	// +optional
	authCredentials string,
	// body is the request body for POST requests (can be empty string for GET requests).
	// +optional
	body string,
) (*dagger.Container, error) {
	headerMap := parseHeaders(headers)

	timeoutDuration, _ := parseTimeout(timeout)

	client := &http.Client{
		Timeout: timeoutDuration,
	}

	req, err := createRequest(ctx, method, baseURL, body)
	if err != nil {
		return nil, WrapError(err, "failed to create HTTP request")
	}

	addHeaders(req, headerMap)

	if err := addAuthentication(req, authType, authCredentials); err != nil {
		return nil, err
	}

	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return nil, WrapError(clientErr, "Error executing HTTP request")
	}
	defer safeCloseBody(resp)

	respBody, bodyErr := io.ReadAll(resp.Body)
	if bodyErr != nil {
		return nil, WrapError(bodyErr, "Error reading HTTP response body")
	}

	m.Ctr = m.Ctr.WithNewFile("/http_response.txt", string(respBody))
	m.Ctr = m.Ctr.WithNewFile("/http_status.txt", fmt.Sprintf("%d", resp.StatusCode))

	return m.Ctr, nil
}

// DoJSONAPICall performs an API call and returns the JSON response as a dagger.File.
//
// Example:
//
//	module := &ModuleTemplate{}
//	jsonFile, err := module.DoJSONAPICall(ctx, "POST", "https://api.example.com/data",
//	                                      []string{"Content-Type=application/json"},
//	                                      time.Second*30,
//	                                      "bearer", "your-token-here",
//	                                      `{"key": "value"}`)
//
// Parameters:
//   - ctx: The context for the HTTP request.
//   - method: The HTTP method (e.g., "GET", "POST").
//   - baseURL: The base URL for the HTTP requests.
//   - headers: A slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
//   - timeout: The timeout duration for the HTTP client.
//   - authType: The type of authentication (e.g., "basic", "bearer").
//   - authCredentials: The credentials for the specified authentication type.
//   - body: The JSON request body for POST requests (can be empty string for GET requests).
//
// Returns:
//   - A dagger.File containing the JSON response.
//   - An error if the request fails or the response is not JSON.
func (m *ModuleTemplate) DoJSONAPICall(
	ctx context.Context,
	// method is the HTTP method (e.g., "GET", "POST").
	method string,
	// baseURL is the base URL for the HTTP requests.
	baseURL string,
	// headers is a slice of headers to include in the HTTP requests, formatted as "Header=Value" strings.
	// +optional
	headers []string,
	// timeout is the timeout duration for the HTTP client.
	// +optional
	timeout string,
	// authType is the type of authentication (e.g., "basic", "bearer").
	// +optional
	authType string,
	// authCredentials is the credentials for the specified authentication type.
	// +optional
	authCredentials string,
	// body is the JSON request body for POST requests (can be empty string for GET requests).
	// +optional
	body string,
) (*dagger.File, error) {
	if authType == "" {
		authType = "none"
	}

	container, err := m.DoHTTPRequest(ctx, method, baseURL, headers, timeout, authType, authCredentials, body)
	if err != nil {
		return nil, WrapError(err, "Failed to perform HTTP request")
	}

	// Get the HTTP status code
	statusCode, catErr := container.WithExec([]string{"cat", "/http_status.txt"}).Stdout(ctx)
	if catErr != nil {
		return nil, WrapError(catErr, "Failed to read HTTP status code")
	}

	content, fileErr := container.File("/http_response.txt").Contents(ctx)
	if fileErr != nil {
		return nil, WrapError(fileErr, "Failed to read HTTP response")
	}

	// If status code is not 2xx, return an error with more details
	statusCodeInt, _ := strconv.Atoi(strings.TrimSpace(statusCode))
	if statusCodeInt < 200 || statusCodeInt >= 300 {
		return nil, Errorf("HTTP request failed with status %s. Response: %s", statusCode, content)
	}

	jsonFile := dag.Directory().WithNewFile("response.json", content)

	return jsonFile.File("response.json"), nil
}

// Helper functions
func parseHeaders(headers []string) map[string]string {
	headerMap := make(map[string]string)

	for _, header := range headers {
		parts := strings.SplitN(header, "=", headerSplitCount)
		if len(parts) == headerSplitCount {
			headerMap[parts[0]] = parts[1]
		} else {
			fmt.Printf("Malformed header: %s\n", header)
		}
	}

	return headerMap
}

func buildCurlCommand(
	baseURL string,
	headers map[string]string,
	timeout time.Duration,
	authType,
	authCredentials string) string {
	curlCmd := fmt.Sprintf("curl -m %d", int(timeout.Seconds()))
	for key, value := range headers {
		curlCmd += fmt.Sprintf(" -H '%s: %s'", key, value)
	}

	switch authType {
	case "basic":
		curlCmd += fmt.Sprintf(" -u '%s'", authCredentials)
	case "bearer":
		curlCmd += fmt.Sprintf(" -H 'Authorization: Bearer %s'", authCredentials)
	}

	curlCmd += fmt.Sprintf(" '%s'", baseURL)

	return curlCmd
}

func createRequest(ctx context.Context, method, baseURL, body string) (*http.Request, error) {
	var (
		req *http.Request
		err error
	)

	if method == "" {
		method = http.MethodGet
	}

	if method == http.MethodGet {
		method = http.MethodPost
	}

	if method == http.MethodPost && body != "" {
		req, err = http.NewRequestWithContext(ctx, method, baseURL, strings.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, baseURL, http.NoBody)
	}

	if err != nil {
		return nil, WrapError(err, "failed to create HTTP request")
	}

	return req, nil
}

func addHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Add(key, value)
	}
}

func addAuthentication(req *http.Request, authType, authCredentials string) error {
	switch authType {
	case "basic":
		creds := strings.SplitN(authCredentials, ":", headerSplitCount)
		if len(creds) != headerSplitCount {
			return WrapError(nil, "Invalid credentials format for basic authentication. Expected format: username:password")
		}

		req.SetBasicAuth(creds[0], creds[1])
	case "bearer":
		req.Header.Add("Authorization", "Bearer "+authCredentials)
	case "none":
		// Do nothing for no authentication
	default:
		return WrapErrorf(nil, "Unsupported authentication type: %s", authType)
	}

	return nil
}
func safeCloseBody(resp *http.Response) {
	if resp != nil {
		if closeErr := resp.Body.Close(); closeErr != nil {
			fmt.Printf("Error closing response body: %v\n", closeErr)
		}
	}
}

func parseTimeout(timeout string) (time.Duration, error) {
	if timeout == "" {
		return 0, nil
	}

	parsedTimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return 0, WrapError(err, "Failed to parse timeout")
	}
	return parsedTimeout, nil
}
