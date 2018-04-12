package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ejsonLookupKeyPair(public string) string {
	var err error

	var byters []byte

	private := ""

	if public == "" || public == publicKeyPlaceHolder {
		log.WithFields(log.Fields{
			"ejson_public_key": public,
		}).Fatal("found placeholder; use --new-key-pair")
	} else if _, err = os.Stat(os.Getenv("EJSON_KEYDIR") + "/" + public); err == nil {
		byters, err = ioutil.ReadFile(os.Getenv("EJSON_KEYDIR") + "/" + public)
		die(err)

		private = string(byters)

		log.WithFields(log.Fields{
			"ejson_public_key": public,
			"directory":        os.Getenv("EJSON_KEYDIR"),
		}).Debug("looked up")
	} else if _, err = os.Stat(public + ".gpg"); err == nil {
		if args["--force"].(bool) {
			var gpgArgs = []string{
				"--decrypt",
				"--batch",
				"--yes",
				public + ".gpg",
			}

			gpgCmd := exec.Command("gpg", gpgArgs...)

			var gpgOut bytes.Buffer
			gpgCmd.Stdout = &gpgOut

			err = gpgCmd.Run()
			die(err)

			private = gpgOut.String()

			dir, err := os.Getwd()
			die(err)

			log.WithFields(log.Fields{
				"ejson_public_key": public,
				"directory":        dir,
			}).Debug("looked up")
		} else {
			log.WithFields(log.Fields{
				"ejson_public_key": public,
				"directory":        os.Getenv("EJSON_KEYDIR"),
			}).Fatal("missing key pair; to decrypt, use the --force")
		}
	} else {
		log.WithFields(log.Fields{
			"ejson_public_key": public,
		}).Fatal("missing key pair; i'm out of ideas")
	}

	return strings.TrimSpace(private)
}
