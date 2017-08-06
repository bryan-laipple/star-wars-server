package etl

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/bryan-laipple/star-wars-server/storage"
)

const imageScraper = "/src/github.com/bryan-laipple/star-wars-server/etl/wookieepediaImageScraper.sh"

type scraperOutput struct {
	Images []string `json:"images"`
}

func NameToWookieepediaLink(name string) storage.Link {
	underscore := strings.Replace(name, " ", "_", -1)
	href := strings.Join([]string{"http://starwars.wikia.com/wiki/", underscore}, "")
	return storage.Link{"wookieepedia", href}
}

func ScrapeImageLinks(wookieepediaUrl string) []string {
	name := strings.Join([]string{os.Getenv("GOPATH"), imageScraper}, "")
	out, err := exec.Command(name, wookieepediaUrl).Output()
	if err != nil {
		log.Fatal(err)
	}
	var dat scraperOutput
	if err := json.Unmarshal(out, &dat); err != nil {
		panic(err)
	}
	return dat.Images
}
