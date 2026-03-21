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
	ServiceDataFetchFailed  LogType = "SERVICE_DATA_FETCH_FAILED"
	EventDataReadFailed     LogType = "EVENT_DATA_READ_FAILED"
	EventDataReceived       LogType = "EVENT_DATA_RECEIVED"
	EventDataSentFailed     LogType = "EVENT_DATA_SENT_FAILED"
	EventDataSentSuccess    LogType = "EVENT_DATA_SENT_SUCCESS"
	// Indicates that an HTTP handler operation completed successfully.
	HandlerOperationSuccess LogType = "HANDLER_OPERATION_SUCCESS"
	// Indicates an internal failure while processing the request.
	HandlerInternalError LogType = "HANDLER_INTERNAL_ERROR"
	// Indicates that the request could not be processed due to invalid client input.
	HandlerBadRequest LogType = "HANDLER_BAD_REQUEST"
	// Indicates that a middleware has received an HTTP request.
	MiddlewareRequestReceived LogType = "MIDDLEWARE_REQUEST_RECEIVED"
	// Indicates that a api has finished an HTTP request.
	MiddlewareRequestFinished LogType = "MIDDLEWARE_REQUEST_FINISHED"
	// MiddlewareFailed indicates that an error occurred while processing a request in a middleware.
	MiddlewareFailed LogType = "MIDDLEWARE_ERROR"
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
	config[RepositoryQueryFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[RepositoryScanFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[ServiceOperationSuccess] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[ServiceRepositoryFailed] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[ServiceValidationFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[ServiceDataFetchFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[EventDataReadFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[EventDataReceived] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[EventDataSentSuccess] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[EventDataSentFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[HandlerOperationSuccess] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[HandlerInternalError] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[HandlerBadRequest] = LogTypeConfig{ToFile: false, ToConsole: false}
	config[MiddlewareRequestReceived] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[MiddlewareRequestFinished] = LogTypeConfig{ToFile: false, ToConsole: true}
	config[MiddlewareFailed] = LogTypeConfig{ToFile: false, ToConsole: true}
}

func getAllLogTypes() []LogType {
	return []LogType{
		RepositoryOperationSuccess,
		RepositoryQueryFailed,
		RepositoryScanFailed,
		ServiceOperationSuccess,
		ServiceRepositoryFailed,
		ServiceValidationFailed,
		ServiceDataFetchFailed,
		EventDataReadFailed,
		EventDataReceived,
		EventDataSentFailed,
		EventDataSentSuccess,
		HandlerOperationSuccess,
		HandlerInternalError,
		HandlerBadRequest,
		MiddlewareRequestReceived,
		MiddlewareRequestFinished,
		MiddlewareFailed,
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
