package optional_test

import (
	"bytes"
	"github.com/clarketm/json"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/stretchr/testify/require"
)

func TestFunc_Optional(t *testing.T) {
	t.Run(
		"Create",
		func(t *testing.T) {
			op := optional.NewOptional("some_value")

			require.NotNil(t, op.Value)
			require.True(t, op.HasValue())

			require.Equal(
				t,
				"some_value",
				op.MustGetValue(),
			)

			getValue, err := op.GetValue()
			require.Equal(
				t,
				"some_value",
				getValue,
			)
			require.NoError(t, err)

			require.Equal(
				t,
				op.GetPointerValue(),
				&getValue,
			)
		},
	)

	t.Run(
		"CreateEmpty",
		func(t *testing.T) {
			op := optional.NewEmptyOptional[string]()
			require.False(
				t,
				op.HasValue(),
			)

			_, err := op.GetValue()
			require.ErrorIs(
				t,
				err,
				optional.NoValue,
			)

			require.Nil(
				t,
				op.GetPointerValue(),
			)
		},
	)

	type someStruct struct {
		Field optional.Optional[string] `json:"field"`
	}

	t.Run(
		"Marshall",
		func(t *testing.T) {
			t.Run(
				"OptionalNotNil",
				func(t *testing.T) {
					str := "value_1"
					s := someStruct{
						Field: optional.Optional[string]{Value: &str},
					}
		
					data, err := json.Marshal(s)
					require.NoError(t, err)
					require.Equal(t, `{"field":"value_1"}`, string(data))
				},
			)

			t.Run(
				"OptionalNil",
				func(t *testing.T) {
					s := someStruct{}
					data, err := json.Marshal(s)
					require.NoError(t, err)

					require.Equal(t, `{"field":null}`, string(data))
				},
			)

			t.Run(
				"OptionalNilAndOmitempty",
				func(t *testing.T) {
					type Struct struct {
						Field optional.Optional[string] `json:"field,omitempty"`
					}

					s := Struct{
						Field: optional.Optional[string]{},
					}
					b := bytes.NewBuffer([]byte{})
					err := json.NewEncoder(b).Encode(s)
					require.NoError(t, err)
				},
			)
		},
	)

	t.Run(
		"Unmarshal",
		func(t *testing.T) {
			t.Run(
				"OptionalNotNil",
				func(t *testing.T) {
					data := `{"field":"value"}`
					s := someStruct{}
					err := json.Unmarshal([]byte(data), &s)
					require.NoError(t, err)
					require.Equal(t, "value", *s.Field.Value)
				},
			)

			t.Run(
				"OptionalNil",
				func(t *testing.T) {
					data := `{"field":null}`
					s := someStruct{}
					err := json.Unmarshal([]byte(data), &s)
					require.NoError(t, err)
					require.Nil(t, s.Field.Value)
				},
			)
		},
	)
}
