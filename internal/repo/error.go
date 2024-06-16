package repo

import (
	"fmt"

	"github.com/ChillyWR/PasswordManager/pkg/pmerror"
	"gorm.io/gorm"
)

func convertError(err error) error {
	if err == nil {
		return nil
	}

	switch err {
	case gorm.ErrRecordNotFound:
		return pmerror.ErrNotFound
	default:
		return fmt.Errorf("%w: %s", pmerror.ErrInternal, err.Error())
	}
}
