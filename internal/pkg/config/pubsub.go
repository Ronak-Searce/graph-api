package config

// Pubsub ..
type Pubsub struct {
	ADC     string
	Project string
}

func pubsubConfig() Pubsub {
	return Pubsub{
		ADC:     String("env.google_adc"),
		Project: String("env.pubsub.project"),
	}
}
