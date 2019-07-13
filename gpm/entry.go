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

import(
  "fmt"
  "time"
  "net/url"
)

// Entry struct have the password informations
type Entry struct {
  Name     string `yaml:"name"`
  ID       string `yaml:"id"`
  URI      string `yaml:"uri"`
  User     string `yaml:"login"`
  Password string `yaml:"password"`
  OTP      string `yaml:"otp"`
  Group    string `yaml:"group"`
  Comment  string `yaml:"comment"`
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
