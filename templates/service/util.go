package main

import (
	"encoding/hex"
	"errors"

	"gopkg.in/urfave/cli.v1"
)

// parseServicekey fetches the wercker-service key from the urfave cli context,
// and decodes it as a hex string to byte array. Returns an error if it is not
// present, or if it is not a valid hex encoded string.
func parseServiceKey(c *cli.Context) ([]byte, error) {
	k := c.String("service-key")
	if k == "" {
		return nil, errors.New("no wercker service key is supplied")
	}

	k, err := hex.DecodeString(k)
	if err != nil {
		return nil, errors.Wrap(err, "invalid hex encoded wercker service key")
	}

	return k, nil
}
