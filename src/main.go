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

package main

import(
  "flag"
  "os"
)

// Options
var(
  ADD     = flag.Bool("add", false, "add a new entry in the wallet")
  UPDATE  = flag.Bool("update", false, "update an entry")
  DELETE  = flag.Bool("delete", false, "delete an entry")
  LIST    = flag.Bool("list", false, "list the entries in a wallet")
  COPY    = flag.Bool("copy", false, "enter an copy mode for an entry")
  CONFIG  = flag.String("config", "", "specify the config file")
  GROUP   = flag.String("group", "", "search the entries in this group ")
  PATTERN = flag.String("pattern", "", "search the entries with this pattern")
  WALLET  = flag.String("wallet", "", "specify the wallet")
  HELP    = flag.Bool("help", false, "print this help message")
)

func init() {
  flag.Parse()

  if *HELP {
    flag.PrintDefaults()
    os.Exit(1)
  }
}

func main() {
  c := Cli{}
  c.Init()
  c.Run()
}
