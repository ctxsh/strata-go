package apex

type ApexError string

const (
	InvalidNamespace   = ApexError("Invalid namepace")
	InvalidSubsystem   = ApexError("Invalid subsystem")
	InvalidMetricName  = ApexError("Invalid metric name")
	RegistrationFailed = ApexError("Unable to register collector")
)

func (e ApexError) Error() string {
	return string(e)
}
