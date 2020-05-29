package args

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	ErrPathLstEmpty    = errors.New("Path to LST not filled in")
	ErrPathIBasesEmpty = errors.New("Path to iBases not filled in")

	pathLstRaw string
	pathiBasesRaw string
	instance     string

	fs *flag.FlagSet
)

func init() {

	fs = flag.NewFlagSet("creator", flag.ContinueOnError)
	fs.Usage = Usage

	fs.StringVar(&instance, "instance", "Agent iBases", "Instance name")
	fs.StringVar(&pathLstRaw, "lst", "", "Path to LST (comma separated)")
	fs.StringVar(&pathiBasesRaw, "ibases", "", "Path to iBases (comma separated)")

	i := 1
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		i = 2
	}

	fs.Parse(os.Args[i:])

}

func PathLst() ([]string, error) {

	if len(pathLstRaw) == 0 {
		return nil, ErrPathLstEmpty
	}

	pathLstTemp := strings.Split(pathLstRaw, ",")
	pathLst := make([]string, len(pathLstTemp))

	for i, p := range pathLstTemp {
		pathLst[i] = strings.Trim(p, " ")
	}

	return pathLst, nil

}

func PathIBases() ([]string, error) {

	if len(pathiBasesRaw) == 0 {
		return nil, ErrPathIBasesEmpty
	}

	pathiBasesTemp := strings.Split(pathiBasesRaw, ",")
	pathiBases := make([]string, len(pathiBasesTemp))

	for i, p := range pathiBasesTemp {
		pathiBases[i] = strings.Trim(p, " ")
	}

	return pathiBases, nil

}

func Instance() (string, error) {
	return instance, nil
}

func Usage() {
	fmt.Println("© TITOVCODE iBases Creator\nGithub: github.com/korableg/OneCIBasesCreator, E-mail: titov-de@yandex.ru, 2020")
	fmt.Println("Help:")
	fs.PrintDefaults()
	os.Exit(2)
}
