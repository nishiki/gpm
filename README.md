# gpm: Go Passwords Manager

[![Version](https://img.shields.io/badge/latest_version-2.0.0-green.svg)](https://git.yaegashi.fr/nishiki/gpm/releases)
[![Build Status](https://travis-ci.org/nishiki/gpm.svg?branch=master)](https://travis-ci.org/nishiki/gpm)
[![GoReport](https://goreportcard.com/badge/git.yaegashi.fr/nishiki/gpm)](https://goreportcard.com/report/git.yaegashi.fr/nishiki/gpm)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://git.yaegashi.fr/nishiki/gpm/src/branch/master/LICENSE)

gpm is passwords manager write in go and use AES-256 to encrypt the wallets

## Features

- generate TOTP code
- copy your login, password or otp in clipboard
- manage multiple wallets
- generate random password

## Install

- Install [golang](https://golang.org/doc/install)
- Add `~/go/bin` in your `PATH`
- Download and build

```text
go get git.yaegashi.fr/nishiki/gpm/cmd/gpm
```

## How to use

### First launch

- Run `gpm`
- Enter the passphrase to encrypt your new wallet
- Press `n` to create your first entry and follow the instructions

### All options

```text
  -config string
    	specify the config file
  -digit
    	use digit to generate a random password
  -export string
    	json file path to export a wallet
  -help
    	print this help message
  -import string
    	json file path to import entries
  -length int
    	specify the password length (default 16)
  -letter
    	use letter to generate a random password
  -password
    	generate and print a random password
  -special
    	use special chars to generate a random password
  -wallet string
    	specify the wallet
```

## License

```text
Copyright (c) 2019 Adrien Waksberg

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
