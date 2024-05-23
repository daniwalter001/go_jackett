package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/daniwalter001/jackett_fiber/types"
	"github.com/daniwalter001/jackett_fiber/types/rd"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	ctx := context.Background()

	errDot := godotenv.Load("./.env")
	if errDot != nil {
		log.Fatalln("Error loading .env file")
	}

	//create redis client instance
	rdClient := RedisClient()
	status, errS := rdClient.Ping(ctx).Result()
	if errS != nil {
		fmt.Print("Error: ")
		fmt.Println(errS.Error())
	} else {
		fmt.Print("OK redis: ")
		fmt.Println(status)
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("Working")
	})

	app.Get("/manifest.json", func(c *fiber.Ctx) error {
		a := types.StreamManifest{ID: "strem.go.beta", Description: "Random Golang version on stremio Addon", Name: "GoDon", Resources: []string{"stream"}, Version: "1.0.9", Types: []string{"movie", "series", "anime"}, Logo: "https://upload.wikimedia.org/wikipedia/commons/2/23/Golang.png", IdPrefixes: []string{"tt", "kitsu"}, Catalogs: []string{}}

		// {
		// 	"catalogs": [  ],
		// 	"description": "VOD from Google Drive.",
		// 	"id": "hy.stremio.googledrive",
		// 	"logo": "https://raw.githubusercontent.com/mik25/stremio-greek-tv/master/pngwing.com.png",
		// 	"name": "GDrive Reborn",
		// 	"resources": [
		// 	  "stream"
		// 	],
		// 	"types": [
		// 	  "movie",
		// 	  "series"
		// 	],
		// 	"version": "2.0.0"
		//   }

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
		//Reading the cache
		streams, err := rdClient.JSONGet(ctx, id, "$").Result()
		if err == nil && streams != "" {
			fmt.Printf("Sending that %s shit from cache\n", id)
			var cachedStreams []types.StreamMeta
			errJson := json.Unmarshal([]byte(streams), &cachedStreams)
			if errJson != nil {
				fmt.Println(errJson)
				return c.Status(fiber.StatusNotFound).SendString("lol")
			} else if len(cachedStreams) > 0 {
				fmt.Printf("Sent from cache %s\n", id)
				return c.Status(fiber.StatusOK).JSON(cachedStreams[len(cachedStreams)-1])
			}
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
		l := 5
		if type_ == "series" {
			if abs == "true" {
				l = l + 2
			}
			if s == 1 {
				l = l + 2
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
				results = append(results, fetchTorrent(fmt.Sprintf("%s integrale", name), type_)...)
			}()
			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s batch", name), type_)...)
			}()
			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s complet", name), type_)...)
			}()

			go func() {
				defer wg.Done()
				results = append(results, fetchTorrent(fmt.Sprintf("%s S%02dE%02d", name, s, e), type_)...)
			}()

			if s == 1 {
				go func() {
					defer wg.Done()
					results = append(results, fetchTorrent(fmt.Sprintf("%s E%02d", name, e), type_)...)
				}()
				go func() {
					defer wg.Done()
					results = append(results, fetchTorrent(fmt.Sprintf("%s %02d", name, e), type_)...)
				}()
			}

			if abs == "true" {
				go func() {
					defer wg.Done()
					results = append(results, fetchTorrent(fmt.Sprintf("%s E%03d", name, abs_episode), type_)...)
				}()

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

		maxRes, _ := strconv.Atoi(os.Getenv("MAX_RES"))

		if len(results) > maxRes {
			results = results[:maxRes]
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

				if len(r.TorrentData) != 0 {
					parsedTorrentFiles = append(parsedTorrentFiles, r)
				}
			}(results[i])
		}
		wg.Wait()

		var parsedSuitableTorrentFiles []torrent.File
		// var parsedSuitableTorrentFilesIndex = make([]int, len(parsedTorrentFiles))
		var parsedSuitableTorrentFilesIndex = make(map[string]int, 0)

		for index, el := range parsedTorrentFiles {
			parsedSuitableTorrentFiles = make([]torrent.File, 0)

			for _index, ell := range el.TorrentData {
				lower := strings.ToLower(ell.DisplayPath())

				if !isVideo(ell.DisplayPath()) {
					continue
				}

				if type_ == "movie" {
					parsedSuitableTorrentFiles = append(parsedSuitableTorrentFiles, ell)
					parsedSuitableTorrentFilesIndex[ell.DisplayPath()] = _index + 1

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
						parsedSuitableTorrentFilesIndex[ell.DisplayPath()] = _index + 1
						break
					}
				}
			}
			parsedTorrentFiles[index].TorrentData = parsedSuitableTorrentFiles
		}

		var ttttt types.StreamMeta

		parsedTorrentFiles = filter[types.ItemsParsed](parsedTorrentFiles, func(ip types.ItemsParsed) bool {
			return len(ip.TorrentData) != 0
		})

		fmt.Printf("Response %d\n", len(parsedTorrentFiles))

		fmt.Println("Parsing that shit")

		wg = sync.WaitGroup{}
		wg.Add(len(parsedTorrentFiles))

		nbreAdded := 0

		for _, el := range parsedTorrentFiles {
			go func(item types.ItemsParsed) {
				defer wg.Done()
				for _, ell := range el.TorrentData {
					fmt.Println(ell.DisplayPath())
					if !isVideo(ell.DisplayPath()) {
						continue
					}

					// ========================== RD =============================
					fmt.Printf("Trynna some RD...\n")

					infoHash := ell.Torrent().InfoHash().String()
					// magnet := fmt.Sprint(ell.Torrent().Metainfo().Magnet(nil, ell.Torrent().Info()))
					var folderId string
					var details []rd.UnrestrictLinkResponse
					var data rd.AddTorrentResponse

					available, err := checkTorrentFileinRD(infoHash)

					if err.Error != "" {
						continue
					}

					v, availableCheck := available[infoHash]

					if !availableCheck {
						continue
					}

					v_ := v["rd"]
					availableCheck = len(v_) > 0

					if availableCheck || nbreAdded < 3 {
						data, err = addTorrentFileinRD2(fmt.Sprintf("magnet:?xt=urn:btih:%s", infoHash))
						if availableCheck {
							fmt.Println("Cached")
						} else {
							fmt.Println("Added")
							nbreAdded = nbreAdded + 1
						}
					}

					folderId = data.ID
					selected, err := selectFilefromRD(folderId, "all")
					if folderId != "" && selected {
						torrentDetails, err_ := getTorrentInfofromRD(folderId)
						//fmt.Println((PrettyPrint(torrentDetails)))
						if err.Error != "" {
							fmt.Println("Error")
							fmt.Println(err_.Error)
						}
						var files []rd.Files
						if len(torrentDetails.Files) > 0 {
							files = filter[rd.Files](torrentDetails.Files, func(f rd.Files) bool {
								return f.Selected == 1
							})
							links := torrentDetails.Links
							selectedIndex := 0

							if len(files) > 1 {
								selectedIndex = slices.IndexFunc[[]rd.Files](files, func(t rd.Files) bool {
									return strings.Contains(strings.ToLower(t.Path), strings.ToLower(ell.DisplayPath()))
								})
							}
							if selectedIndex == -1 || len(links) <= selectedIndex {
								selectedIndex = 0
							}
							if len(links) > 0 {
								unrestrictLink, errun := unrestrictLinkfromRD(links[selectedIndex])
								details = append(details, unrestrictLink)
								// fmt.Println(PrettyPrint(errun.Error))
								// fmt.Println(PrettyPrint(details[len(details)-1]))
								if errun.Error != "" {
									continue
								}
							}
						}

					}

					if len(details) > 0 {
						ttttt.Streams = append(ttttt.Streams, types.TorrentStreams{Title: fmt.Sprintf("%s\n%s\n%s | %s", ell.Torrent().Name(), ell.DisplayPath(), getQuality(ell.DisplayPath()), getSize(int(ell.Length()))), Name: fmt.Sprintf("RD.%s\n S:%s, P:%s", item.Tracker, item.Seeders, item.Peers), Type: type_, BehaviorHints: types.BehaviorHints{BingeGroup: fmt.Sprintf("Jackett|%s", ell.Torrent().InfoHash().String()), NotWebReady: true}, URL: details[0].Download})

						// ========================== END RD =============================
					} else if os.Getenv("PUBLIC") == "1" {
						announceList := make([]string, 0)
						for i := 0; i < len(ell.Torrent().Metainfo().AnnounceList); i++ {
							for j := 0; j < len(ell.Torrent().Metainfo().AnnounceList[i]); j++ {
								announceList = append(announceList, fmt.Sprintf("tracker:%s", ell.Torrent().Metainfo().AnnounceList[i][j]))
							}
						}
						announceList = append(announceList, fmt.Sprintf("dht:%s", ell.Torrent().InfoHash().String()))
						ttttt.Streams = append(ttttt.Streams, types.TorrentStreams{Title: fmt.Sprintf("%s\n%s\n%s | %s", ell.Torrent().Name(), ell.DisplayPath(), getQuality(ell.DisplayPath()), getSize(int(ell.Length()))), Name: fmt.Sprintf("%s\n S:%s, P:%s", item.Tracker, item.Seeders, item.Peers), Type: type_, InfoHash: ell.Torrent().InfoHash().String(), Sources: announceList, BehaviorHints: types.BehaviorHints{BingeGroup: fmt.Sprintf("Jackett|%s", ell.Torrent().InfoHash().String()), NotWebReady: true}, FileIdx: parsedSuitableTorrentFilesIndex[ell.DisplayPath()] - 1})
					}

				}
			}(el)
		}

		wg.Wait()

		if len(ttttt.Streams) > 0 {
			jsonBytes, errttt := json.Marshal(ttttt)
			if errttt == nil {
				_, errrrr := rdClient.JSONSet(ctx, id, "$", jsonBytes).Result()
				if errrrr == nil {
					rdClient.Expire(ctx, id, time.Hour*24*7).Result()
				}
			}
		}

		fmt.Println("Sending that shit")
		return c.Status(fiber.StatusOK).JSON(ttttt)
	})

	app.Listen(":8000")
}
