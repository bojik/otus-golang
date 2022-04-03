package sqlstorage

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type Duration struct {
	time.Duration
}

func (d *Duration) Value() (driver.Value, error) {
	return d.Duration / time.Second, nil
}

func (d *Duration) Scan(src interface{}) error {
	if src == nil {
		d.Duration = 0
		return nil
	}
	v, ok := src.(int64)
	if !ok {
		d.Duration = 0
		return nil
	}

	d.Duration = time.Duration(v) * time.Second

	return nil
}

var (
	_ sql.Scanner   = (*Duration)(nil)
	_ driver.Valuer = (*Duration)(nil)
)
