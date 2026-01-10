package http

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	adapterhttp "github.com/zevwings/workflow/internal/infrastructure/http"
	"github.com/zevwings/workflow/internal/logging"
)

var (
	// globalClient global HTTP client singleton
	globalClient *httpClient
	globalOnce   sync.Once
)

// Client HTTP client interface
//
// Provides unified HTTP request interface, encapsulates underlying resty client.
// Supports GET, POST, PUT, DELETE, PATCH and other HTTP methods, as well as streaming requests and multipart requests.
type Client interface {
	// SetAuth sets authentication token
	SetAuth(token string)
	// SetBasicAuth sets Basic Auth
	SetBasicAuth(username, password string)
	// SetProxy sets proxy
	SetProxy(proxyURL string)
	// Get sends GET request (legacy API, maintained for backward compatibility)
	Get(url string) (*resty.Response, error)
	// Post sends POST request (legacy API, maintained for backward compatibility)
	Post(url string, body interface{}) (*resty.Response, error)
	// Put sends PUT request (legacy API, maintained for backward compatibility)
	Put(url string, body interface{}) (*resty.Response, error)
	// Delete sends DELETE request (legacy API, maintained for backward compatibility)
	Delete(url string) (*resty.Response, error)
	// Patch sends PATCH request (legacy API, maintained for backward compatibility)
	Patch(url string, body interface{}) (*resty.Response, error)
	// GetWithConfig sends GET request (new API, supports RequestConfig)
	GetWithConfig(url string, config *RequestConfig) (*HttpResponse, error)
	// PostWithConfig sends POST request (new API, supports RequestConfig)
	PostWithConfig(url string, config *RequestConfig) (*HttpResponse, error)
	// PutWithConfig sends PUT request (new API, supports RequestConfig)
	PutWithConfig(url string, config *RequestConfig) (*HttpResponse, error)
	// DeleteWithConfig sends DELETE request (new API, supports RequestConfig)
	DeleteWithConfig(url string, config *RequestConfig) (*HttpResponse, error)
	// PatchWithConfig sends PATCH request (new API, supports RequestConfig)
	PatchWithConfig(url string, config *RequestConfig) (*HttpResponse, error)
	// Stream streaming request
	Stream(method HttpMethod, url string, config *RequestConfig) (io.ReadCloser, error)
	// PostMultipart POST Multipart request
	PostMultipart(url string, config *MultipartRequestConfig) (*HttpResponse, error)
	// GetRestyClient gets underlying resty client (for advanced usage)
	GetRestyClient() *resty.Client
}

// httpClient HTTP client implementation
type httpClient struct {
	client *resty.Client
}

// DefaultRetryCondition default retry condition function (exported for use by other packages)
func DefaultRetryCondition(r *resty.Response, err error) bool {
	// Network errors or connection errors should be retried
	if err != nil {
		return isRetryableNetworkError(err)
	}

	// HTTP status code judgment
	statusCode := r.StatusCode()
	// 5xx server errors and 429 Too Many Requests are retryable
	if statusCode >= 500 && statusCode < 600 {
		return true
	}
	if statusCode == 429 {
		return true
	}

	// 4xx client errors are not retryable
	return false
}

// DefaultRetryAfter default retry delay function (exported for use by other packages)
func DefaultRetryAfter(client *resty.Client, resp *resty.Response) (time.Duration, error) {
	// If response contains Retry-After header, use it
	if retryAfter := resp.Header().Get("Retry-After"); retryAfter != "" {
		if duration, err := parseRetryAfter(retryAfter); err == nil {
			return duration, nil
		}
	}

	// Otherwise use exponential backoff: calculate delay based on retry count
	// resty will automatically handle exponential backoff, return 0 here to let resty use default exponential backoff
	return 0, nil
}

// newClient creates a new HTTP client (internal function, not exported)
//
// External code should use Global() to get singleton client, rather than directly creating new instances.
func newClient() *httpClient {
	client := resty.New()
	client.SetTimeout(30 * time.Second)

	// Set Logrus Logger (implemented through adapter)
	client.SetLogger(adapterhttp.NewLogrusLogger())

	// Configure default retry strategy
	client.SetRetryCount(3)
	client.SetRetryWaitTime(1 * time.Second)
	client.SetRetryMaxWaitTime(30 * time.Second)
	client.AddRetryCondition(DefaultRetryCondition)
	client.SetRetryAfter(DefaultRetryAfter)

	// Add logging hooks
	setupLoggingHooks(client)

	return &httpClient{client: client}
}

