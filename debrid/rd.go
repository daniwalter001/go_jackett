package debrid

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/daniwalter001/jackett_fiber/types/rd"
	"github.com/gofiber/fiber/v2"
)

var keys = []string{"N4LZTV4IC4MUUJBFQMNFHSR4OCT4EV2SL4OTDPDPBMQYEFQML3GQ",
	"MUEFTLYM4XHWZ7MWFUYUHMS27SAOCX7D5W2HY67TZJXHDINYW5MQ",
	"YODGFXPGMJO6MEHZPILLK5OOKQ3X7W4LHBJFEGSMJMAWTDJKADBQ",
	"ZDHDDYHFB67Q2QAE3FBWKDPQ3AFGNWA2JALNWCDZR2E4QDMRS5TQ",
	"SKY26TG4NFOE4QJX55ISNVITB7Q2S7ZGKFOB34RUUJOBWHABJUHQ",
	"2ILZOV4OXUWH2V3D276BLKVV6XRACRRVH4DPL4XSDRPB2V6QALXA",
	"W665ORDQWJBT7OUT2UARA3SZYFARHVNIBBQ6ZBBEBCBBFHW5GECQ"}

func getApiKey() string {
	return keys[rand.Intn(len(keys))]
}

func bearer() string {
	return fmt.Sprintf("Bearer %s", getApiKey())
}

func checkTorrentFileinRD(hash string) (rd.AvailabilityResponse, rd.RdError) {
	if len(hash) == 0 {
		return rd.AvailabilityResponse{}, rd.RdError{}

	}
	api := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/instantAvailability/%s", hash)

	request := fiber.Get(api).Timeout(5 * time.Second)
	request.Set("Authorization", bearer())

	status, data, err := request.Bytes()

	if err != nil {
		return rd.AvailabilityResponse{}, rd.RdError{}
	}

	if status >= 400 {
		var resErr rd.RdError
		json.Unmarshal(data, &resErr)
		return rd.AvailabilityResponse{}, resErr
	}
	var resJson rd.AvailabilityResponse
	json.Unmarshal(data, &resJson)
	return resJson, rd.RdError{}

}

func addTorrentFileinRD(magnet string) (rd.AddTorrentResponse, rd.RdError) {
	if len(magnet) == 0 {
		return rd.AddTorrentResponse{}, rd.RdError{}
	}
	// magnet:?xt=urn:btih:fc30f2a7628a28330ba9e84d04992b6ea3a8d637

	api := "https://api.real-debrid.com/rest/1.0/torrents/addMagnet"

	request := fiber.Post(api).Timeout(5 * time.Second).Body([]byte(fmt.Sprintf("magnet=%s", magnet)))
	request.Set("Authorization", bearer())
	request.Set("Content-Type", "application/x-www-form-urlencoded")

	status, data, err := request.Bytes()
	if err != nil {
		return rd.AddTorrentResponse{}, rd.RdError{}
	}
	if status >= 400 {
		var resErr rd.RdError
		json.Unmarshal(data, &resErr)
		return rd.AddTorrentResponse{}, resErr
	}
	var resJson rd.AddTorrentResponse
	json.Unmarshal(data, &resJson)
	return resJson, rd.RdError{}

}

func getTorrentInfofromRD(id string) (rd.TorrentInfoResponse, rd.RdError) {
	if len(id) == 0 {
		return rd.TorrentInfoResponse{}, rd.RdError{}

	}
	api := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/info/%s", id)

	request := fiber.Get(api).Timeout(5 * time.Second)
	request.Set("Authorization", bearer())

	status, data, err := request.Bytes()

	if err != nil {
		return rd.TorrentInfoResponse{}, rd.RdError{}
	}

	if status != fiber.StatusOK {
		var resErr rd.RdError
		json.Unmarshal(data, &resErr)
		return rd.TorrentInfoResponse{}, resErr
	}
	var resJson rd.TorrentInfoResponse
	json.Unmarshal(data, &resJson)
	return resJson, rd.RdError{}

}

func selectFilefromRD(id string, files string) (bool, rd.RdError) {
	if len(id) == 0 {
		return false, rd.RdError{}
	}
	if len(files) == 0 {
		files = "all"
	}

	api := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/selectFiles/%s", id)

	request := fiber.Post(api).Timeout(5 * time.Second).Body([]byte(fmt.Sprintf("files=%s", files)))
	request.Set("Authorization", bearer())
	request.Set("Content-Type", "application/x-www-form-urlencoded")

	status, data, err := request.Bytes()
	if err != nil {
		return false, rd.RdError{}
	}
	if status >= 400 {
		var resErr rd.RdError
		json.Unmarshal(data, &resErr)
		return false, resErr
	}
	return true, rd.RdError{}
}

func unrestrictLinkfromRD(link string) (rd.UnrestrictLinkResponse, rd.RdError) {
	if len(link) == 0 {
		return rd.UnrestrictLinkResponse{}, rd.RdError{}
	}

	api := "https://api.real-debrid.com/rest/1.0/unrestrict/link"

	request := fiber.Post(api).Timeout(5 * time.Second).Body([]byte(fmt.Sprintf("link=%s", link)))
	request.Set("Authorization", bearer())
	request.Set("Content-Type", "application/x-www-form-urlencoded")

	status, data, err := request.Bytes()
	if err != nil {
		return rd.UnrestrictLinkResponse{}, rd.RdError{}
	}
	if status >= 400 {
		var resErr rd.RdError
		json.Unmarshal(data, &resErr)
		return rd.UnrestrictLinkResponse{}, resErr
	}

	var resJson rd.UnrestrictLinkResponse
	json.Unmarshal(data, &resJson)
	return resJson, rd.RdError{}
}
