package cmd

import (
	"context"
	"errors"

	"code.cloudfoundry.org/uaa-cli/cli"
	"code.cloudfoundry.org/uaa-cli/config"
	"code.cloudfoundry.org/uaa-cli/help"
	"github.com/cloudfoundry-community/go-uaa"
	"github.com/spf13/cobra"
)

func GetPasswordTokenValidations(cfg config.Config, args []string, username, password string) error {
	if err := cli.EnsureTargetInConfig(cfg); err != nil {
		return err
	}
	if len(args) < 1 {
		return cli.MissingArgumentError("client_id")
	}
	if password == "" {
		return cli.MissingArgumentError("password")
	}
	if username == "" {
		return cli.MissingArgumentError("username")
	}
	return validateTokenFormatError(tokenFormat)
}

func GetPasswordTokenCmd(cfg config.Config, clientId, clientSecret, username, password, tokenFormat string) error {
	requestedType := uaa.OpaqueToken
	if tokenFormat == uaa.JSONWebToken.String() {
		requestedType = uaa.JSONWebToken
	}

	api, err := uaa.New(
		cfg.GetActiveTarget().BaseUrl,
		uaa.WithPasswordCredentials(
			clientId,
			clientSecret,
			username,
			password,
			requestedType,
		),
		uaa.WithZoneID(cfg.ZoneSubdomain),
		uaa.WithSkipSSLValidation(cfg.GetActiveTarget().SkipSSLValidation),
		uaa.WithVerbosity(verbose),
	)
	if err != nil {
		return errors.New("An error occurred while fetching token.")
	}

	token, err := api.Token(context.Background())
	if err != nil {
		log.Info("Unable to retrieve token")
		return err
	}

	activeContext := cfg.GetActiveContext()
	activeContext.ClientId = clientId
	activeContext.GrantType = config.PASSWORD
	activeContext.Username = username

	activeContext.Token = *token
	cfg.AddContext(activeContext)
	config.WriteConfig(cfg)
	log.Info("Access token successfully fetched and added to context.")
	return nil
}

var getPasswordToken = &cobra.Command{
	Use:   "get-password-token CLIENT_ID -s CLIENT_SECRET -u USERNAME -p PASSWORD",
	Short: "Obtain an access token using the password grant type",
	Long:  help.PasswordGrant(),
	PreRun: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		cli.NotifyValidationErrors(GetPasswordTokenValidations(cfg, args, username, password), cmd, log)
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := GetSavedConfig()
		clientId := args[0]
		cli.NotifyErrorsWithRetry(GetPasswordTokenCmd(cfg, clientId, clientSecret, username, password, tokenFormat), log, GetSavedConfig())
	},
}

func init() {
	RootCmd.AddCommand(getPasswordToken)
	getPasswordToken.Annotations = make(map[string]string)
	getPasswordToken.Annotations[TOKEN_CATEGORY] = "true"
	getPasswordToken.Flags().StringVarP(&clientSecret, "client_secret", "s", "", "client secret")
	getPasswordToken.Flags().StringVarP(&username, "username", "u", "", "username")
	getPasswordToken.Flags().StringVarP(&password, "password", "p", "", "user password")
	getPasswordToken.Flags().StringVarP(&tokenFormat, "format", "", "jwt", "available formats include "+availableFormatsStr())
}
