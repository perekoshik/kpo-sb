package config

type Config struct {
	InteractiveKindnessThreshold int

	MaxFoodPerDayForAdmission int
}

func Default() Config {
	return Config{
		InteractiveKindnessThreshold: 6,
		MaxFoodPerDayForAdmission:    0,
	}
}

func (c Config) AdmissionLimitEnabled() bool {
	return c.MaxFoodPerDayForAdmission > 0
}
