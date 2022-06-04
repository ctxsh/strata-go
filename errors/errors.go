package errors

type ApexError string

const (
	ErrInvalidMetricName  = ApexError("Invalid metric name")
	ErrRegistrationFailed = ApexError("Unable to register collector")
	ErrAlreadyRegistered  = ApexError("metric is already registered")
)

func (e ApexError) Error() string {
	return string(e)
}
