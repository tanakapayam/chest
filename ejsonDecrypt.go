package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/Shopify/ejson"
	log "github.com/sirupsen/logrus"
)

// partially-encrypted-doc oddity:
// - cannot be fully decrypted
// - can be fully encrypted
func ejsonDecrypt(doc string) {
	var err error

	var byters []byte
	byters, err = ioutil.ReadFile(doc)
	die(err)

	var j map[string]interface{}
	err = json.Unmarshal(byters, &j)
	die(err)

	if _, ok := j["_public_key"]; !ok {
		j["_public_key"] = publicKeyPlaceHolder
	}

	if j["_public_key"].(string) == publicKeyPlaceHolder || j["_public_key"].(string) == "" {
		// fresh document; nothing to decrypt
		return
	} else if _, err := os.Stat(os.Getenv("EJSON_KEYDIR") + "/" + j["_public_key"].(string)); err == nil {
		// found key pair in $EJSON_KEYDIR
		_, err = ejson.EncryptFileInPlace(doc)
		die(err)

		byters, err = ejson.DecryptFile(
			doc,
			os.Getenv("EJSON_KEYDIR"),
			"",
		)
		die(err)

		decrypted := &bytes.Buffer{}
		err = json.Indent(decrypted, byters, "", "  ")
		die(err)

		err = ioutil.WriteFile(
			doc,
			decrypted.Bytes(),
			0666,
		)
		die(err)
	} else if _, err = os.Stat(j["_public_key"].(string) + ".gpg"); err == nil {
		// key pair needs to be decrypted
		log.WithFields(log.Fields{
			"ejson_public_key": j["_public_key"].(string),
		}).Fatal("missing in $EJSON_KEYDIR/; decrypt gpg ciphertext")
	} else {
		// no key pair found
		log.WithFields(log.Fields{
			"ejson_public_key": j["_public_key"].(string),
		}).Fatal("missing key pair; what now?")
	}
}
