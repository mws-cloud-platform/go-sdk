package retry

type config struct {
	retryer                 Retryer
	failOnExhaustedAttempts bool
}

func defaultConfig() *config {
	return &config{
		retryer: NoRetry{},
	}
}
