package optional

import (
	"github.com/clarketm/json"
	"errors"
)

var (
	NoValue = errors.New("Has no value")
)

type Optional[T any] struct {
	Value *T `json:",omitempty"`
}

func NewOptional[T any](value T) *Optional[T] {
	return &Optional[T]{
		Value: &value,
	}
}

func NewEmptyOptional[T any]() *Optional[T] {
	return &Optional[T] {
	}
}

func (o Optional[T]) HasValue() bool {
	return o.Value != nil
}

func (o *Optional[T]) SetValue(value T) {
	o.Value = &value
}

// GetValue if not have value return error
// 	NoValue
func (o *Optional[T]) GetValue() (T, error) {
	if o.HasValue() {
		return *o.Value, nil
	} else {
		return *new(T), NoValue
	}
}

// MustGetValue return value if have
// if not return a empty struct
func (o *Optional[T]) MustGetValue() T {
	if o.HasValue() {
		return *o.Value
	} else {
		return *new(T)
	}
}

// Return Pointer so can be nil
func (o *Optional[T]) GetPointerValue() *T {
	return o.Value
}

func (o Optional[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	val := new(T)
	err := json.Unmarshal(data, &val)
	if err != nil {
		return err
	}
	
	o.Value = val
	return nil
}





