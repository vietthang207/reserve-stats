package client

import "github.com/urfave/cli"

const (
	coreURLFlag    = "core-url"
	coreSigningKey = "core-signing-key"
)

// NewCliFlags returns cli flags to configure a core client.
func NewCliFlags(prefix string) ([]cli.Flag) {
	return []cli.Flag{
		cli.StringFlag{
			Name:   coreURLFlag,
			Usage:  "Core API URL",
			EnvVar: prefix + "CORE_URL",
		},
		cli.StringFlag{
			Name:   coreSigningKey,
			Usage:  "Core API URL",
			EnvVar: prefix + "CORE_SIGNING_KEY",
		},
	}
}

// ValidateCliFlags validates core input flags.
func ValidateCliFlags(c *cli.Context) error {

}
