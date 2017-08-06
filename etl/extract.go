package etl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/bryan-laipple/star-wars-server/storage"
	"github.com/leejarvis/swapi"
)

type swData struct {
	Characters []storage.Character `json:"characters"`
	Starships  []storage.Starship  `json:"starships"`
	Planets    []storage.Planet    `json:"planets"`
}

func Extract() {
	fmt.Println("Extracting data from swapi.co and starwars.wikia.com...")
	swData := &swData{}
	extract(swData)
	jsonBytes, _ := json.MarshalIndent(swData, "", "  ")
	timestamp := time.Now().Format(time.RFC3339)
	filename := strings.Join([]string{"./etl/extracted-", timestamp, ".json"}, "")
	err := ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		panic(err)
	}
	//fmt.Println(string(jsonBytes))
}

func urlToId(url string) string {
	offset := 1
	if strings.HasSuffix(url, "/") {
		offset = 2
	}
	split := strings.Split(url, "/")
	return split[len(split)-offset]
}

func extract(data *swData) {
	urlToName := make(map[string]string)
	characters, err := storage.GetCharacters()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, character := range characters {
		urlToName[character.URL] = character.Name
	}

	films, err := storage.GetFilms()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, film := range films {
		urlToName[film.URL] = film.Title
	}

	starships, err := storage.GetStarships()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship.URL] = starship.Name
	}

	species, err := storage.GetSpecies()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, sp := range species {
		urlToName[sp.URL] = sp.Name
	}

	planets, err := storage.GetPlanets()
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet.URL] = planet.Name
	}

	for i, _ := range characters {
		character := &characters[i]
		character.Id = urlToId(character.URL)
		character.Type = "character"
		character.Homeworld = urlToName[character.Homeworld]
		wookieepediaLink := NameToWookieepediaLink(character.Name)
		character.Links = append(character.Links, wookieepediaLink)
		character.Images = ScrapeImageLinks(wookieepediaLink.Href)
		for i, film := range character.Films {
			character.Films[i] = urlToName[film]
		}
		for i, sp := range character.Species {
			character.Species[i] = urlToName[sp]
		}
		for i, starship := range character.Starships {
			character.Starships[i] = urlToName[starship]
		}
	}

	for i, _ := range starships {
		starship := &starships[i]
		starship.Id = urlToId(starship.URL)
		starship.Type = "starship"
		wookieepediaLink := NameToWookieepediaLink(starship.Name)
		starship.Links = append(starship.Links, wookieepediaLink)
		starship.Images = ScrapeImageLinks(wookieepediaLink.Href)
		for i, film := range starship.Films {
			starship.Films[i] = urlToName[film]
		}
		for i, pilot := range starship.Pilots {
			starship.Pilots[i] = urlToName[pilot]
		}
	}

	for i, _ := range planets {
		planet := &planets[i]
		planet.Id = urlToId(planet.URL)
		planet.Type = "planet"
		wookieepediaLink := NameToWookieepediaLink(planet.Name)
		planet.Links = append(planet.Links, wookieepediaLink)
		planet.Images = ScrapeImageLinks(wookieepediaLink.Href)
		for i, film := range planet.Films {
			planet.Films[i] = urlToName[film]
		}
		for i, resident := range planet.Residents {
			planet.Residents[i] = urlToName[resident]
		}
	}

	data.Characters = characters
	data.Starships = starships
	data.Planets = planets
}

//
// Some experimenting below with generic results structure
//
type pagedResponse struct {
	Count    int                      `json:"count"`
	Next     string                   `json:"next"`
	Previous string                   `json:"previous"`
	Results  []map[string]interface{} `json:"results"`
}

func getList(url string) (list []map[string]interface{}, err error) {
	for url != "" {
		var res pagedResponse
		if err = swapi.Get(url, &res); err != nil {
			return
		}
		url = res.Next
		list = append(list, res.Results...)
	}
	return
}

func extractToGenericMap() map[string]string {
	urlToName := make(map[string]string)
	characters, err := getList("https://swapi.co/api/people/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, character := range characters {
		urlToName[character["url"].(string)] = character["name"].(string)
	}

	//films, err := getList("https://swapi.co/api/films/")
	//if err != nil {
	//	fmt.Printf("some error occured")
	//}
	//for _, film := range films {
	//	urlToName[film["url"].(string)] = film["title"].(string)
	//}

	starships, err := getList("https://swapi.co/api/starships/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, starship := range starships {
		urlToName[starship["url"].(string)] = starship["name"].(string)
	}

	//species, err := getList("https://swapi.co/api/species/")
	//if err != nil {
	//	fmt.Printf("some error occured")
	//}
	//for _, sp := range species {
	//	urlToName[sp["url"].(string)] = sp["name"].(string)
	//}

	planets, err := getList("https://swapi.co/api/planets/")
	if err != nil {
		fmt.Printf("some error occured")
	}
	for _, planet := range planets {
		urlToName[planet["url"].(string)] = planet["name"].(string)
	}

	//fmt.Printf("%+v\n", urlToName)
	return urlToName
}
