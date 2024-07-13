package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/daniwalter001/jackett_fiber/types"
	"github.com/gofiber/fiber/v2"
	gotorrentparser "github.com/j-muller/go-torrent-parser"
)

func getMeta(id string, type_ string) (string, string) {

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	splitedId := strings.Split(id, ":")
	// api := "https://v3-cinemeta.strem.io/meta/" + type_ + "/" + splitedId[0] + ".json"
	api := "https://cinemeta-live.strem.io/meta/" + type_ + "/" + splitedId[0] + ".json"
	fmt.Println(api)
	request := fiber.Get(api)

	status, data, err := request.Bytes()

	if err != nil {
		fmt.Println(err)
		return "", ""
	}

	fmt.Printf("Status code: %d\n", status)

	if status >= 400 {
		return "", ""
	}

	var res types.IMDBMeta

	jsonErr := json.Unmarshal(data, &res)

	if jsonErr != nil {
		return "", ""
	}

	var year string

	if res.Meta.Year != nil {
		year = *res.Meta.Year
	} else if res.Meta.ReleaseInfo != nil {
		year = (*res.Meta.ReleaseInfo)[:4]
	} else {
		year = ""
	}

	return *res.Meta.Name, year
}

func getImdbFromKitsu(id string) []string {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	splitedId := strings.Split(id, ":")
	api := "https://anime-kitsu.strem.fun/meta/anime/" + splitedId[0] + ":" + splitedId[1] + ".json"
	request := fiber.Get(api)
	status, data, err := request.Bytes()
	if err != nil {
		// panic(err)
		fmt.Println(PrettyPrint(err))
		return make([]string, 0)
	}

	fmt.Printf("Status code: %d\n", status)

	if status >= 400 {
		return make([]string, 0)
	}

	var res types.KitsuMeta

	jsonErr := json.Unmarshal(data, &res)

	if jsonErr != nil {
		panic(jsonErr)
	}

	imdb := res.Meta.ImdbID
	var meta types.Videos

	for i := 0; i < len(res.Meta.Videos); i++ {
		a := res.Meta.Videos[i]

		if a.ID == id {
			meta = res.Meta.Videos[i]
		}
	}
	var resArray []string

	var e int
	var abs string

	if meta.Episode != meta.ImdbSeason || meta.ImdbSeason == 1 {
		abs = "true"
	} else {
		abs = "false"
	}

	if meta.ImdbSeason == 1 {
		e = meta.ImdbEpisode
	} else {
		e = meta.Episode
	}

	resArray = append(resArray, imdb, fmt.Sprint(meta.ImdbSeason), fmt.Sprint(meta.ImdbEpisode), fmt.Sprint(meta.Season), fmt.Sprint(e), abs)

	return resArray

}

func fetchTorrent(query string, type_ string) []types.ItemsParsed {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	servers := getServers()
	randomInt := rand.Intn(len(servers))
	host := servers[randomInt].Host
	apiKey := servers[randomInt].ApiKey
	//
	category := "5000"
	if type_ == "movie" {
		category = "2000"
	}
	query = removeAccents(strings.ReplaceAll(query, " ", "+"))

	override := os.Getenv("OVERRIDE_API_URL")
	api := fmt.Sprintf("%s/api/v2.0/indexers/yggtorrent/results/torznab/api?cache=false&cat=%s&apikey=%s&q=%s", host, category, apiKey, query)

	if override != "" {
		api = fmt.Sprintf("%s%s&apikey=%s&q=%s", host, override, apiKey, query)
	}

	fmt.Println(api)

	request := fiber.Get(api)

	status, data, err := request.Bytes()
	if err != nil {
		return make([]types.ItemsParsed, 0)
	}
	fmt.Printf("Status code: %d\n", status)
	if status >= 400 {
		return make([]types.ItemsParsed, 0)
	}

	var res types.JackettRssReponse

	xmlErr := xml.Unmarshal(data, &res)

	if xmlErr != nil {
		panic(xmlErr)
	}

	items := res.Channel.Item
	var parsedItems []types.ItemsParsed
	for i := 0; i < len(items); i++ {
		var a types.ItemsParsed
		a.Title = items[i].Title
		a.Link = items[i].Enclosure.URL
		a.Tracker = items[i].Jackettindexer.Text
		a.MagnetURI = items[i].Link
		attr := items[i].Attr
		for ii := 0; ii < len(attr); ii++ {
			if attr[ii].Name == "seeders" {
				a.Seeders = attr[ii].Value
			}
			if attr[ii].Name == "peers" {
				a.Peers = attr[ii].Value
			}
		}
		parsedItems = append(parsedItems, a)
	}

	// fmt.Println(PrettyPrint(parsedItems))

	return parsedItems
}

func readTorrent(item types.ItemsParsed) types.ItemsParsed {
	url := item.Link

	request := fiber.Get(url).Timeout(15 * time.Second)

	status, data, err := request.Bytes()

	if status >= 400 {
		return item
	}

	if err != nil {
		fmt.Printf("%s\n", err)
		return item
	}

	if err != nil {
		fmt.Println(err)
		return item
	}

	var files []types.TorrentFile

	fileReader := bytes.NewReader(data)
	torrentFile, _ := gotorrentparser.Parse(fileReader)

	for _, file := range torrentFile.Files {

		files = append(files, types.TorrentFile{
			Name:         file.Path[len(file.Path)-1],
			TorrentName:  file.Path[len(file.Path)-1],
			Path:         "/" + file.Path[len(file.Path)-1],
			Length:       int(file.Length),
			AnnounceList: torrentFile.Announce,
			InfoHash:     torrentFile.InfoHash,
		})
	}

	item.TorrentData = files

	return item

}

func readTorrentFromMagnet(item types.ItemsParsed) types.ItemsParsed {

	// fmt.Println(item.Link)

	c := TorrentClient()

	t, addErr := c.AddMagnet(item.MagnetURI)

	if addErr != nil {
		fmt.Printf("ErrMagnet: %s\n", addErr)
		return item
	}

	ed := make(chan string, 1)
	go func() {
		<-t.GotInfo()
		ed <- "done"
	}()

	select {
	case <-time.After(15 * time.Second):
		return item
	case res := <-ed:
		if res == "done" {

			var files []types.TorrentFile
			for i := 0; i < len(t.Files()); i++ {
				file := t.Files()[i]
				var announceList []string
				for i := 0; i < len(file.Torrent().Metainfo().AnnounceList); i++ {
					for j := 0; j < len(file.Torrent().Metainfo().AnnounceList[i]); j++ {
						announceList = append(announceList, fmt.Sprintf("tracker:%s", file.Torrent().Metainfo().AnnounceList[i][j]))
					}
				}
				files = append(files, types.TorrentFile{
					Name:         file.DisplayPath(),
					TorrentName:  file.Torrent().Name(),
					Path:         file.Path(),
					Length:       int(file.Length()),
					AnnounceList: announceList,
					InfoHash:     file.Torrent().InfoHash().String(),
				})
			}
			item.TorrentData = files
		}
		return item
	}

}
