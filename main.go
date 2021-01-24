package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "jwtd",
		Usage: "TODO",
		Action: func(c *cli.Context) error {
			return decodeJWT(c.Args().Get(0))
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}

func decodeJWT(jwt string) error {
	if len(jwt) == 0 {
		return errors.New("You need to specify a JWT.")
	}

	parts := strings.Split(jwt, ".")
	if len(parts) != 3 {
		return errors.New("Your JWT format is incorrect.")
	}

	header, err := decodePart(parts[0])
	if err != nil {
		return err
	}

	payload, err := decodePart(parts[2])
	if err != nil {
		return err
	}

	signature, err := decodePart(parts[2])
	if err != nil {
		return err
	}

	fmt.Println("Header:")
	fmt.Println(string(header))
	fmt.Println("Payload:")
	fmt.Println(string(payload))
	fmt.Println("Signature:")
	fmt.Println(string(signature))
	return nil
}

func decodePart(part string) ([]byte, error) {
	switch len(part) % 4 {
	case 2:
		part = part + "=="
	case 3:
		part = part + "="
	}

	return base64.StdEncoding.DecodeString(part)
}
