package keycloak

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/iam"
)

func (i *keyCloakIAM) GetUser(ctx context.Context, username string) (*iam.User, error) {
	accessToken, err := i.refreshAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	users, err := i.kc.GetUsers(
		ctx,
		accessToken,
		i.kc.Realm,
		gocloak.GetUsersParams{
			Username: gocloak.StringP(username),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "keycloak")
	}
	for _, user := range users {
		if gocloak.PString(user.Username) == username {
			return &iam.User{
				ID:        gocloak.PString(user.ID),
				Username:  gocloak.PString(user.Username),
				Email:     gocloak.PString(user.Email),
				FirstName: gocloak.PString(user.FirstName),
				LastName:  gocloak.PString(user.FirstName),
				Enabled:   gocloak.PBool(user.Enabled),
			}, nil
		}
	}
	return nil, iam.ErrUserNotFound
}
