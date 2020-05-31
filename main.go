package main

import (
	"OneCIBasesCreator/OneCIBasesCreator"
	"OneCIBasesCreator/args"
	"fmt"
)

func main() {

	pathLst, err := args.PathLst()
	handleError(err)

	pathIBases, err := args.PathIBases()
	handleError(err)

	err = OneCIBasesCreator.Create(pathLst, pathIBases)
	handleError(err)

}

func handleError(err error) {
	if err != nil {
		fmt.Println("Error: " + err.Error())
		args.Usage()
	}
}
