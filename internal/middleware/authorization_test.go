package middleware

import (
	"testing"

	"shortener/internal/logger"
)

func BenchmarkGetUserID(b *testing.B) {
	tokenString := `"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJleHAiOjE3MjY0Mjc5NjksIlVzZXJJRCI6ImUwOTI2OGQ2LWVlMTQtNGU3Yi04MWZiLTUxOGU4M2JmMDM0NiJ9.
yTmWk0mALkC1Lb2j9Qcz70GqY-RA-BOUWX_0e_TbA0U"`
	secretKey := "!@#$%^YdBg0DS"
	log := &logger.Log{}
	log.Initialize("INFO")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		getUserID(tokenString, secretKey, log)
	}
}
