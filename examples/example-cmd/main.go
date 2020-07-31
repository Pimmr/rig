package main

import (
	"encoding/json"
	"net/url"
	"os"

	"github.com/Pimmr/rig"
	"github.com/Pimmr/rig/command"
)

type Configuration struct {
	Foo int
	Bar string
}

type ConfCmdA struct {
	Fizz string
	Buzz []string
}

type ConfCmdB struct {
	Hello int
	World []int
}

type ConfCmdC struct {
	Ahoy *url.URL
}

func main() {
	var config Configuration

	flags, err := rig.StructToFlags(&config)
	if err != nil {
		panic(err)
	}

	err = command.Parse(
		flags,
		command.New("cmda", func(remainder command.Args) error {
			var cmdaConfig ConfCmdA

			err = remainder.ParseStruct(&cmdaConfig)
			if err != nil {
				return err
			}

			pprint(config)
			pprint(cmdaConfig)

			return nil
		}),
		command.New("cmdb", func(remainder command.Args) error {
			var cmdbConfig ConfCmdB

			flags, err := rig.StructToFlags(&cmdbConfig)
			if err != nil {
				return err
			}

			return remainder.ParseCommands(flags, command.New("cmdc", func(remainder command.Args) error {
				var cmdcConfig ConfCmdC

				err := remainder.ParseStruct(&cmdcConfig)
				if err != nil {
					return err
				}

				pprint(config)
				pprint(cmdbConfig)
				pprint(cmdcConfig)

				return nil
			}))
		}),
	)
	if err != nil {
		panic(err)
	}
}

func pprint(i interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	err := enc.Encode(i)
	if err != nil {
		panic(err)
	}
}
