package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
)

func (i *keyCloakIAM) resolve(
	ctx context.Context,
	username string,
) (userID string, err error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return "", err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{},
	)
	if err != nil {
		return "", errors.Wrap(err, "keycloak")
	}
	for _, user := range users {
		if gocloak.PString(user.Username) == username {
			return gocloak.PString(user.ID), nil
		}
	}
	return "", iam.ErrUserNotFound
}
