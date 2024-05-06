package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/daniwalter001/jackett_fiber/types"
	"github.com/gofiber/fiber/v2"
)

func main() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)

	//read and parse cache file content
	mapCache := make(map[string]types.StreamMeta)
	cacheFile, _ := os.ReadFile("./persistence/cache.json")
	if len(cacheFile) > 0 {
		json.Unmarshal(cacheFile, &mapCache)
	}

	app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Working")
	// })

	app.Get("/manifest.json", func(c *fiber.Ctx) error {
		a := types.StreamManifest{ID: "strem.go.beta", Description: "Random Golang version on stremio Addon", Name: "GoDon", Resources: []string{"stream"}, Version: "1.0.9", Types: []string{"movie", "series", "anime"}, Logo: "https://upload.wikimedia.org/wikipedia/commons/2/23/Golang.png"}

		u, err := json.Marshal(a)
		if err != nil {
			return c.SendStatus(fiber.StatusOK)
		}

		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "*")
		c.Set("Content-Type", "application/json")

		return c.Status(fiber.StatusOK).SendString(string(u))
	})

	app.Get("/stream/:type/:id.json", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Headers", "*")
		c.Set("Content-Type", "application/json")

		fmt.Printf("Id: %s\n", c.Params("id"))
		fmt.Printf("Type: %s\n", c.Params("type"))

		var s, e, abs_season, abs_episode int
		var tt string
		abs := "false"

		id := c.Params("id")
		id = strings.ReplaceAll(id, "%3A", ":")

		//Reading the cache
		streams, exists := mapCache[id]
		if exists {
			fmt.Printf("Sending that %s shit from cache\n", id)
			return c.Status(fiber.StatusOK).JSON(streams)
		}

		type_ := c.Params("type")

		var tmp []string

		if strings.Contains(id, "kitsu") {
			tmp = getImdbFromKitsu(id)
		} else {
			tmp = strings.Split(id, ":")
		}

		fmt.Println(tmp)

		tt = tmp[0]
		if len(tmp) > 1 {
			s, _ = strconv.Atoi(tmp[1])
			e, _ = strconv.Atoi(tmp[2])
			if len(tmp) > 3 {
				abs_season, _ = strconv.Atoi(tmp[3])
				abs_episode, _ = strconv.Atoi(tmp[4])
				abs = tmp[5]
			}
		}

		// fmt.Println("----------------------------")
		// fmt.Println(tt)
		// fmt.Println(strconv.Itoa(s))
		// fmt.Println(strconv.Itoa(e))
		// fmt.Println(abs)
		// fmt.Println(strconv.Itoa(abs_season))
		// fmt.Println(strconv.Itoa(abs_episode))
		// fmt.Println("----------------------------")

		name, year := getMeta(tt, type_)

		var results []types.ItemsParsed

		wg := sync.WaitGroup{}
		l := 4
		if type_ == "series" {
			if abs == "true" {
				l = l + 1
			}
			if s == 1 {
				l = l + 1
			}
		} else if type_ == "movie" {
			l = 1
		}
		fmt.Printf("Requests: %d\n", l)

		wg.Add(l)

		//=========================================================

		if type_ == "movie" {
			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s %s", name, year), type_)...)
			}()
		} else {

			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s S%02d", name, s), type_)...)
			}()
			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s batch", name), type_)...)
			}()

			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s complete", name), type_)...)
			}()

			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s S%02dE%02d", name, s, e), type_)...)
			}()

			if s == 1 {
				go func() {
					defer wg.Done()
					results = append(results, fetchTorrent(fmt.Sprintf("%s %02d", name, e), type_)...)
				}()
			}

			if abs == "true" {
				go func() {
					defer wg.Done()
					results = append(results, fetchTorrent(fmt.Sprintf("%s %03d", name, abs_episode), type_)...)
				}()
			}
		}

		//=========================================================

		wg.Wait()

		sort.Slice(results, func(i, j int) bool {
			iv, _ := strconv.Atoi(results[i].Peers)
			jv, _ := strconv.Atoi(results[j].Peers)
			return iv > jv
		})

		results = removeDuplicates(results)

		fmt.Printf("Results:%d\n", len(results))

		if len(results) > 100 {
			results = results[:100]
		}

		// for index, el := range results {
		// 	fmt.Printf("%d. %s => %s\n", index, el.MagnetURI, el.Seeders)
		// }

		fmt.Printf("Retenus:%d\n", len(results))

		var parsedTorrentFiles []types.ItemsParsed

		wg = sync.WaitGroup{}
		wg.Add(len(results))
		for i := 0; i < len(results); i++ {
			go func(item types.ItemsParsed) {
				defer wg.Done()
				r := item
				if strings.Contains(item.MagnetURI, "magnet:?xt") {
					r = readTorrentFromMagnet(item)
				} else {
					r = readTorrent(item)
				}
				parsedTorrentFiles = append(parsedTorrentFiles, r)
			}(results[i])
		}
		wg.Wait()

		fmt.Printf("Response %d\n", len(parsedTorrentFiles))

		var parsedSuitableTorrentFiles []torrent.File
		var parsedSuitableTorrentFilesIndex = make([]int, len(parsedTorrentFiles))

		for index, el := range parsedTorrentFiles {
			parsedSuitableTorrentFiles = make([]torrent.File, 0)

			parsedSuitableTorrentFilesIndex[index] = 0

			for _index, ell := range el.TorrentData {
				lower := strings.ToLower(ell.DisplayPath())

				if !isVideo(ell.DisplayPath()) {
					continue
				}

				if type_ == "movie" {
					parsedSuitableTorrentFiles = append(parsedSuitableTorrentFiles, ell)
					parsedSuitableTorrentFilesIndex[index] = _index + 1

					break
				} else {
					if isVideo(ell.DisplayPath()) && (containEandS(lower, strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(abs_season), strconv.Itoa(abs_episode)) ||
						containE_S(lower, strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(abs_season), strconv.Itoa(abs_episode)) ||
						(s == 1 && (containsAbsoluteE(lower, strconv.Itoa(s), strconv.Itoa(e), true, strconv.Itoa(s), strconv.Itoa(e)) ||
							containsAbsoluteE_(lower, strconv.Itoa(s), strconv.Itoa(e), true, strconv.Itoa(s), strconv.Itoa(e)))) ||
						// false ||
						(((abs == "true" && containsAbsoluteE(lower, strconv.Itoa(s), strconv.Itoa(e), true, strconv.Itoa(abs_season), strconv.Itoa(abs_episode))) ||
							(abs == "true" && containsAbsoluteE_(lower, strconv.Itoa(s), strconv.Itoa(e), true, strconv.Itoa(abs_season), strconv.Itoa(abs_episode)))) &&
							!(strings.Contains(lower, "s0") && strings.Contains(lower, "e0") && strings.Contains(lower, "season") && strings.Contains(lower, fmt.Sprintf("s%d", abs_season)) && strings.Contains(lower, fmt.Sprintf("e%d", abs_episode))))) {
						parsedSuitableTorrentFiles = append(parsedSuitableTorrentFiles, ell)
						parsedSuitableTorrentFilesIndex[index] = _index + 1
						break
					}
				}
			}
			parsedTorrentFiles[index].TorrentData = parsedSuitableTorrentFiles
		}

		var ttttt types.StreamMeta
		fmt.Println("Parsing that shit")

		for ind, el := range parsedTorrentFiles {
			for _, ell := range el.TorrentData {

				if !isVideo(ell.DisplayPath()) {
					continue
				}

				announceList := make([]string, 0)
				// fmt.Println(PrettyPrint(ell.Torrent().Metainfo().AnnounceList))
				for i := 0; i < len(ell.Torrent().Metainfo().AnnounceList); i++ {
					for j := 0; j < len(ell.Torrent().Metainfo().AnnounceList[i]); j++ {
						announceList = append(announceList, fmt.Sprintf("tracker:%s", ell.Torrent().Metainfo().AnnounceList[i][j]))
					}
				}
				// if len(announceList) > 0 {
				// 	announceList = defaultTracker()
				// }
				announceList = append(announceList, fmt.Sprintf("dht:%s", ell.Torrent().InfoHash().String()))

				ttttt.Streams = append(ttttt.Streams, types.TorrentStreams{Title: fmt.Sprintf("%s\n%s\n%s | %s", ell.Torrent().Name(), ell.DisplayPath(), getQuality(ell.DisplayPath()), getSize(int(ell.Length()))), Name: fmt.Sprintf("%s\n S:%s, P:%s", el.Tracker, el.Seeders, el.Peers), Type: type_, InfoHash: ell.Torrent().InfoHash().String(), Sources: announceList, BehaviorHints: types.BehaviorHints{BingeGroup: fmt.Sprintf("Jackett|%s", ell.Torrent().InfoHash().String()), NotWebReady: true}, FileIdx: parsedSuitableTorrentFilesIndex[ind] - 1})
				break
			}
		}

		if len(ttttt.Streams) > 0 {
			mapCache[id] = ttttt
			toFile, _ := json.MarshalIndent(mapCache, "", " ")
			os.WriteFile("./persistence/cache.json", toFile, 0666)
		}

		fmt.Println("Sending that shit")
		return c.Status(fiber.StatusOK).JSON(ttttt)
	})

	app.Listen(":3000")
}
