package env

import (
	"context"
	"os/user"
)

// HasUser is checking if the current user is available.
func HasUser() Check {
	return func(ctx context.Context) error {
		_, err := user.Current()
		if err != nil {
			return err
		}

		return nil
	}
}
