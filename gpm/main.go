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
  "flag"
  "os"
)

// Options
var(
  ADD     = flag.Bool("add", false, "add a new entry in the wallet")
  UPDATE  = flag.Bool("update", false, "update an entry")
  DELETE  = flag.Bool("delete", false, "delete an entry")
  LIST    = flag.Bool("list", false, "list the entries in a wallet")
  LENGTH  = flag.Int("length", 16, "specify the password length")
  COPY    = flag.Bool("copy", false, "enter an copy mode for an entry")
  CONFIG  = flag.String("config", "", "specify the config file")
  GROUP   = flag.String("group", "", "search the entries in this group ")
  WALLET  = flag.String("wallet", "", "specify the wallet")
  PATTERN = flag.String("pattern", "", "search the entries with this pattern")
  RANDOM  = flag.Bool("random", false, "generate a random password for a new entry or an update")
  PASSWD  = flag.Bool("password", false, "generate and print a random password")
  DIGIT   = flag.Bool("digit", false, "use digit to generate a random password")
  LETTER  = flag.Bool("letter", false, "use letter to generate a random password")
  SPECIAL = flag.Bool("special", false, "use special chars to generate a random password")
  HELP    = flag.Bool("help", false, "print this help message")
)

// Run the cli interface
func Run() {
  var cli Cli
  cli.Config.Load(*CONFIG)

  flag.Parse()
  if *HELP {
    flag.PrintDefaults()
    os.Exit(1)
  } else if *PASSWD {
    fmt.Println(RandomString(*LENGTH, *LETTER, *DIGIT, *SPECIAL))
  } else if *LIST {
    cli.listEntry()
  } else if *COPY {
    cli.copyEntry()
  } else if *ADD {
    cli.addEntry()
  } else if *UPDATE {
    cli.updateEntry()
  } else if *DELETE {
    cli.deleteEntry()
  }
}
