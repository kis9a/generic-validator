## Example

```go
package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/kis9a/generic-validator/pkg/validator"
)

func validateLength(v string) (bool, error) {
	if len(v) > 8 {
		return true, nil
	}
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

func main() {
	validatorsMap := map[string]validator.Validators[string]{
		"name": {validateLength, validateCharacter},
	}
	dataMap := map[string]string{
		"name": "hellowor8d",
	}
	validateDataMap := validator.BindValidators(validatorsMap)
	validatedDataMap := validateDataMap(dataMap)
	isValidDataMap := validator.Every(validatedDataMap, func(field validator.ValidatedField[string]) bool { return field.IsValid })
	if !isValidDataMap {
		validator.Map(validatedDataMap, func(field validator.ValidatedField[string]) validator.ValidatedField[string] {
			for _, err := range field.Errors {
				fmt.Printf("ERROR: %s\n", err.Error())
			}
			return field
		})
	}
}
```
