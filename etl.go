package main

import (
	"fmt"
	"io/ioutil"

	"github.com/bryan-laipple/star-wars-server/etl"
)

func print(filename string) {
	jsonBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}

func main() {
	//etl.Extract()
	etl.Load(etl.Transform("./etl/extracted.json"))
	print("./etl/extracted.json")
}
