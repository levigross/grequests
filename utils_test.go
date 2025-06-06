package grequests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContainsString(t *testing.T) {
	// Test case sensitivity (current function is case-sensitive)
	assert.True(t, containsString("Hello", []string{"Hello", "World"}), "Expected to find 'Hello'")
	assert.False(t, containsString("hello", []string{"Hello", "World"}), "Expected not to find 'hello' (case-sensitive)")

	// Test with empty slice
	assert.False(t, containsString("Hello", []string{}), "Expected not to find in empty slice")

	// Test with nil slice
	assert.False(t, containsString("Hello", nil), "Expected not to find in nil slice")

	// Matches at beginning, middle, end
	assert.True(t, containsString("A", []string{"A", "B", "C"}), "Expected to find 'A' at beginning")
	assert.True(t, containsString("B", []string{"A", "B", "C"}), "Expected to find 'B' in middle")
	assert.True(t, containsString("C", []string{"A", "B", "C"}), "Expected to find 'C' at end")

	// No match
	assert.False(t, containsString("D", []string{"A", "B", "C"}), "Expected not to find 'D'")
}

func TestIsStringInSlice(t *testing.T) {
	// Test case sensitivity (current function is case-sensitive)
	assert.True(t, isStringInSlice("Hello", []string{"Hello", "World"}), "Expected to find 'Hello'")
	assert.False(t, isStringInSlice("hello", []string{"Hello", "World"}), "Expected not to find 'hello' (case-sensitive)")

	// Test with empty slice
	assert.False(t, isStringInSlice("Hello", []string{}), "Expected not to find in empty slice")

	// Test with nil slice
	assert.False(t, isStringInSlice("Hello", nil), "Expected not to find in nil slice")

	// Matches at beginning, middle, end
	assert.True(t, isStringInSlice("A", []string{"A", "B", "C"}), "Expected to find 'A' at beginning")
	assert.True(t, isStringInSlice("B", []string{"A", "B", "C"}), "Expected to find 'B' in middle")
	assert.True(t, isStringInSlice("C", []string{"A", "B", "C"}), "Expected to find 'C' at end")

	// No match
	assert.False(t, isStringInSlice("D", []string{"A", "B", "C"}), "Expected not to find 'D'")
}

func TestIsRedirect(t *testing.T) {
	redirectCodes := []int{301, 302, 303, 307, 308}
	for _, code := range redirectCodes {
		assert.True(t, isRedirect(code), fmt.Sprintf("Expected status code %d to be a redirect", code))
	}

	nonRedirectCodes := []int{200, 201, 400, 404, 500, 502}
	for _, code := range nonRedirectCodes {
		assert.False(t, isRedirect(code), fmt.Sprintf("Expected status code %d not to be a redirect", code))
	}
}
