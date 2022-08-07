package directusapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type PrimitiveStruct struct {
	StrVal   string  `directus:"str-val"`
	FloatVal float64 `directus:"float-val"`
	IntVal   int8    `directus:"int-val"`
	UintVal  uint32  `directus:"uint-val"`
}

type SliceStruct struct {
	StringSlice []string          `directus:"str-slice-val"`
	IntSlice    []int             `directus:"int-slice-val"`
	BoolSlice   []bool            `directus:"bool-slice-val"`
	StructSlice []PrimitiveStruct `directus:"struct-slice-val"`
	// todo: add slice of pointers
}

func TestCustomMarshal(t *testing.T) {
	t.Run("primitive struct", func(t *testing.T) {
		primitiveStruct := PrimitiveStruct{
			StrVal:   "abcd",
			FloatVal: 13.2,
			IntVal:   -43,
			UintVal:  2978,
		}
		jsonBytes, err := jsonMarshal(primitiveStruct)
		require.NoError(t, err)

		expectedResult := `{"str-val":"abcd","float-val":13.2,"int-val":-43,"uint-val":2978}`
		assert.Equal(t, expectedResult, string(jsonBytes))
	})

	t.Run("struct with slice", func(t *testing.T) {
		sliceStruct := SliceStruct{
			StringSlice: []string{"a", "b", "c"},
			IntSlice:    []int{12, 45, 2},
			BoolSlice:   []bool{true, true, false, true},
			StructSlice: []PrimitiveStruct{
				{
					StrVal:   "1",
					FloatVal: 1.1,
					IntVal:   1,
					UintVal:  1,
				},
				{
					StrVal:   "2",
					FloatVal: 2.2,
					IntVal:   2,
					UintVal:  2,
				},
			},
		}
		jsonBytes, err := jsonMarshal(sliceStruct)
		require.NoError(t, err)

		expectedResult := `{"str-slice-val":["a","b","c"],"int-slice-val":[12,45,2],"bool-slice-val":[true,true,false,true],"struct-slice-val":[{"str-val":"1","float-val":1.1,"int-val":1,"uint-val":1},{"str-val":"2","float-val":2.2,"int-val":2,"uint-val":2}]}`
		assert.Equal(t, expectedResult, string(jsonBytes))
	})
}
