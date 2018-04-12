package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Shopify/ejson"
	log "github.com/sirupsen/logrus"
)

// looks up current key pair or generates new one based on --update-public-key PUBLIC_KEY or --new-key-pair switches
func ejsonGetKeyPairs() []string {
	var err error

	var byters []byte

	var publicKeys []string
	var privateKeys []string

	switch {
	case args["--new-key-pair"].(bool):
		// generate new key pair
		publicKeys = append(publicKeys, "")
		privateKeys = append(privateKeys, "")

		publicKeys[0], privateKeys[0], err = ejson.GenerateKeypair()
		die(err)

		var dir string
		dir, err = os.Getwd()
		die(err)

		log.WithFields(log.Fields{
			"ejson_public_key": publicKeys[0],
			"directory":        dir,
		}).Debug("generated new ejson key pair")
	case args["--grab"].(bool):
		// grab ejson key pairs
		for i, doc := range docs {
			byters, err = ioutil.ReadFile(doc)
			die(err)

			var j map[string]interface{}
			err = json.Unmarshal(byters, &j)
			die(err)

			publicKeys = append(publicKeys, j["_public_key"].(string))
			privateKeys = append(privateKeys, ejsonLookupKeyPair(publicKeys[i]))

			log.WithFields(log.Fields{
				"ejson_public_key": publicKeys[i],
				"document":         doc,
			}).Debug("grabbed ejson key pair")
		}
	case args["--update-public-key"].(bool):
		// look up key pair
		publicKeys = append(publicKeys, args["<ejson_public_key>"].(string))
		privateKeys = append(privateKeys, ejsonLookupKeyPair(publicKeys[0]))

		log.WithFields(log.Fields{
			"ejson_public_key": publicKeys[0],
			"directory":        os.Getenv("EJSON_KEYDIR"),
		}).Debug("looked up ejson key pair")
	default:
		// no op
		os.Exit(0)
	}

	// put a copy of key pair in CWD
	for i := range publicKeys {
		if _, err = os.Stat(publicKeys[i] + ".gpg"); err == nil {
			if !args["--force"].(bool) {
				continue
			}
		}

		err = ioutil.WriteFile(
			publicKeys[i],
			[]byte(privateKeys[i]),
			0666,
		)
		die(err)

		dir, err := os.Getwd()
		die(err)

		log.WithFields(log.Fields{
			"ejson_public_key": publicKeys[i],
			"directory":        dir,
		}).Debug("upserted public key")

		// put a copy of key pair in CWD in EJSON_KEYDIR
		err = ioutil.WriteFile(
			os.Getenv("EJSON_KEYDIR")+"/"+publicKeys[i],
			[]byte(privateKeys[i]),
			0600,
		)
		die(err)

		log.WithFields(log.Fields{
			"ejson_public_key": publicKeys[i],
			"directory":        os.Getenv("EJSON_KEYDIR"),
		}).Debug("upserted public key")
	}

	return publicKeys
}
