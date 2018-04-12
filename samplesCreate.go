package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func samplesCreate() {
	var err error

	if !args["--samples"].(bool) {
		return
	}

	_, err = fmt.Fprintln(
		os.Stderr,
		bold("# SAMPLES")+`

* "jq" and "yq" are required for command-line JSON and YAML processing.
  * For Darwin:
    `+green("brew install jq")+`
    `+green("brew install python # brew install pypy")+`
    `+green("pip install yq      # pip_pypy install yq")+`

* sample.json and samples.yaml were created
    `+yellow("ls -l sample.json samples.yaml")+`
* Edit sample.json or samples.yaml
    `+green("vi sample.json")+`
    `+green("vi samples.yaml")+`
  * _PUBLIC_KEY will be generated after "`+prog+` -n" on document
  * The following are placeholders -- replace or remove them:
    * _SECRET_FILENAME
    * _YOUR_KEY, _YOUR_SINGLE_OR_MULTILINE_VALUE, _YOUR_SINGLELINE_VALUE, _YOUR_MULTILINE_VALUE
    * _CONTEXT, _NAMESPACE
* Replace YOUR-FILENAME with your secrets filename
    `+green("_SECRET_EJSON=\"YOUR-FILENAME.ejson\"")+`
* If sample.json was chosen, standardize it to EJSON
    `+green("jq . < sample.json > \"$_SECRET_EJSON\"")+`
* If samples.yaml was chosen, convert it to EJSON
    `+green("yq . < samples.yaml > \"$_SECRET_EJSON\"")+`
* Verify
    `+green("jq . < \"$_SECRET_EJSON\"")+`
* Encrypt with --new-key-pair
    `+green(prog+" -n \"$_SECRET_EJSON\"")+`
* Verify
    `+green("jq . < \"$_SECRET_EJSON\"")+`
    `+green("ls -l \"$EJSON_KEYDIR/$(jq -r ._public_key < \"$_SECRET_EJSON\")\" \"./$(jq -r ._public_key < \"$_SECRET_EJSON\").gpg\"")+`
* delete artifacts
    `+red("rm -f sample.json samples.yaml")+`
* If the secret is meant for use in a K8s cluster, save the EJSON key pair in the cluster:
    `+yellow(`kubectl --context=_CONTEXT --namespace=_NAMESPACE create secret generic ejson-keys --from-literal="$(jq -r ._public_key < "$_SECRET_EJSON")=$(gpg --decrypt "$(jq -r ._public_key < "$_SECRET_EJSON").gpg")"`),
	)
	die(err)

	err = ioutil.WriteFile(
		"sample.json",
		[]byte(sampleJSON),
		0666,
	)
	die(err)

	err = ioutil.WriteFile(
		"samples.yaml",
		[]byte(samplesYAML),
		0666,
	)
	die(err)

	os.Exit(0)
}
