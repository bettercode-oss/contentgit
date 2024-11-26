package persistence

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
)

type JSONB map[string]any

// Value Marshal
func (jsonField *JSONB) Value() (driver.Value, error) {
	return json.Marshal(jsonField)
}

// Scan Unmarshal
func (jsonField *JSONB) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(data, &jsonField)
}
