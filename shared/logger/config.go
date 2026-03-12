package logger

// LogType represents the type/category of a log entry.
type LogType string

const (
	// Indicates that a repository operation completed successfully.
	RepositoryOperationSuccess LogType = "REPOSITORY_OPERATION_SUCCESS"
	// Indicates a failure while executing a database query or interacting with the database
	RepositoryQueryFailed LogType = "REPOSITORY_QUERY_FAILED"
	// Indicates a failure while reading data from the database result set
	RepositoryScanFailed LogType = "REPOSITORY_SCAN_FAILED"
	// Indicates that a service-level operation completed successfully.
	ServiceOperationSuccess LogType = "SERVICE_OPERATION_SUCCESS"
	// Indicates a failure during interaction with the repository layer.
	ServiceRepositoryFailed LogType = "SERVICE_REPOSITORY_FAILED"
	// Indicates that the provided data failed validation or did not satisfy service-level business rules.
	ServiceValidationFailed LogType = "SERVICE_VALIDATION_FAILED"
	// Indicates that an HTTP handler operation completed successfully.
	HandlerOperationSuccess LogType = "HANDLER_OPERATION_SUCCESS"
	// Indicates an internal failure while processing the request.
	HandlerInternalError LogType = "HANDLER_INTERNAL_ERROR"
	// Indicates that the request could not be processed due to invalid client input.
	HandlerBadRequest LogType = "HANDLER_BAD_REQUEST"
	// Indicates that a middleware has received an HTTP request.
	MiddlewareRequestReceived LogType = "MIDDLEWARE_REQUEST_RECEIVED"
	// Indicates a failure during logging.
	LoggerError LogType = "LOGGER_ERROR"
)

var config map[LogType]LogTypeConfig

// LogTypeConfig defines where a log of a specific type should be written.
type LogTypeConfig struct {
	ToFile    bool // Whether to write logs to a file
	ToConsole bool // Whether to write logs to the console
}

// LoadDefaultLogConfig initializes the global config map with default settings
// for log types, specifying which logs go to file and/or console.
func LoadDefaultLogConfig() {
	config = make(map[LogType]LogTypeConfig, 10)
	config[RepositoryOperationSuccess] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[RepositoryQueryFailed] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[RepositoryScanFailed] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[ServiceOperationSuccess] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[ServiceRepositoryFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[ServiceValidationFailed] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[HandlerOperationSuccess] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[HandlerInternalError] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[HandlerBadRequest] = LogTypeConfig{ToFile: true, ToConsole: true}
	config[MiddlewareRequestReceived] = LogTypeConfig{ToFile: true, ToConsole: true}
}

func getAllLogTypes() []LogType {
	return []LogType{
		RepositoryOperationSuccess,
		RepositoryQueryFailed,
		RepositoryScanFailed,
		ServiceOperationSuccess,
		ServiceRepositoryFailed,
		ServiceValidationFailed,
		HandlerOperationSuccess,
		HandlerInternalError,
		HandlerBadRequest,
		MiddlewareRequestReceived,
		LoggerError,
	}
}

// LoadAllInclusiveConfig sets all log types to be logged to file and console.
func LoadAllInclusiveConfig() {
	config = make(map[LogType]LogTypeConfig, 10)
	AllLogTypes := getAllLogTypes()

	for _, lType := range AllLogTypes {
		config[lType] = LogTypeConfig{ToFile: true, ToConsole: true}
	}
}

// LoadOnlyConsoleConfig sets all log types to be logged only to console.
func LoadOnlyConsoleConfig() {
	config = make(map[LogType]LogTypeConfig, 10)
	AllLogTypes := getAllLogTypes()

	for _, lType := range AllLogTypes {
		config[lType] = LogTypeConfig{ToFile: false, ToConsole: true}
	}
}
