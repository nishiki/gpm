# gpm: Go Passwords Manager

[![Version](https://img.shields.io/badge/latest_version-1.2.1-green.svg)](https://git.yaegashi.fr/nishiki/gpm/releases)
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

### First steps

- Add new entry `gpm -add`

```text
Enter the passphrase to unlock the wallet: 
Enter the name: Test
Enter the group: MyGroup
Enter the URI: http://localhost
Enter the username: lastname
Enter the new password: 
Enter the OTP key: 
Enter a comment: My first entry
the entry has been added
```

- Search and copy `gpm -copy`

```text
Enter the passphrase to unlock the wallet: 

MyGroup

    | NAME |       URI        |   USER   | OTP |    COMMENT      
----+------+------------------+----------+-----+-----------------
  0 | Test | http://localhost | lastname | X   | My first entry  

select one action: p
select one action: l
select one action: q
```

### All options

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
  -digit
        use digit to generate a random password
  -export
        export a wallet in json format
  -group string
        search the entries in this group 
  -help
        print this help message
  -import string
        import entries from a json file
  -length int
        specify the password length (default 16)
  -letter
        use letter to generate a random password
  -list
        list the entries in a wallet
  -password
        generate and print a random password
  -pattern string
        search the entries with this pattern
  -random
        generate a random password for a new entry or an update
  -special
        use special chars to generate a random password
  -update
        update an entry
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
