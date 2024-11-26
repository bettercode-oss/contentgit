package serializer

import (
	jsoniter "github.com/json-iterator/go"
	"io"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Unmarshal(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return json.NewDecoder(r)
}

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return json.NewEncoder(w)
}
