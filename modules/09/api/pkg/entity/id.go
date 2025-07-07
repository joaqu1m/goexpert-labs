package entity

import (
	"database/sql/driver"
	"errors"

	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(s string) (ID, error) {
	id, err := uuid.Parse(s)
	return ID(id), err
}

func (id ID) Value() (driver.Value, error) {
	return uuid.UUID(id).String(), nil
}

func (id *ID) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		parsed, err := uuid.ParseBytes(v)
		if err != nil {
			return err
		}
		*id = ID(parsed)
		return nil
	case string:
		parsed, err := uuid.Parse(v)
		if err != nil {
			return err
		}
		*id = ID(parsed)
		return nil
	default:
		return errors.New("invalid type for ID")
	}
}
