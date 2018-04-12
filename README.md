# chest

## Motivation

[ejson](https://github.com/Shopify/ejson) is a brilliant command-line utility. But generating new keys, or grabbing existing key pairs, can be burdensome. `chest` attempts to ease the pain.

## WARNING

`admin-public-keys.asc` and `all-public-keys.asc` are copies of the author's PGP key! Replace them with your own. (For this utility to be useful, you need to replace each of these public keys with a group of recipients' keys.)

## Usage

```
  chest (-h | --version)
  chest -s
  chest [-df] -n <ejson_document> ...
  chest [-df] -g <ejson_document> ...
  chest [-df] -k <ejson_public_key> <ejson_document> ...
```

## Description

### To Hold Precious Secret Bits

```
  * Generates new EJSON key pair, ${EJSON_KEYDIR}/_PUBLIC_KEY; and
    adds EJSON public key to documents, ["_public_key"];
  * Grabs existing EJSON public keys from documents, ["_public_key"];
    puts plaintext pairs in ${EJSON_KEYDIR}/_PUBLIC_KEY; and
    puts cyphertext pairs in CWD;
  * Updates EJSON public key in documents, ["_public_key"];
  * Generates sample documents, --samples.
```

## Arguments

```
  <ejson_document>    optional EJSON document
  <ejson_public_key>  existing public-key, needs -k and <ejson_document>
```

## Options

```
  -d, --debug
  -f, --force              use the force to exec gpg decrypt PUBLIC_KEY.gpg
  -g, --grab               grab public keys from ejson documents
  -h, --help
  -k, --update-public-key  use existing ejson key pair to encrypt documents
  -n, --new-key-pair       generate new ejson key pair to use and gpg-encrypt them
  -s, --samples            generate sample secret yaml documents,
                           sample.json and samples.yaml
  --version
```

## Installation

```
  go get -u -v -ldflags="-s -w" github.com/tanakapayam/chest
```

## DOCKER

### BUILD

```
make docker-build
```

### PULL

```
docker pull tanakapayam/chest
```

### RUN

```
docker run --tty --volume ${HOME}/.ejson:/ejson --volume ${PWD}:/chest tanakapayam/chest:latest --help
```
