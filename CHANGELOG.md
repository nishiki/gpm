# CHANGELOG

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/).
Which is based on [Keep A Changelog](http://keepachangelog.com/)

## Unreleased

### Added

- Use go module to get this software
- Generate random password
- Print the expiration time of TOTP code
- Export a wallet in json
- Import entries from a json file

## Changed

- Prefix error message with ERROR 
- Fix new line with clear input
- Replace sha1 to sha512 in pbkdf2.Key function
- Replace default config directory

## v1.0.0 - 2019-07-12

### Added

- Save the wallet in AES-256 encrypted file
- Search entries with a pattern and/or by group
- Copy login, password and OTP code in clipboard
- Manage multiple wallets
