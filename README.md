# gpm: Go Passwords Manager

[![Version](https://img.shields.io/badge/latest_version-1.0.0-green.svg)](https://git.yaegashi.fr/nishiki/gpm/releases)
[![License](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://git.yaegashi.fr/nishiki/gpm/src/branch/master/LICENSE)

gpm is passwords manager write in go and use AES-256 to encrypt the wallets

## Features

- generate OTP code
- copy your login, password or otp in clipboard
- manage multiple wallets
- generate random password

## Install

### Build

Download the sources and build

```text
git clone https://git.yaegashi.fr/nishiki/gpm.git
cd gpm
go build -o bin/gpm src/*.go
```

Copy the binary in PATH:

```text
sudo cp bin/gpm /usr/local/bin/gpm
```

## How to use

```text
gpm -help
  -add
        add a new entry in the wallet
  -config string
        specify the config file
  -copy
        enter an copy mode for an entry
  -delete
        delete an entry
  -group string
        search the entries in this group
  -help
        print this help message
  -list
        list the entries in a wallet
  -pattern string
        search the entries with this pattern
  -update
        update an entry
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
