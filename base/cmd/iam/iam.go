package iam

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thavlik/bjjvb/base/pkg/base"
	"github.com/thavlik/bjjvb/base/pkg/iam"
	cognito_iam "github.com/thavlik/bjjvb/base/pkg/iam/cognito"
	keycloak_iam "github.com/thavlik/bjjvb/base/pkg/iam/keycloak"
)

var defaultTimeout = 10 * time.Second

var iamCmd = &cobra.Command{
	Use: "iam",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("please choose a subcommand")
	},
}

func AddIAMSubCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(iamCmd)
}

func InitIAM(o *base.IAMOptions) iam.IAM {
	switch o.Driver {
	case base.CognitoDriver:
		return cognito_iam.NewCognitoIAM(
			o.Cognito.AllowTokenUseBeforeIssue,
			base.Log,
		)
	case base.KeyCloakDriver:
		return keycloak_iam.NewKeyCloakIAM(
			base.ConnectKeyCloak(&o.KeyCloak),
			base.Log,
		)
	default:
		panic(fmt.Errorf("unrecognized iam driver '%s'", o.Driver))
	}
}
