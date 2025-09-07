package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	// Base62 characters for URL encoding
	base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomCode generates a random alphanumeric string of length n
func GenerateRandomCode(n int) string {
	fmt.Println("GenerateRandomCode called with length:", n)
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = base62Chars[rand.Intn(len(base62Chars))]
	}
	return string(b)
}

// ValidateURL performs basic URL validation
func ValidateURL(url string) bool {
	fmt.Println("ValidateURL called with url:", url)
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// NormalizeURL adds https:// prefix if no protocol is specified
func NormalizeURL(url string) string {
	fmt.Println("NormalizeURL called with url:", url)
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		fmt.Println("Normalized URL:", "https://"+url)
		return "https://" + url
	}
	fmt.Println("Normalized URL:", url)
	return url
}
