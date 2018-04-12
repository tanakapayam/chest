package main

func main() {
	// get or generate ejson key pair
	// put them in EJSON_KEYDIR and CWD
	publicKeys := ejsonGetKeyPairs()

	// encrypt the ejson key pair in CWD
	gpgEncryptEjsonKeyPairs(publicKeys)

	// take the first public key
	// process documents
	processEjsonDocs(publicKeys[0])
}
