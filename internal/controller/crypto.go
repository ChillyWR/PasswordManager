package controller

import (
	"fmt"

	"github.com/okutsen/PasswordManager/model"
	"github.com/okutsen/PasswordManager/pkg/pmcrypto"
)

func (c *Controller) encryptCredentialRecord(record *model.CredentialRecord) error {
	if record.Notes != nil {
		v, err := pmcrypto.Encrypt(*record.Notes, Salt)
		if err != nil {
			return err
		}

		record.Notes = &v
	}

	return nil
}

func (c *Controller) encryptLogin(record *model.LoginRecord) error {
	if err := c.encryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("encrypt core: %w", err)
	}

	if record.Username != nil {
		v, err := pmcrypto.Encrypt(*record.Username, Salt)
		if err != nil {
			return err
		}

		record.Username = &v
	}

	if record.Password != nil {
		v, err := pmcrypto.Encrypt(*record.Password, Salt)
		if err != nil {
			return err
		}

		record.Password = &v
	}

	return nil
}

func (c *Controller) encryptCard(record *model.CardRecord) error {
	if err := c.encryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("encrypt core: %w", err)
	}

	if record.Number != nil {
		v, err := pmcrypto.Encrypt(*record.Number, Salt)
		if err != nil {
			return err
		}

		record.Number = &v
	}

	if record.ExpirationMonth != nil {
		v, err := pmcrypto.Encrypt(*record.ExpirationMonth, Salt)
		if err != nil {
			return err
		}

		record.ExpirationMonth = &v
	}

	if record.ExpirationYear != nil {
		v, err := pmcrypto.Encrypt(*record.ExpirationYear, Salt)
		if err != nil {
			return err
		}

		record.ExpirationYear = &v
	}

	if record.CVV != nil {
		v, err := pmcrypto.Encrypt(*record.CVV, Salt)
		if err != nil {
			return err
		}

		record.CVV = &v
	}

	return nil
}

func (c *Controller) encryptIdentity(record *model.IdentityRecord) error {
	if err := c.encryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("encrypt core: %w", err)
	}

	if record.FirstName != nil {
		v, err := pmcrypto.Encrypt(*record.FirstName, Salt)
		if err != nil {
			return err
		}

		record.FirstName = &v
	}

	if record.MiddleName != nil {
		v, err := pmcrypto.Encrypt(*record.MiddleName, Salt)
		if err != nil {
			return err
		}

		record.MiddleName = &v
	}

	if record.LastName != nil {
		v, err := pmcrypto.Encrypt(*record.LastName, Salt)
		if err != nil {
			return err
		}

		record.LastName = &v
	}

	if record.Address != nil {
		v, err := pmcrypto.Encrypt(*record.Address, Salt)
		if err != nil {
			return err
		}

		record.Address = &v
	}

	if record.Email != nil {
		v, err := pmcrypto.Encrypt(*record.Email, Salt)
		if err != nil {
			return err
		}

		record.Email = &v
	}

	if record.PhoneNumber != nil {
		v, err := pmcrypto.Encrypt(*record.PhoneNumber, Salt)
		if err != nil {
			return err
		}

		record.PhoneNumber = &v
	}

	if record.PassportNumber != nil {
		v, err := pmcrypto.Encrypt(*record.PassportNumber, Salt)
		if err != nil {
			return err
		}

		record.PassportNumber = &v
	}

	return nil
}

func (c *Controller) decryptCredentialRecord(record *model.CredentialRecord) error {
	if record.Notes != nil {
		v, err := pmcrypto.Decrypt(*record.Notes, Salt)
		if err != nil {
			return err
		}

		record.Notes = &v
	}

	return nil
}

func (c *Controller) decryptLogin(record *model.LoginRecord) error {
	if err := c.decryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("decrypt core: %w", err)
	}

	if record.Username != nil {
		v, err := pmcrypto.Decrypt(*record.Username, Salt)
		if err != nil {
			return err
		}

		record.Username = &v
	}

	if record.Password != nil {
		v, err := pmcrypto.Decrypt(*record.Password, Salt)
		if err != nil {
			return err
		}

		record.Password = &v
	}

	return nil
}

func (c *Controller) decryptCard(record *model.CardRecord) error {
	if err := c.decryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("decrypt core: %w", err)
	}

	if record.Number != nil {
		v, err := pmcrypto.Decrypt(*record.Number, Salt)
		if err != nil {
			return err
		}

		record.Number = &v
	}

	if record.ExpirationMonth != nil {
		v, err := pmcrypto.Decrypt(*record.ExpirationMonth, Salt)
		if err != nil {
			return err
		}

		record.ExpirationMonth = &v
	}

	if record.ExpirationYear != nil {
		v, err := pmcrypto.Decrypt(*record.ExpirationYear, Salt)
		if err != nil {
			return err
		}

		record.ExpirationYear = &v
	}

	if record.CVV != nil {
		v, err := pmcrypto.Decrypt(*record.CVV, Salt)
		if err != nil {
			return err
		}

		record.CVV = &v
	}

	return nil
}

func (c *Controller) decryptIdentity(record *model.IdentityRecord) error {
	if err := c.decryptCredentialRecord(&record.CredentialRecord); err != nil {
		return fmt.Errorf("decrypt core: %w", err)
	}

	if record.FirstName != nil {
		v, err := pmcrypto.Decrypt(*record.FirstName, Salt)
		if err != nil {
			return err
		}

		record.FirstName = &v
	}

	if record.MiddleName != nil {
		v, err := pmcrypto.Decrypt(*record.MiddleName, Salt)
		if err != nil {
			return err
		}

		record.MiddleName = &v
	}

	if record.LastName != nil {
		v, err := pmcrypto.Decrypt(*record.LastName, Salt)
		if err != nil {
			return err
		}

		record.LastName = &v
	}

	if record.Address != nil {
		v, err := pmcrypto.Decrypt(*record.Address, Salt)
		if err != nil {
			return err
		}

		record.Address = &v
	}

	if record.Email != nil {
		v, err := pmcrypto.Decrypt(*record.Email, Salt)
		if err != nil {
			return err
		}

		record.Email = &v
	}

	if record.PhoneNumber != nil {
		v, err := pmcrypto.Decrypt(*record.PhoneNumber, Salt)
		if err != nil {
			return err
		}

		record.PhoneNumber = &v
	}

	if record.PassportNumber != nil {
		v, err := pmcrypto.Decrypt(*record.PassportNumber, Salt)
		if err != nil {
			return err
		}

		record.PassportNumber = &v
	}

	return nil
}
