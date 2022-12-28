package keycloak

import (
	"context"
	"net/http"

	"github.com/Nerzal/gocloak/v12"
	"github.com/pkg/errors"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	"go.uber.org/zap"
)

func (i *keyCloakIAM) Login(ctx context.Context, username string, password string) (string, error) {
	token, err := i.kc.Login(
		ctx,
		i.kc.ClientID,
		i.kc.ClientSecret,
		i.kc.Realm,
		username,
		password,
	)
	if kcErr, ok := err.(*gocloak.APIError); ok && kcErr.Code == http.StatusUnauthorized {
		i.log.Error("failed keycloak login",
			zap.Error(err),
			zap.String("username", username))
		return "", iam.ErrInvalidCredentials
	} else if err != nil {
		return "", errors.Wrap(err, "keycloak")
	}
	//return strings.Join(strings.Split(token.AccessToken, ".")[:2], "."), nil
	return token.AccessToken, nil
}
