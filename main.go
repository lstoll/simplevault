package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

var v Vault

func main() {
	app := cli.NewApp()
	app.Name = "simplevault"
	app.Usage = "Simple vault app for storing things on S3"
	app.Before = validate
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
			Name:   "password",
			Usage:  "Password to access with",
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
			Action:  cliSet,
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
	var data []byte
	var err error

	key := c.Args().Get(0)
	if key == "" {
		cliErrChk(errors.New("You need to provide the key to save under"))
	}
	filename := c.Args().Get(1)
	if filename == "" {
		// do we have piped data?
		if incomingPipe() {
			data, err = getStdin()
			cliErrChk(err)
		} else {
			cliErrChk(errors.New("You need to specifiy a source file or pipe data in"))
		}
	} else {
		data, err = getFile(filename)
		cliErrChk(err)
	}

	accessPass, err := v.PutItem(key, getPassword(c), data)
	cliErrChk(err)
	fmt.Printf("Item stored, its access key is %s", accessPass)
}

func cliGet(c *cli.Context) {
	key := c.Args().Get(1)
	if key == "" {
		cliErrChk(errors.New("You need to provide the key to read from"))
	}

	fmt.Println(key)

	filename := c.Args().Get(0)
	if filename == "" {
		cliErrChk(errors.New("You need to provide a filename or - for stdout"))
	}

	fmt.Println(filename)

	data, err := v.GetItem(key, getPassword(c))
	cliErrChk(err)
	err = putSource(filename, data)
	cliErrChk(err)
}

func getPassword(c *cli.Context) string {
	return c.String("password")
}

func keyToEnviron(key string) string {
	// do this later
	return "AAAAA"
}

func cliErrChk(err error) {
	if err != nil {
		fmt.Printf("Error!: %v\n", err)
		os.Exit(1)
	}
}

func incomingPipe() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) == 0
}

func getStdin() ([]byte, error) {
	return ioutil.ReadAll(os.Stdin)
}

func getFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func putSource(filename string, data []byte) error {
	if filename == "-" {
		fmt.Print(string(data))
	} else {
		return ioutil.WriteFile(filename, data, 0600)
	}
	return nil
}

func setUp(c *cli.Context) (Vault, error) {
	encryptor := NewEnc()
	fmt.Println(c.String("aws-access-key"))
	s3client := NewS3(c.String("aws-access-key"), c.String("aws-secret-key"), c.String("bucket"), c.String("bucket-prefix"))
	return NewVault(s3client, encryptor), nil
}

func validate(c *cli.Context) error {
	if c.String("aws-access-key") == "" ||
		c.String("aws-secret-key") == "" ||
		c.String("bucket") == "" {
		return errors.New("The aws-access-key, aws-secret-key and bucket parameters are mandatory")
	}
	v, _ = setUp(c)
	return nil
}
