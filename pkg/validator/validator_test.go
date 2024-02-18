package validator

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

type TestCaseBindValidators[K comparable, T any] struct {
	name          string
	validatorsMap map[K]Validators[T]
	dataMap       map[K]T
	expectedValid bool
}

func runTestCasesBindValidators[K comparable, T any](t *testing.T, tests []TestCaseBindValidators[K, T]) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validateDataMap := BindValidators(tt.validatorsMap)
			validatedDataMap := validateDataMap(tt.dataMap)
			isValidDataMap := Every[K, ValidatedField[T]](validatedDataMap, func(field ValidatedField[T]) bool {
				return field.IsValid
			})

			if isValidDataMap != tt.expectedValid {
				t.Errorf("Expected dataMap validity to be %v, got %v", tt.expectedValid, isValidDataMap)
			}
		})
	}
}

func TestBindValidators(t *testing.T) {
	runTestCasesBindValidators(t, []TestCaseBindValidators[string, string]{
		{
			name: "Valid string values",
			validatorsMap: map[string]Validators[string]{
				"name": {validateLength, validateCharacter},
			},
			dataMap: map[string]string{
				"name": "helloWorld",
			},
			expectedValid: true,
		},
		{
			name: "Invalid string value - name too short",
			validatorsMap: map[string]Validators[string]{
				"name": {validateLength, validateCharacter},
			},
			dataMap: map[string]string{
				"name": "hell",
			},
			expectedValid: false,
		},
		{
			name: "Invalid name - contains numbers",
			validatorsMap: map[string]Validators[string]{
				"name": {validateLength, validateCharacter},
			},
			dataMap: map[string]string{
				"name": "hellowor8d",
			},
			expectedValid: false,
		},
	})
	runTestCasesBindValidators(t, []TestCaseBindValidators[string, int]{
		{
			name: "Valid int values",
			validatorsMap: map[string]Validators[int]{
				"age": {validatePositive},
			},
			dataMap: map[string]int{
				"age": 30,
			},
			expectedValid: true,
		},
		{
			name: "Invalid int value - age negative",
			validatorsMap: map[string]Validators[int]{
				"age": {validatePositive},
			},
			dataMap: map[string]int{
				"age": -1,
			},
			expectedValid: false,
		},
	})
}

func TestEvery(t *testing.T) {
	var tests = []struct {
		name     string
		dataMap  map[string]ValidatedField[string]
		expected bool
	}{
		{
			name: "All fields valid",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "test", IsValid: true},
				"b": {Value: "example", IsValid: true},
			},
			expected: true,
		},
		{
			name: "Some fields invalid",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "test", IsValid: false},
				"b": {Value: "example", IsValid: true},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValidDataMap := Every(tt.dataMap, func(v ValidatedField[string]) bool {
				return v.IsValid
			})
			if isValidDataMap != tt.expected {
				t.Errorf("Expected result for '%s': %v, got: %v", tt.name, tt.expected, isValidDataMap)
			}
		})
	}
}

func TestSome(t *testing.T) {
	tests := []struct {
		name     string
		dataMap  map[string]ValidatedField[string]
		expected bool
	}{
		{
			name: "At least one valid field",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "test", IsValid: false},
				"b": {Value: "example", IsValid: true},
			},
			expected: true,
		},
		{
			name: "No valid fields",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "fail", IsValid: false},
				"b": {Value: "fail", IsValid: false},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Some(tt.dataMap, func(v ValidatedField[string]) bool {
				return v.IsValid
			})
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		dataMap  map[string]ValidatedField[string]
		expected int
	}{
		{
			name: "Filter out invalid fields",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "test", IsValid: false},
				"b": {Value: "example", IsValid: true},
			},
			expected: 1,
		},
		{
			name: "All fields are valid",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "hello", IsValid: true},
				"b": {Value: "world", IsValid: true},
			},
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.dataMap, func(v ValidatedField[string]) bool {
				return v.IsValid
			})
			if len(result) != tt.expected {
				t.Errorf("Expected %d valid fields, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		dataMap  map[string]ValidatedField[string]
		reducer  func(string, ValidatedField[string]) string
		init     string
		expected string
	}{
		{
			name: "Concatenate valid fields",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "hello", IsValid: true},
				"b": {Value: "world", IsValid: true},
			},
			reducer: func(acc string, v ValidatedField[string]) string {
				return acc + v.Value
			},
			init:     "",
			expected: "helloworld",
		},
		{
			name: "Concatenate with separator",
			dataMap: map[string]ValidatedField[string]{
				"a": {Value: "hello", IsValid: true},
				"b": {Value: "world", IsValid: true},
			},
			reducer: func(acc string, v ValidatedField[string]) string {
				if acc != "" {
					return acc + ", " + v.Value
				}
				return v.Value
			},
			init:     "",
			expected: "hello, world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.dataMap, tt.reducer, tt.init)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func validateLength(v string) (bool, error) {
	if len(v) > 8 {
		return true, nil
	}
	return false, errors.New("Must be longer than 8 characters")
}

func validatePositive(v int) (bool, error) {
	if v > 0 {
		return true, nil
	}
	return false, errors.New("Value must be positive")
}

func validateTrue(v bool) (bool, error) {
	if v {
		return true, nil
	}
	return false, errors.New("Value must be true")
}

func validateEmptyString(v string) (bool, error) {
	return false, errors.New("Must be longer than 8 characters")
}

func validateCharacter(s string) (bool, error) {
	pattern := "^[A-Za-z]+$"
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false, fmt.Errorf("Failed to compile regexp string, %s", pattern)
	}
	if matched {
		return true, nil
	}
	return false, errors.New("Must be an alphabetic string")
}
