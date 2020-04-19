package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/interpreter"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "covid",
		Authors:   []*cli.Author{{Name: "DHL", Email: "dhl1402@gmail.com"}},
		HelpName:  "covid",
		Usage:     "an useless tool for managing an useless language source code",
		Version:   "0.0.1-alpha.2",
		UsageText: "covid example.covs",
		Action: func(c *cli.Context) error {
			fileName := c.Args().First()
			if fileName == "" {
				cli.ShowAppHelp(c)
				return nil
			}
			b, err := ioutil.ReadFile(fileName)
			if err != nil {
				return err
			}
			err = interpreter.Interpret(string(b), config.Config{
				Writer: os.Stdout,
			})
			return err
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
