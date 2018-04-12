package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func replaceEjsonPublicKey(doc string, public string) {
	var err error

	var byters []byte
	byters, err = ioutil.ReadFile(doc)
	die(err)

	var j map[string]interface{}
	err = json.Unmarshal(byters, &j)
	die(err)

	j["_public_key"] = public

	byters, err = json.Marshal(j)
	die(err)

	replaced := &bytes.Buffer{}
	err = json.Indent(replaced, byters, "", "  ")
	die(err)

	err = ioutil.WriteFile(
		doc,
		replaced.Bytes(),
		0666,
	)
	die(err)
}
