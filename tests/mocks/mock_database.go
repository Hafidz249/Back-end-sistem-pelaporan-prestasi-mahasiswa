package mocks

import (
	"database/sql"
	"database/sql/driver"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockDatabase provides mock database connections
type MockDatabase struct {
	mock.Mock
	PostgresDB *sql.DB
	PostgresMock sqlmock.Sqlmock
	MongoDB    *MockMongoDB
}

// NewMockDatabase creates a new mock database instance
func NewMockDatabase() (*MockDatabase, error) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		return nil, err
	}

	return &MockDatabase{
		PostgresDB:   db,
		PostgresMock: mock,
		MongoDB:      NewMockMongoDB(),
	}, nil
}

// Close closes the mock database connections
func (m *MockDatabase) Close() error {
	return m.PostgresDB.Close()
}

// ExpectationsWereMet checks if all expectations were met
func (m *MockDatabase) ExpectationsWereMet() error {
	return m.PostgresMock.ExpectationsWereMet()
}

// MockMongoDB provides mock MongoDB operations
type MockMongoDB struct {
	mock.Mock
}

// NewMockMongoDB creates a new mock MongoDB instance
func NewMockMongoDB() *MockMongoDB {
	return &MockMongoDB{}
}

// Collection mocks MongoDB collection
func (m *MockMongoDB) Collection(name string) *MockMongoCollection {
	args := m.Called(name)
	return args.Get(0).(*MockMongoCollection)
}

// MockMongoCollection provides mock MongoDB collection operations
type MockMongoCollection struct {
	mock.Mock
}

// InsertOne mocks MongoDB InsertOne operation
func (m *MockMongoCollection) InsertOne(ctx interface{}, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, document)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// FindOne mocks MongoDB FindOne operation
func (m *MockMongoCollection) FindOne(ctx interface{}, filter interface{}) *MockSingleResult {
	args := m.Called(ctx, filter)
	return args.Get(0).(*MockSingleResult)
}

// Find mocks MongoDB Find operation
func (m *MockMongoCollection) Find(ctx interface{}, filter interface{}) (*MockCursor, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockCursor), args.Error(1)
}

// UpdateOne mocks MongoDB UpdateOne operation
func (m *MockMongoCollection) UpdateOne(ctx interface{}, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

// DeleteOne mocks MongoDB DeleteOne operation
func (m *MockMongoCollection) DeleteOne(ctx interface{}, filter interface{}) (*mongo.DeleteResult, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*mongo.DeleteResult), args.Error(1)
}

// MockSingleResult mocks MongoDB SingleResult
type MockSingleResult struct {
	mock.Mock
}

// Decode mocks SingleResult Decode operation
func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// MockCursor mocks MongoDB Cursor
type MockCursor struct {
	mock.Mock
}

// Next mocks Cursor Next operation
func (m *MockCursor) Next(ctx interface{}) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

// Decode mocks Cursor Decode operation
func (m *MockCursor) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

// All mocks Cursor All operation
func (m *MockCursor) All(ctx interface{}, results interface{}) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// Close mocks Cursor Close operation
func (m *MockCursor) Close(ctx interface{}) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockFileSystem provides mock file system operations
type MockFileSystem struct {
	mock.Mock
}

// WriteFile mocks file writing
func (m *MockFileSystem) WriteFile(filename string, data []byte, perm int) error {
	args := m.Called(filename, data, perm)
	return args.Error(0)
}

// ReadFile mocks file reading
func (m *MockFileSystem) ReadFile(filename string) ([]byte, error) {
	args := m.Called(filename)
	return args.Get(0).([]byte), args.Error(1)
}

// DeleteFile mocks file deletion
func (m *MockFileSystem) DeleteFile(filename string) error {
	args := m.Called(filename)
	return args.Error(0)
}

// FileExists mocks file existence check
func (m *MockFileSystem) FileExists(filename string) bool {
	args := m.Called(filename)
	return args.Bool(0)
}

// MockHTTPClient provides mock HTTP client operations
type MockHTTPClient struct {
	mock.Mock
}

// Get mocks HTTP GET request
func (m *MockHTTPClient) Get(url string) (*MockHTTPResponse, error) {
	args := m.Called(url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockHTTPResponse), args.Error(1)
}

// Post mocks HTTP POST request
func (m *MockHTTPClient) Post(url string, contentType string, body interface{}) (*MockHTTPResponse, error) {
	args := m.Called(url, contentType, body)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockHTTPResponse), args.Error(1)
}

// MockHTTPResponse provides mock HTTP response
type MockHTTPResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
}

// MockTimeProvider provides mock time operations
type MockTimeProvider struct {
	mock.Mock
}

// Now mocks current time
func (m *MockTimeProvider) Now() interface{} {
	args := m.Called()
	return args.Get(0)
}

// Sleep mocks sleep operation
func (m *MockTimeProvider) Sleep(duration interface{}) {
	m.Called(duration)
}

// AnyValue is a helper for matching any value in mock expectations
type AnyValue struct{}

// Match implements driver.Valuer interface
func (a AnyValue) Match(v driver.Value) bool {
	return true
}