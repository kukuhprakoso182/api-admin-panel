package mysql

import (
	"errors"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var me *driver.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}