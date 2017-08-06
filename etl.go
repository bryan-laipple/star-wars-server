package main

import (
	"github.com/bryan-laipple/star-wars-server/etl"
)

func main() {
	//etl.Extract()
	etl.Load(etl.Transform("./etl/extracted.json"))
}
