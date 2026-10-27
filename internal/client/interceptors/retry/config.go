package retry

type config struct {
	retryer                 Retryer
	failOnExhaustedAttempts bool
	maxAttempts             int
}

func defaultConfig() *config {
	return &config{
		retryer: NoRetry{},
	}
}
