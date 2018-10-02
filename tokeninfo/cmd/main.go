package main

import (
	"encoding/json"
	"log"
	"os"

	libapp "github.com/KyberNetwork/reserve-stats/lib/app"
	"github.com/KyberNetwork/reserve-stats/lib/core/client"
	"github.com/KyberNetwork/reserve-stats/tokeninfo"
	"github.com/urfave/cli"
)

const (
	nodeURLFlag         = "node"
	nodeURLDefaultValue = "https://mainnet.infura.io"
	outputFlag          = "output"
)

func main() {
	app := libapp.NewApp()
	app.Name = "token reserve fetcher"
	app.Usage = "fetching token reserve mapping information"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "reserve",
			Aliases: []string{"r"},
			Usage:   "report which reserves provides which token",
			Action:  reserve,
			Flags: append(client.NewCliFlags("TOKEN_INFO_"),
				cli.StringFlag{
					Name:  nodeURLFlag,
					Usage: "Ethereum node provider URL",
					Value: nodeURLDefaultValue,
				},
				cli.StringFlag{
					Name:  outputFlag,
					Usage: "output file location",
					Value: "./output.json",
				},
			),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func reserve(c *cli.Context) error {
	logger, err := libapp.NewLogger(c)
	if err != nil {
		return err
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	coreClient, err := client.NewClientFromContext(c)
	if err != nil {
		return err
	}

	tokens, err := coreClient.Tokens()
	if err != nil {
		return err
	}

	log.Println(tokens)

	f, err := tokeninfo.NewReserveCrawler(
		sugar,
		c.String(nodeURLFlag))
	if err != nil {
		return err
	}

	output, err := os.Create(c.String(outputFlag))
	if err != nil {
		return err
	}

	result, err := f.Fetch()
	if err != nil {
		return err
	}

	return json.NewDecoder(output).Decode(result)
}
