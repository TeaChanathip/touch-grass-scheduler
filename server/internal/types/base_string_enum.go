package types

import (
	"database/sql/driver"
	"fmt"
)

// BaseStringEnum provides common Scan/Value methods
type BaseStringEnum string

func (e *BaseStringEnum) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("cannot scan %T into BaseStringEnum", value)
	}

	*e = BaseStringEnum(str)
	return nil
}

func (e BaseStringEnum) Value() (driver.Value, error) {
	return string(e), nil
}