// setupLoggingHooks sets up HTTP request logging hooks
//
// Uses Resty's Hook mechanism to record various stages of HTTP requests:
//   - OnBeforeRequest: Record before request is sent (Info level)
//   - OnAfterResponse: Record request success/failure (Info/Warn level)
//   - OnError: Record request errors (Error level)
//
// All log output automatically filters sensitive information before output (API keys in URLs, sensitive information in request headers, etc.).
func setupLoggingHooks(client *resty.Client) {
	// Before request hook
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Filter sensitive information in URL
		filteredURL := FilterSensitiveURL(req.URL)
		logging.GetLogger().Infof("HTTP %s request to %s", req.Method, filteredURL)
		return nil
	})

	// After response hook
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		method := resp.Request.Method
		// Filter sensitive information in URL
		filteredURL := FilterSensitiveURL(resp.Request.URL)
		statusCode := resp.StatusCode()

		if resp.IsSuccess() {
			// 2xx status code: success
			logging.GetLogger().Infof("HTTP %s request to %s succeeded: %d", method, filteredURL, statusCode)
		} else if statusCode >= 400 && statusCode < 500 {
			// 4xx status code: client error
			logging.GetLogger().Warnf("HTTP %s request to %s returned client error: %d", method, filteredURL, statusCode)
		} else if statusCode >= 500 {
			// 5xx status code: server error
			logging.GetLogger().Warnf("HTTP %s request to %s returned server error: %d", method, filteredURL, statusCode)
		}

		return nil
	})

	// Error hook
	client.OnError(func(req *resty.Request, err error) {
		if err != nil {
			// Filter sensitive information in URL
			filteredURL := FilterSensitiveURL(req.URL)
			logging.GetLogger().Errorf("HTTP %s request to %s failed: %v", req.Method, filteredURL, err)
		}
	})
}

// isRetryableNetworkError determines if network error is retryable
func isRetryableNetworkError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()

	// Network connection errors are retryable
	retryableKeywords := []string{
		"timeout",
		"connection",
		"network",
		"dial",
		"connection refused",
		"connection reset",
		"no such host",
		"temporary failure",
	}

	for _, keyword := range retryableKeywords {
		if contains(errStr, keyword) {
			return true
		}
	}

	return false
}

// parseRetryAfter parses Retry-After header value
func parseRetryAfter(value string) (time.Duration, error) {
	// Retry-After can be seconds (number) or HTTP date
	// Here simplified processing, only supports seconds
	if seconds, err := parseInt(value); err == nil {
		return time.Duration(seconds) * time.Second, nil
	}

	// If unable to parse, return error
	return 0, fmt.Errorf("invalid Retry-After value: %s", value)
}

// contains checks if string contains substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// parseInt attempts to parse string as integer
func parseInt(s string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(s))
}

// Global gets global HttpClient singleton
//
// Returns process-level HttpClient singleton with default configuration.
// Singleton is initialized on first call, subsequent calls reuse the same instance.
//
// This is the recommended way to get HTTP client, external code should use this method rather than directly creating new instances.
//
// Advantages:
//   - Reuse connection pool: All requests share the same connection pool, improving performance
//   - Reduce resource consumption: Avoid duplicate client instance creation
//   - Thread-safe: Can be safely used in multi-threaded environments
//   - Encapsulation: Returns interface type, hides implementation details, prevents external direct access to internal structure
//
// Returns:
//   - Client: Global HTTP client interface instance
func Global() Client {
	globalOnce.Do(func() {
		globalClient = newClient()
	})
	return globalClient
}

// SetAuth sets authentication token
func (c *httpClient) SetAuth(token string) {
	c.client.SetAuthToken(token)
}

// SetBasicAuth sets Basic Auth
func (c *httpClient) SetBasicAuth(username, password string) {
	c.client.SetBasicAuth(username, password)
}

// SetProxy sets proxy
func (c *httpClient) SetProxy(proxyURL string) {
	c.client.SetProxy(proxyURL)
}

// Get sends GET request (legacy API, maintained for backward compatibility)
func (c *httpClient) Get(url string) (*resty.Response, error) {
	return c.client.R().Get(url)
}

// Post sends POST request (legacy API, maintained for backward compatibility)
func (c *httpClient) Post(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Post(url)
}

// Put sends PUT request (legacy API, maintained for backward compatibility)
func (c *httpClient) Put(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Put(url)
}

// Delete sends DELETE request (legacy API, maintained for backward compatibility)
func (c *httpClient) Delete(url string) (*resty.Response, error) {
	return c.client.R().Delete(url)
}

