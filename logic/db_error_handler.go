package logic

import (
	"errors"
	"fmt"

	"github.com/copernet/whcwallet/api/errs"
	"github.com/go-sql-driver/mysql"
)

func process(err error, args ...interface{}) error {
	if err == nil {
		return nil
	}

	e, ok := err.(*mysql.MySQLError)
	if ok && e.Number == 1062 {
		return errors.New(fmt.Sprintf(errs.DuplicateEmail, args))
	}

	return errors.New(errs.DefaultMessage)
}
