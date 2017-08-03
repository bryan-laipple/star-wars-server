package main

import (
	//"github.com/bryan-laipple/star-wars-server/server"
	"fmt"
	"github.com/bryan-laipple/star-wars-server/etl"
)

func main() {
	//server.Start(8080);
	//getAndPrint("https://swapi.co/api/people/")
	etl.DataStuff()
}

func getAndPrint(url string) {
	var list []map[string]interface{}
	var err error
	if list, err = etl.GetList(url); err != nil {
		fmt.Printf("some error occurred")
		return
	}
	fmt.Printf("%+v\n", list)
	//for _, one := range list {
	//	fmt.Printf("%+v\n", one)
	//}
}
