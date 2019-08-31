// Copyright 2019 Adrien Waksberg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gpm

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pquerna/otp/totp"
)

// Entry struct have the password informations
type Entry struct {
	Name       string
	ID         string
	URI        string
	User       string
	Password   string
	OTP        string
	Group      string
	Comment    string
	Create     int64
	LastUpdate int64
}

// Verify if the item have'nt error
func (e *Entry) Verify() error {
	if e.ID == "" {
		return fmt.Errorf("you must generate an ID")
	}

	if e.Name == "" {
		return fmt.Errorf("you must define a name")
	}

	uri, _ := url.Parse(e.URI)
	if e.URI != "" && uri.Host == "" {
		return fmt.Errorf("the uri isn't a valid uri")
	}

	return nil
}

// GenerateID create a new id for the entry
func (e *Entry) GenerateID() {
	e.ID = fmt.Sprintf("%d", time.Now().UnixNano())
}

// OTPCode generate an OTP Code
func (e *Entry) OTPCode() (string, int64, error) {
	code, err := totp.GenerateCode(e.OTP, time.Now())
	time := 30 - (time.Now().Unix() % 30)
	if err != nil {
		return "", 0, err
	}

	return code, time, nil
}
