package main

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func yq(doc string) string {
	var err error

	if !strings.HasSuffix(doc, ".yaml") && !strings.HasSuffix(doc, ".yml") {
		return doc
	}

	if !args["--force"].(bool) {
		log.WithFields(log.Fields{
			"document": doc,
			"format":   "yaml",
		}).Fatal("not json; use the --force")
	}

	jsonCmd := exec.Command("yq", "-S", ".", doc)

	var jsonOut bytes.Buffer
	jsonCmd.Stdout = &jsonOut

	err = jsonCmd.Run()
	die(err)

	// filename: *.yaml -> *.ejson
	var re = regexp.MustCompile(`(.yaml|.yml)$`)
	newDoc := re.ReplaceAllString(doc, ".ejson")

	err = ioutil.WriteFile(
		newDoc,
		jsonOut.Bytes(),
		0666,
	)
	die(err)

	log.WithFields(log.Fields{
		"old": doc,
		"new": newDoc,
	}).Warn("converted from yaml to ejson")

	return newDoc
}
