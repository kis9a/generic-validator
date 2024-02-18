package validator

type Validator[T any] func(T) (bool, error)
type Validators[T any] []Validator[T]
type ValidatedField[T any] struct {
	Value   T
	IsValid bool
	Errors  []error
}

func ApplyValidators[T any](validators Validators[T], value T) ValidatedField[T] {
	validatedField := ValidatedField[T]{Value: value, IsValid: true}
	for _, validator := range validators {
		ok, err := validator(value)
		if !ok {
			validatedField.IsValid = false
			validatedField.Errors = append(validatedField.Errors, err)
		}
	}
	return validatedField
}

func BindValidators[K comparable, T any](validatorsMap map[K]Validators[T]) func(map[K]T) map[K]ValidatedField[T] {
	return func(dataMap map[K]T) map[K]ValidatedField[T] {
		validatedFieldMap := make(map[K]ValidatedField[T])
		for key, validators := range validatorsMap {
			if data, ok := dataMap[key]; ok {
				validatedFieldMap[key] = ApplyValidators(validators, data)
			}
		}
		return validatedFieldMap
	}
}

func Map[K comparable, V any](dataMap map[K]V, fn func(V) V) map[K]V {
	rmap := make(map[K]V)
	for key, data := range dataMap {
		rmap[key] = fn(data)
	}
	return rmap
}

func Every[K comparable, V any](dataMap map[K]V, fn func(V) bool) bool {
	for _, value := range dataMap {
		if !fn(value) {
			return false
		}
	}
	return true
}

func Some[K comparable, V any](dataMap map[K]V, fn func(V) bool) bool {
	for _, value := range dataMap {
		if fn(value) {
			return true
		}
	}
	return false
}

func Filter[K comparable, V any](dataMap map[K]V, fn func(V) bool) map[K]V {
	rmap := make(map[K]V)
	for key, value := range dataMap {
		if fn(value) {
			rmap[key] = value
		}
	}
	return rmap
}

func Reduce[K comparable, V any, R any](dataMap map[K]V, fn func(R, V) R, value R) R {
	rv := value
	for _, value := range dataMap {
		rv = fn(rv, value)
	}
	return rv
}
