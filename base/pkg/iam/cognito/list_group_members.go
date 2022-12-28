package cognito

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/pkg/errors"
	"github.com/thavlik/t4vd/base/pkg/iam"
)

func (i *cognitoIAM) ListGroupMembers(
	ctx context.Context,
	groupID string,
) ([]*iam.User, error) {
	result, err := i.cognito.ListUsersInGroup(
		&cognitoidentityprovider.ListUsersInGroupInput{
			UserPoolId: aws.String(i.userPoolID),
			GroupName:  aws.String(groupID),
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "cognito")
	}
	n := len(result.Users)
	users := make([]*iam.User, n)
	for i, user := range result.Users {
		users[i] = &iam.User{
			ID:       aws.StringValue(user.Username),
			Username: aws.StringValue(user.Username),
		}
		applyUserAttributes(users[i], user.Attributes)
	}
	return users, nil
}
