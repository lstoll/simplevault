# simple vault

[![Build Status](https://travis-ci.org/lstoll/simplevault.svg)](https://travis-ci.org/lstoll/simplevault)

:warning: **This is not really secured or vetted. Use it for things that aren't important** :warning:

This is a really simple vault tool that will store blobs in S3 under a given key. It was designed to simply share docker-macchine certs and .env files with places like Travis CI or deployment systems

Blobs are encryped and have two keys. One is unique to this blob to be used to decrypt specific items on your target clients. The other is a main key, that can decrypt anything and save items. All the configuration is via environment variables.

## Usage

The following environment variables are expected to be set

* `SIMPLEVAULT_AWS_ACCESS_KEY_ID` or `AWS_SECRET_ACCESS_KEY`
* `SIMPLEVAULT_AWS_SECRET_ACCESS_KEY` or `AWS_SECRET_ACCESS_KEY`
* `SIMPLEVAULT_PASSWORD_<PATH_TO_KEY_CAPS_UNDERSCORE>` or `SIMPLEVAULT_PASSWORD` - password to unlock items with. This can be a comma seperated list, if you want to store multiple keys.

`simplevault set <keypath> [filename]`. Will set the content of keypath to content of file, or from incoming piped stdin.

`simplevault get <keypath> [filename]`. Will write the content of keypath to a file, or to stdout if no file specified

`simplevault delete <keypath>`. Will delete the item at keypath. No password is required for this operation.

## Caveats
* No credential rolling. If you're compromised, delete and re-upload
* Can't choose non-master password. These are randomly generated
