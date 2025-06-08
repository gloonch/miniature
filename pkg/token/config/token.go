package config

import (
	"time"
)

var (
	JWTSecret = []byte("your_secret_key_here")
	TokenTTL  = time.Hour * 72 // 72 hours expiration
)
