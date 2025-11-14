package api

import (
	"net/http"
	"regexp"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidationMiddleware provides request validation
func ValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("validator", validate)
		c.Next()
	}
}

// ValidateStruct validates a struct against validation tags
func ValidateStruct(obj interface{}) error {
	return validate.Struct(obj)
}

// Custom validation functions
func init() {
	// Register custom validators if needed
	validate.RegisterValidation("strong_password", validateStrongPassword)
}

// validateStrongPassword checks if password meets minimum security requirements
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	
	// At least 8 characters
	if len(password) < 8 {
		return false
	}
	
	// Contains at least one uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return false
	}
	
	// Contains at least one lowercase letter
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return false
	}
	
	// Contains at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	if !hasDigit {
		return false
	}
	
	// Contains at least one special character
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)
	if !hasSpecial {
		return false
	}
	
	return true
}

// SanitizeInput middleware to sanitize input data
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				values[i] = strings.TrimSpace(value)
			}
			c.Request.URL.Query()[key] = values
		}
		
		c.Next()
	}
}

// RateLimitMiddleware (simple implementation)
func RateLimitMiddleware() gin.HandlerFunc {
	// This is a placeholder for rate limiting
	// In production, you'd want to use something like redis-based rate limiting
	return func(c *gin.Context) {
		// Simple rate limiting logic would go here
		c.Next()
	}
}

// SecurityHeaders middleware adds security headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// In production, you might want to add CSP headers
		// c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}