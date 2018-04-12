package main

import log "github.com/sirupsen/logrus"

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
