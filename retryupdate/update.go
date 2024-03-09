//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

type GetResponseWrapper struct {
	Value   *string
	Version uuid.UUID
}

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	var getResp GetResponseWrapper
	var authErr *kvapi.AuthError

loop:
	for {
		res, err := c.Get(&kvapi.GetRequest{Key: key})

		switch true {
		case errors.As(err, &authErr):
			return err

		case err == nil:
			getResp.Value = &res.Value
			getResp.Version = res.Version
			fallthrough

		case errors.Is(err, kvapi.ErrKeyNotFound):
			break loop
		}
	}

	var conflictErr *kvapi.ConflictError
	newVersion := uuid.Must(uuid.NewV4())

	for {
		updatedValue, err := updateFn(getResp.Value)

		if err != nil {
			return err
		}

		_, err = c.Set(&kvapi.SetRequest{Key: key, Value: updatedValue, OldVersion: getResp.Version, NewVersion: newVersion})

		switch true {
		case errors.As(err, &authErr) || err == nil:
			return err

		case errors.Is(err, kvapi.ErrKeyNotFound):
			getResp.Value = nil
			getResp.Version = uuid.UUID{}
			continue

		case errors.As(err, &conflictErr):
			if conflictErr.ExpectedVersion == newVersion {
				return nil
			}

			return UpdateValue(c, key, updateFn)
		}
	}
}
