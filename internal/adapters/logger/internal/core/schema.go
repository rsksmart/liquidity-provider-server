package core

// TimestampLayout renders ISO-8601 with a timezone offset and millisecond
// precision (e.g. 2026-04-14T15:30:00.123Z for UTC). Layout is hardcoded as
// this specific version is not part of the standard library.
const TimestampLayout = "2006-01-02T15:04:05.000Z07:00"

// CanonicalField represents the canonical field names. The base fields are mandatory on every record.
type CanonicalField string

const (
	FieldTimestamp   CanonicalField = "timestamp"
	FieldLevel       CanonicalField = "level"
	FieldMessage     CanonicalField = "message"
	FieldService     CanonicalField = "service"
	FieldEnvironment CanonicalField = "environment"
	FieldVersion     CanonicalField = "version"
	FieldTraceID     CanonicalField = "traceId"
	FieldSpanID      CanonicalField = "spanId"
)

// ErrorField represents the error field names. The error fields are well-known optional fields defined by the standard.
type ErrorField string

const (
	FieldErrorMessage ErrorField = "errorMessage"
	FieldErrorCode    ErrorField = "errorCode"
	FieldErrorStack   ErrorField = "errorStack"
	FieldErrorContext ErrorField = "errorContext"
)
