package mysql

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// UUID is a MySQL-specific type that stores UUIDs as binary(16).
// PostgreSQL users should use uuid.UUID directly instead.
type UUID uuid.UUID

func NewUUID() UUID {
	return UUID(uuid.New())
}

func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID(id), err
}

func MustParseUUID(s string) UUID {
	return UUID(uuid.MustParse(s))
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

func (id UUID) IsZero() bool {
	return uuid.UUID(id) == uuid.Nil
}

func (id UUID) GormDataType() string {
	return "binary(16)"
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func (id *UUID) UnmarshalJSON(by []byte) error {
	str := string(by)
	if len(str) == 0 || str == "null" {
		*id = UUID{}
		return nil
	}
	if len(str) >= 2 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	if len(str) == 0 {
		*id = UUID{}
		return nil
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return err
	}
	*id = UUID(u)
	return nil
}

func (id *UUID) Scan(value any) error {
	if value == nil {
		*id = UUID{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan type %T into mysql.UUID", value)
	}
	parsed, err := uuid.FromBytes(bytes)
	if err != nil {
		return err
	}
	*id = UUID(parsed)
	return nil
}

func (id UUID) Value() (driver.Value, error) {
	return uuid.UUID(id).MarshalBinary()
}
