package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/daniwalter001/jackett_fiber/types/rd"
	"github.com/gofiber/fiber/v2"
)

var keys = []string{"JMY424SDLO42UT46TXHBSSIFJSDYZJO3PVL5JDA7CKTNIG7YWFWA",
	"3ZQITZP34YX3M2DKHF6TUPRKD5FUD7EPPORFV65ZZVBQBAI6BQAA",
	"TUCIGWCX5VJCPB5YPAD64NB25TZFFGAWGDVHELHZDLNUJEGX45BA",
}

var rdApikey = keys[rand.Intn(len(keys))]

func bearer() string {
	return fmt.Sprintf("Bearer %s", rdApikey)
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

func addTorrentFileinRD2(magnet string) (rd.AddTorrentResponse, rd.RdError) {
	if len(magnet) == 0 {
		return rd.AddTorrentResponse{}, rd.RdError{Error: "magnet not defined"}
	}

	url := "https://api.real-debrid.com/rest/1.0/torrents/addMagnet"

	payload := strings.NewReader(fmt.Sprintf("magnet=%s", magnet))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "insomnia/8.6.1")
	req.Header.Add("Authorization", bearer())

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		var resErr rd.RdError
		json.Unmarshal(body, &resErr)
		return rd.AddTorrentResponse{}, resErr
	}
	var resJson rd.AddTorrentResponse
	json.Unmarshal(body, &resJson)
	defer res.Body.Close()
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
		return false, rd.RdError{Error: "id not defined"}
	}
	if len(files) == 0 {
		files = "all"
	}

	api := fmt.Sprintf("https://api.real-debrid.com/rest/1.0/torrents/selectFiles/%s", id)

	payload := strings.NewReader(fmt.Sprintf("files=%s", files))

	req, _ := http.NewRequest("POST", api, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "insomnia/8.6.1")
	req.Header.Add("Authorization", bearer())

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		var resErr rd.RdError
		json.Unmarshal(body, &resErr)
		return false, resErr
	}
	defer res.Body.Close()
	return true, rd.RdError{}

}

func unrestrictLinkfromRD(link string) (rd.UnrestrictLinkResponse, rd.RdError) {
	if len(link) == 0 {
		return rd.UnrestrictLinkResponse{}, rd.RdError{}
	}

	api := "https://api.real-debrid.com/rest/1.0/unrestrict/link"

	payload := strings.NewReader(fmt.Sprintf("link=%s", link))

	req, _ := http.NewRequest("POST", api, payload)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "insomnia/8.6.1")
	req.Header.Add("Authorization", bearer())

	res, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 400 {
		var resErr rd.RdError
		json.Unmarshal(body, &resErr)
		return rd.UnrestrictLinkResponse{}, resErr
	}

	var resJson rd.UnrestrictLinkResponse
	json.Unmarshal(body, &resJson)
	defer res.Body.Close()
	return resJson, rd.RdError{}

}
