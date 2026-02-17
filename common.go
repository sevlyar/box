package box

import (
	"encoding"
	"encoding/json"
)

func marshalJSON[T any](v T) ([]byte, error) {
	if casted, ok := any(v).(json.Marshaler); ok {
		return casted.MarshalJSON()
	}

	if casted, ok := any(v).(encoding.TextMarshaler); ok {
		return casted.MarshalText()
	}

	return json.Marshal(v)
}

func unmarshalJSON[T any](data []byte, ptr *T) error {
	if casted, ok := any(ptr).(json.Unmarshaler); ok {
		return casted.UnmarshalJSON(data)
	}

	if casted, ok := any(ptr).(encoding.TextUnmarshaler); ok {
		return casted.UnmarshalText(data)
	}

	return json.Unmarshal(data, ptr)
}
