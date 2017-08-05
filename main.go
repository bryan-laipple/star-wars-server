package main

import (
	//"github.com/bryan-laipple/star-wars-server/server"
	"github.com/bryan-laipple/star-wars-server/etl"
)

func main() {
	//server.Start(8080);
	etl.BuildStarWarsDB()
}
