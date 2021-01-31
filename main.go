package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

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
	payload, err := decodePart(parts[1])
	if err != nil {
		return err
	}

	var headerParams map[string]interface{}
	if err := json.Unmarshal(header, &headerParams); err != nil {
		return err
	}
	fmtHeader, err := json.MarshalIndent(headerParams, "", " ")
	if err != nil {
		return err
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return err
	}
	fmtPayload, err := json.MarshalIndent(claims, "", " ")
	if err != nil {
		return err
	}

	var expHuman, nbfHuman, iatHuman time.Time
	if v, ok := claims["exp"]; ok {
		if vv, ok := v.(float64); ok {
			expHuman = time.Unix(int64(vv), 0)
		}
	}
	if v, ok := claims["nbf"]; ok {
		if vv, ok := v.(float64); ok {
			nbfHuman = time.Unix(int64(vv), 0)
		}
	}
	if v, ok := claims["iat"]; ok {
		if vv, ok := v.(float64); ok {
			iatHuman = time.Unix(int64(vv), 0)
		}
	}

	fmt.Println("\n=== Header ===")
	fmt.Println(string(fmtHeader))
	fmt.Println("\n===Payload===")
	fmt.Println(string(fmtPayload))
	fmt.Println("\n===TimeHuman===")
	if !isZero(iatHuman) {
		fmt.Println("iat:", iatHuman.String())
	}
	if !isZero(nbfHuman) {
		fmt.Println("nbf:", nbfHuman.String())
	}
	if !isZero(expHuman) {
		fmt.Println("exp:", expHuman.String())
	}
	return nil
}

func decodePart(part string) ([]byte, error) {
	for {
		if len(part)%4 == 0 {
			break
		}
		part = part + "="
	}

	part = strings.ReplaceAll(part, "-", "+")
	part = strings.ReplaceAll(part, "_", "/")
	return base64.StdEncoding.DecodeString(part)
}

func isZero(t time.Time) bool {
	return t == time.Time{}
}

func main() {
	app := &cli.App{
		Name:  "jwtd",
		Usage: "JWT decoding tool.",
		Commands: []*cli.Command{{
			Name:    "decode",
			Aliases: []string{"d"},
			Usage:   "Decodes a JWT in JWS format.",
			Action: func(c *cli.Context) error {
				return decodeJWT(c.Args().Get(0))
			},
		}},
		Copyright: "MIT",
		Action: func(c *cli.Context) error {
			return decodeJWT(c.Args().Get(0))
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println("error:", err)
	}
}
