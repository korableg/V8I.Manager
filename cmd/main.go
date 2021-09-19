package main

import (
	"log"
	"os"

	"github.com/korableg/V8I.Manager/internal/worker"
	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:    "V8I.Manager",
		Usage:   "Создает файл *.v8i со списком баз 1С на основании *.lst файла кластера серверов",
		Version: "1.1.0",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "lst",
				Usage:    "Путь до 1CV8Clst.lst файла",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:     "v8i",
				Usage:    "Путь до iBases.v8i файла",
				Required: true,
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "Dmitry Titov",
				Email: "dim@titovcode.com",
			},
		},
		Copyright: "© Dmitry Titov | titovcode.com",
		Action:    cliMain,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func cliMain(c *cli.Context) error {

	lst := c.StringSlice("lst")
	v8i := c.StringSlice("v8i")

	return worker.LstToV8i(lst, v8i)

}