// Patch sends PATCH request (legacy API, maintained for backward compatibility)
func (c *httpClient) Patch(url string, body interface{}) (*resty.Response, error) {
	return c.client.R().SetBody(body).Patch(url)
}

// doRequest common method to execute HTTP request
//
// Parameters:
//   - method: HTTP method
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) doRequest(method HttpMethod, url string, config *RequestConfig) (*HttpResponse, error) {
	if config == nil {
		config = NewRequestConfig()
	}

	// If custom retry configuration is provided, create temporary client
	var client *resty.Client
	if config.Retry != nil {
		client = applyRetryConfig(c.client, config.Retry)
	} else {
		client = c.client
	}

	req := client.R()
	req = config.applyToRequest(req)

	var resp *resty.Response
	var err error

	switch method {
	case MethodGet:
		resp, err = req.Get(url)
	case MethodPost:
		resp, err = req.Post(url)
	case MethodPut:
		resp, err = req.Put(url)
	case MethodDelete:
		resp, err = req.Delete(url)
	case MethodPatch:
		resp, err = req.Patch(url)
	default:
		return nil, &InvalidMethodError{Method: string(method)}
	}

	if err != nil {
		return nil, err
	}

	return FromRestyResponse(resp)
}

// GetWithConfig sends GET request (new API, supports RequestConfig)
//
// Parameters:
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) GetWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodGet, url, config)
}

// PostWithConfig sends POST request (new API, supports RequestConfig)
//
// Parameters:
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) PostWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPost, url, config)
}

// PutWithConfig sends PUT request (new API, supports RequestConfig)
//
// Parameters:
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) PutWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPut, url, config)
}

// DeleteWithConfig sends DELETE request (new API, supports RequestConfig)
//
// Parameters:
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) DeleteWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodDelete, url, config)
}

// PatchWithConfig sends PATCH request (new API, supports RequestConfig)
//
// Parameters:
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) PatchWithConfig(url string, config *RequestConfig) (*HttpResponse, error) {
	return c.doRequest(MethodPatch, url, config)
}

// Stream streaming request
//
// Sends request and returns response stream, used for handling large files or streaming data.
//
// Parameters:
//   - method: HTTP method
//   - url: Request URL
//   - config: Request configuration (optional, if nil uses default configuration)
//
// Returns:
//   - io.ReadCloser: Response stream
//   - error: Returns error if request fails
func (c *httpClient) Stream(method HttpMethod, url string, config *RequestConfig) (io.ReadCloser, error) {
	if config == nil {
		config = NewRequestConfig()
	}

	req := c.client.R()
	req = config.applyToRequest(req)
	// Set not to automatically parse response to support streaming read
	req.SetDoNotParseResponse(true)

	var resp *resty.Response
	var err error

	switch method {
	case MethodGet:
		resp, err = req.Get(url)
	case MethodPost:
		resp, err = req.Post(url)
	case MethodPut:
		resp, err = req.Put(url)
	case MethodDelete:
		resp, err = req.Delete(url)
	case MethodPatch:
		resp, err = req.Patch(url)
	default:
		return nil, &InvalidMethodError{Method: string(method)}
	}

	if err != nil {
		return nil, err
	}

	return resp.RawBody(), nil
}

// PostMultipart POST Multipart request
//
// Sends multipart/form-data request, typically used for file uploads.
//
// Parameters:
//   - url: Request URL
//   - config: Multipart request configuration
//
// Returns:
//   - *HttpResponse: Encapsulated HTTP response
//   - error: Returns error if request fails
func (c *httpClient) PostMultipart(url string, config *MultipartRequestConfig) (*HttpResponse, error) {
	if config == nil {
		return nil, &ConfigError{Message: "MultipartRequestConfig is required for multipart requests"}
	}

	// If custom retry configuration is provided, create temporary client
	var client *resty.Client
	if config.Retry != nil {
		client = applyRetryConfig(c.client, config.Retry)
	} else {
		client = c.client
	}

	req := client.R()
	req = config.applyToRequest(req)

	resp, err := req.Post(url)
	if err != nil {
		return nil, err
	}

	return FromRestyResponse(resp)
}

// GetRestyClient gets underlying resty client (for advanced usage)
func (c *httpClient) GetRestyClient() *resty.Client {
	return c.client
}

// InvalidMethodError invalid HTTP method error
type InvalidMethodError struct {
	Method string
}

func (e *InvalidMethodError) Error() string {
	return "invalid HTTP method: " + e.Method
}

// ConfigError configuration error
type ConfigError struct {
	Message string
}

func (e *ConfigError) Error() string {
	return "config error: " + e.Message
}
