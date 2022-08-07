package directusapi

import (
	"bytes"
	"encoding/json"
)

type KeyVal struct {
	Key string
	Val interface{}
}

// OrderedMap guarantees orded for JSON marshaling and unsmarshaling
type OrderedMap []KeyVal

func (omap OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	if _, err := buf.WriteString("{"); err != nil {
		return nil, err
	}
	for i, kv := range omap {
		if i != 0 {
			if _, err := buf.WriteString(","); err != nil {
				return nil, err
			}
		}

		// marshal key
		key, err := json.Marshal(kv.Key)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(key); err != nil {
			return nil, err
		}
		if _, err := buf.WriteString(":"); err != nil {
			return nil, err
		}

		// marshal value
		val, err := json.Marshal(kv.Val)
		if err != nil {
			return nil, err
		}
		if _, err := buf.Write(val); err != nil {
			return nil, err
		}
	}

	if _, err := buf.WriteString("}"); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
