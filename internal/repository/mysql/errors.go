package mysql

import (
	"errors"

	driver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"api-admin-panel/internal/domain"
)

const (
	mysqlErrDuplicateEntry  = 1062
	mysqlErrRowIsReferenced = 1451 // hapus/ubah baris induk yang masih direferensikan
	mysqlErrNoReferencedRow = 1452 // insert/ubah baris anak ke induk yang tidak ada
)

func mysqlErrNumber(err error) uint16 {
	var me *driver.MySQLError
	if errors.As(err, &me) {
		return me.Number
	}
	return 0
}

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	return mysqlErrNumber(err) == mysqlErrDuplicateEntry
}

// translateWriteError menerjemahkan error tulis MySQL (duplicate key, foreign key)
// menjadi domain error agar handler bisa memberi status 4xx, bukan 500.
// Error lain dikembalikan apa adanya.
func translateWriteError(err error) error {
	switch {
	case err == nil:
		return nil
	case isDuplicateKey(err):
		return domain.ErrAlreadyExists
	case mysqlErrNumber(err) == mysqlErrNoReferencedRow:
		return domain.ErrInvalidReference
	case mysqlErrNumber(err) == mysqlErrRowIsReferenced:
		return domain.ErrInUse
	default:
		return err
	}
}
