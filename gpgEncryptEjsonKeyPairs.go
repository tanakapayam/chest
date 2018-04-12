package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
)

// gpg-encrypts ejson private key
// into ejson-public-key.gpg; removes plaintext ejson-public-key
func gpgEncryptEjsonKeyPairs(publicKeys []string) {
	var err error

	var byters []byte

	if len(publicKeys) == 0 {
		return
	}

	for _, public := range publicKeys {
		if _, err = os.Stat(public + ".gpg"); err == nil {
			if !args["--force"].(bool) {
				continue
			}
		}

		byters, err = ioutil.ReadFile(os.Getenv("EJSON_KEYDIR") + "/" + public)
		private := string(byters)
		die(err)

		entityList, err := openpgp.ReadArmoredKeyRing(
			strings.NewReader(allPublicKeysAsc),
		)
		die(err)

		if args["--admin"].(bool) {
			entityList, err = openpgp.ReadArmoredKeyRing(
				strings.NewReader(adminPublicKeysAsc),
			)
		}
		die(err)

		buf := new(bytes.Buffer)
		w, err := openpgp.Encrypt(
			buf,
			entityList,
			nil,
			nil,
			nil,
		)
		die(err)

		_, err = w.Write([]byte(private))
		die(err)

		err = w.Close()
		die(err)

		byters, err = ioutil.ReadAll(buf)
		die(err)

		err = ioutil.WriteFile(
			public+".gpg",
			byters,
			0666,
		)
		die(err)

		err = os.Remove(public)
		die(err)

		log.Debug("created ./" + public + ".gpg")
	}
}
