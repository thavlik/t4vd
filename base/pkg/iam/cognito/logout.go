package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
)

func (i *cognitoIAM) Logout(
	ctx context.Context,
	token string,
) error {
	hdr, err := retrieveAuthHeader(token)
	if err != nil {
		return err
	}
	if _, err := i.cognito.RevokeToken(
		&cognitoidentityprovider.RevokeTokenInput{
			Token:        aws.String(hdr.RefreshToken),
			ClientId:     aws.String(i.clientID),
			ClientSecret: aws.String(i.clientSecret),
		},
	); err != nil {
		return errors.Wrap(err, "cognito")
	}
	return nil
}
