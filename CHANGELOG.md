# CHANGELOG

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](http://semver.org/).
Which is based on [Keep A Changelog](http://keepachangelog.com/)

## Unreleased

### Fixed

- Bug in select item

## v1.2.0 - 2019-09-01

### Added

- Test build with travis
- Add entry's fields Create and LastUpdate
- Add timeout to close wallet for the security

### Changed

- Export in a file
- Use RandomString function for wallet's salt
- Search is case insensite
- change space to tab for goreport

## v1.1.0 - 2019-07-23

### Added

- Use go module to get this software
- Generate random password
- Print the expiration time of TOTP code
- Export a wallet in json
- Import entries from a json file
- Create the wallets directory

### Changed

- Prefix error message with ERROR 
- Fix new line with clear input
- Replace sha1 to sha512 in pbkdf2.Key function
- Replace default config directory
- Fix bug in the wallet choice

## v1.0.0 - 2019-07-12

### Added

- Save the wallet in AES-256 encrypted file
- Search entries with a pattern and/or by group
- Copy login, password and OTP code in clipboard
- Manage multiple wallets
