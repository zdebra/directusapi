package directusapi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomMarshal(t *testing.T) {
	primitiveStruct := struct {
		StrVal   string  `directus:"str-val"`
		FloatVal float64 `directus:"float-val"`
		IntVal   int8    `directus:"int-val"`
		UintVal  uint32  `directus:"uint-val"`
	}{
		StrVal:   "abcd",
		FloatVal: 13.2,
		IntVal:   -43,
		UintVal:  2978,
	}
	jsonBytes, err := jsonMarshal(primitiveStruct)
	require.NoError(t, err)

	expectedResult := `{"str-val":"abcd","float-val":13.2,"int-val":-43,"uint-val":2978}`
	assert.Equal(t, expectedResult, string(jsonBytes))

}
