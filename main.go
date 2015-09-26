package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "simplevault"
	app.Usage = "Simple vault app for storing things on S3"
	app.Before = setUp
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "aws-access-key",
			Usage:  "AWS Access Key",
			EnvVar: "SIMPLEVAULT_AWS_ACCESS_KEY_ID,AWS_ACCESS_KEY_ID",
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			Usage:  "AWS Secret Key",
			EnvVar: "SIMPLEVAULT_AWS_SECRET_ACCESS_KEY,AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "master-password",
			Usage:  "The Master Password",
			EnvVar: "SIMPLEVAULT_AWS_SECRET_ACCESS_KEY,AWS_SECRET_ACCESS_KEY",
		},
		cli.StringFlag{
			Name:   "bucket",
			Usage:  "S3 bucket to store things in",
			EnvVar: "SIMPLEVAULT_BUCKET",
		},
		cli.StringFlag{
			Name:   "bucket-prefix",
			Usage:  "Prefix to store everything under inside the bucket",
			EnvVar: "SIMPLEVAULT_BUCKET_PREFIX",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "set",
			Aliases: []string{"s"},
			Usage:   "stash some data in the vault",
			Action:  cliGet,
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get some data from the vault",
			Action:  cliGet,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("app error %v\n", err)
	}
}

func cliSet(c *cli.Context) {

}

func cliGet(c *cli.Context) {

}

func setUp(c *cli.Context) error {
	if c.String("aws-access-key") == "" ||
		c.String("aws-secret-key") == "" ||
		c.String("bucket") == "" {
		return errors.New("The aws-access-key, aws-secret-key and bucket parameters are mandatory")
	}
	return nil
}
