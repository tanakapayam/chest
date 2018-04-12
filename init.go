package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docopt/docopt-go"
	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
)

func init() {
	var err error

	log.SetLevel(log.WarnLevel)
	log.SetFormatter(
		&log.TextFormatter{
			FullTimestamp: true,
		},
	)

	usage := `Usage:
  ` + prog + ` (-h | --version)
  ` + prog + ` -s
  ` + prog + ` [-df] [--admin] -g <ejson_document> ...
  ` + prog + ` [-df] [--admin] -n <ejson_document> ...
  ` + prog + ` [-df] [--admin] -k <ejson_public_key> <ejson_document> ...

Description
  To Hold Precious Secret Bits
  * Generates new EJSON key pair, ${EJSON_KEYDIR}/_PUBLIC_KEY; and
    adds EJSON public key to documents, ["_public_key"];
  * Grabs existing EJSON public keys from documents, ["_public_key"];
    puts plaintext pairs in ${EJSON_KEYDIR}/_PUBLIC_KEY; and
    puts cyphertext pairs in CWD;
  * Updates EJSON public key in documents, ["_public_key"];
  * Generates sample documents, --samples.

Arguments:
  <ejson_document>    optional EJSON document
  <ejson_public_key>  existing public-key, needs -k and <ejson_document>

Options:
  --admin                  gpg-encrypt ejson key pair using admin keys only
  -d, --debug
  -f, --force              use the force to exec gpg decrypt PUBLIC_KEY.gpg
  -g, --grab               grab public keys from ejson documents
  -h, --help
  -k, --update-public-key  use existing ejson key pair to encrypt documents
  -n, --new-key-pair       generate new ejson key pair to use and gpg-encrypt them
  -s, --samples            generate sample secret yaml documents,
                           sample.json and samples.yaml
  --version

Installation
  go get -u -v -ldflags="-s -w" github.com/tanakapayam/` + prog

	// http://docopt.org/
	parser := &docopt.Parser{
		HelpHandler: func(err error, usage string) {
			var re = regexp.MustCompile(`((?:^|\n)\S+.*\n)`)
			usage = re.ReplaceAllString(usage, bold("$1"))

			if err != nil {
				_, err = fmt.Fprintln(os.Stderr, usage)
				die(err)
				os.Exit(1)
			} else {
				fmt.Println(usage)
				os.Exit(0)
			}
		},
	}

	args, err = parser.ParseArgs(
		usage,
		os.Args[1:],
		version,
	)
	die(err)

	if args["--debug"].(bool) {
		log.SetLevel(log.DebugLevel)
	}

	log.WithFields(log.Fields{
		"args": args,
	}).Debug("positional argument")

	if os.Getenv("EJSON_KEYDIR") == "" {
		log.Fatal("not set: EJSON_KEYDIR")
	} else {
		err = os.MkdirAll(
			os.Getenv("EJSON_KEYDIR"),
			0700,
		)
		die(err)
	}

	// create samples, if required
	samplesCreate()

	// glob
	for _, arg := range args["<ejson_document>"].([]string) {
		doc, err := filepath.Glob(arg)
		die(err)

		log.WithFields(log.Fields{
			"document": doc,
		}).Debug("parsed positional argument")

		if len(doc) > 0 {
			docs = append(docs, doc...)
		}
	}
	docs = uniq(docs)

	// convert yaml documents to ejson
	for i := range docs {
		docs[i] = yq(docs[i])
	}

	log.WithFields(log.Fields{
		"documents": docs,
	}).Debug("globbed")

	// read in pgp public keys
	b, err := ioutil.ReadFile("admin-public-keys.asc")
	die(err)
	adminPublicKeysAsc = string(b)

	b, err = ioutil.ReadFile("all-public-keys.asc")
	die(err)
	allPublicKeysAsc = string(b)
}

const (
	publicKeyPlaceHolder = "_PUBLIC_KEY"
)

var (
	args    docopt.Opts
	bold    = color.New(color.Bold).SprintFunc()
	docs    []string
	prog    = os.Args[0]
	version = "1.0.0"

	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgHiRed).SprintFunc()

	sampleJSON = `{
  "_public_key": "_PUBLIC_KEY",
  "_YOUR_KEY_1": "_YOUR_SINGLELINE_VALUE_1",
  "_YOUR_KEY_2": "_YOUR_SINGLELINE_VALUE_2"
}`
	samplesYAML = `---
# Generic Secrets
# - https://github.com/Shopify/ejson#format
# - https://stackoverflow.com/questions/3790454/in-yaml-how-do-i-break-a-string-over-multiple-lines/21699210#21699210
_public_key: "_PUBLIC_KEY"
_YOUR_KEY_1: "_YOUR_SINGLELINE_VALUE"
_YOUR_KEY_2: |-
  _YOUR_MULTILINE_VALUE
  _YOUR_MULTILINE_VALUE

---
# K8s Secrets
# - https://github.com/Shopify/ejson#format
# - https://github.com/Shopify/kubernetes-deploy#deploying-kubernetes-secrets-from-ejson
# - https://stackoverflow.com/questions/3790454/in-yaml-how-do-i-break-a-string-over-multiple-lines/21699210#21699210
# - kubectl create secret generic ejson-keys --from-literal="$(jq -r ._public_key < "$_SECRET_EJSON")=$(gpg --decrypt "$(jq -r ._public_key < "$_SECRET_EJSON").gpg")" --context=_CONTEXT --namespace=_NAMESPACE
_public_key: "_PUBLIC_KEY"
kubernetes_secrets:
  catphotoscom:
    _type: kubernetes.io/tls
    data:
      tls.crt: |-
        _YOUR_MULTILINE_VALUE
        _YOUR_MULTILINE_VALUE
      tls.key: |-
        _YOUR_MULTILINE_VALUE
        _YOUR_MULTILINE_VALUE
  monitoring-token:
    _type: Opaque
    data:
      _YOUR_KEY: |-
        _YOUR_SINGLE_OR_MULTILINE_VALUE

`

	adminPublicKeysAsc = ""
	allPublicKeysAsc   = ""
)
