package main

import (
	"sync"

	"github.com/Shopify/ejson"
	log "github.com/sirupsen/logrus"
)

func updateEjsonDoc(doc string, public string, wg *sync.WaitGroup) {
	var err error

	ejsonDecrypt(doc)

	replaceEjsonPublicKey(doc, public)

	_, err = ejson.EncryptFileInPlace(doc)
	die(err)

	log.Debug("processed ", doc)
	wg.Done()
}
