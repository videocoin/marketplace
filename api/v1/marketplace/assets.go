package marketplace

import (
	"database/sql/driver"
	"errors"
)

func (s AssetStatus) Value() (driver.Value, error) {
	return AssetStatus_name[int32(s)], nil
}

func (s *AssetStatus) Scan(src interface{}) error {
	id, ok := src.(string)
	if !ok {
		return errors.New("type assertion .(string) failed.")
	}

	*s = AssetStatus(AssetStatus_value[id])

	return nil
}
