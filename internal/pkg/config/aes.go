package config

// AES ..
type AES struct {
	Secret string
	Salt   string
}

func aesConfig() AES {
	return AES{
		Secret: String("env.aes.secret"),
		Salt:   String("env.aes.salt"),
	}
}
