package helpers

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new TestHelper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateFiberApp creates a new Fiber app for testing
func (h *TestHelper) CreateFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

// CreateJSONRequest creates an HTTP request with JSON body
func (h *TestHelper) CreateJSONRequest(method, url string, body interface{}) *httptest.Request {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		assert.NoError(h.t, err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// CreateAuthenticatedRequest creates an HTTP request with JWT token
func (h *TestHelper) CreateAuthenticatedRequest(method, url string, body interface{}, token string) *httptest.Request {
	req := h.CreateJSONRequest(method, url, body)
	req.Header.Set("Authorization", "Bearer "+token)
	return req
}

// SetFiberLocals sets locals in Fiber context for testing
func (h *TestHelper) SetFiberLocals(c *fiber.Ctx, userID uuid.UUID, username, email, role string) {
	c.Locals("user_id", userID)
	c.Locals("username", username)
	c.Locals("email", email)
	c.Locals("role", role)
}

// CreateMiddleware creates a middleware that sets user context
func (h *TestHelper) CreateMiddleware(userID uuid.UUID, username, email, role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h.SetFiberLocals(c, userID, username, email, role)
		return c.Next()
	}
}

// AssertJSONResponse asserts that the response contains expected JSON
func (h *TestHelper) AssertJSONResponse(resp *httptest.ResponseRecorder, expectedStatus int, expectedBody map[string]interface{}) {
	assert.Equal(h.t, expectedStatus, resp.Code)
	
	var actualBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &actualBody)
	assert.NoError(h.t, err)
	
	for key, expectedValue := range expectedBody {
		assert.Equal(h.t, expectedValue, actualBody[key])
	}
}

// AssertErrorResponse asserts that the response contains an error
func (h *TestHelper) AssertErrorResponse(resp *httptest.ResponseRecorder, expectedStatus int, expectedError string) {
	assert.Equal(h.t, expectedStatus, resp.Code)
	
	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(h.t, err)
	
	assert.Contains(h.t, body, "error")
	if expectedError != "" {
		assert.Equal(h.t, expectedError, body["error"])
	}
}

// AssertSuccessResponse asserts that the response is successful
func (h *TestHelper) AssertSuccessResponse(resp *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(h.t, expectedStatus, resp.Code)
	
	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(h.t, err)
	
	// Check for success indicators
	if message, exists := body["message"]; exists {
		assert.NotEmpty(h.t, message)
	}
}

// GenerateUUID generates a new UUID for testing
func (h *TestHelper) GenerateUUID() uuid.UUID {
	return uuid.New()
}

// GenerateUUIDs generates multiple UUIDs for testing
func (h *TestHelper) GenerateUUIDs(count int) []uuid.UUID {
	uuids := make([]uuid.UUID, count)
	for i := 0; i < count; i++ {
		uuids[i] = uuid.New()
	}
	return uuids
}

// MockJWTToken creates a mock JWT token for testing
func (h *TestHelper) MockJWTToken() string {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
}

// CompareJSON compares two JSON objects for equality
func (h *TestHelper) CompareJSON(expected, actual interface{}) {
	expectedJSON, err := json.Marshal(expected)
	assert.NoError(h.t, err)
	
	actualJSON, err := json.Marshal(actual)
	assert.NoError(h.t, err)
	
	assert.JSONEq(h.t, string(expectedJSON), string(actualJSON))
}

// AssertPaginationResponse asserts pagination structure in response
func (h *TestHelper) AssertPaginationResponse(resp *httptest.ResponseRecorder, expectedPage, expectedPerPage int) {
	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(h.t, err)
	
	data, exists := body["data"].(map[string]interface{})
	assert.True(h.t, exists, "Response should contain data field")
	
	pagination, exists := data["pagination"].(map[string]interface{})
	assert.True(h.t, exists, "Data should contain pagination field")
	
	assert.Equal(h.t, float64(expectedPage), pagination["current_page"])
	assert.Equal(h.t, float64(expectedPerPage), pagination["per_page"])
	assert.Contains(h.t, pagination, "total_pages")
	assert.Contains(h.t, pagination, "total_items")
}

// CreateFormRequest creates a multipart form request for file uploads
func (h *TestHelper) CreateFormRequest(method, url string, fields map[string]string, files map[string][]byte) *httptest.Request {
	// This would be implemented for file upload testing
	// For now, return a basic request
	req := httptest.NewRequest(method, url, nil)
	req.Header.Set("Content-Type", "multipart/form-data")
	return req
}

// AssertValidationError asserts that the response contains validation errors
func (h *TestHelper) AssertValidationError(resp *httptest.ResponseRecorder, expectedFields []string) {
	assert.Equal(h.t, fiber.StatusBadRequest, resp.Code)
	
	var body map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &body)
	assert.NoError(h.t, err)
	
	assert.Contains(h.t, body, "error")
	
	// If there are validation details, check them
	if details, exists := body["details"].([]interface{}); exists {
		fieldMap := make(map[string]bool)
		for _, detail := range details {
			if detailMap, ok := detail.(map[string]interface{}); ok {
				if field, ok := detailMap["field"].(string); ok {
					fieldMap[field] = true
				}
			}
		}
		
		for _, expectedField := range expectedFields {
			assert.True(h.t, fieldMap[expectedField], "Expected validation error for field: %s", expectedField)
		}
	}
}

// MockTime provides a consistent time for testing
func (h *TestHelper) MockTime() string {
	return "2024-01-15T10:30:00Z"
}

// AssertContainsKeys asserts that a map contains specific keys
func (h *TestHelper) AssertContainsKeys(data map[string]interface{}, keys []string) {
	for _, key := range keys {
		assert.Contains(h.t, data, key, "Response should contain key: %s", key)
	}
}

// AssertNotContainsKeys asserts that a map does not contain specific keys
func (h *TestHelper) AssertNotContainsKeys(data map[string]interface{}, keys []string) {
	for _, key := range keys {
		assert.NotContains(h.t, data, key, "Response should not contain key: %s", key)
	}
}