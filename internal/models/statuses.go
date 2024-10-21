package models

const (
	ServiceStatusNew        = "NEW"
	ServiceStatusProcessing = "PROCESSING"
	ServiceStatusInvalid    = "INVALID"
	ServiceStatusProcessed  = "PROCESSED"
)

var ServiceStatusesInProgress = []string{
	ServiceStatusNew,
	ServiceStatusProcessing,
}

const (
	AccrualStatusRegistered = "REGISTERED"
	AccrualStatusInvalid    = "INVALID"
	AccrualStatusProcessing = "PROCESSING"
	AccrualStatusProcessed  = "PROCESSED"
)

var StatusesMap = map[string]string{
	AccrualStatusRegistered: ServiceStatusProcessing,
	AccrualStatusProcessing: ServiceStatusProcessing,
	AccrualStatusInvalid:    ServiceStatusInvalid,
	AccrualStatusProcessed:  ServiceStatusProcessed,
}
