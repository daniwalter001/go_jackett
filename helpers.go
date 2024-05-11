package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"unicode"

	"github.com/daniwalter001/jackett_fiber/types"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func PrettyPrint(i interface{}) string {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	s, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return ""
	}
	return string(s)
}

func containEandS(name string, s string, e string, abs bool, abs_season string, abs_episode string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, fmt.Sprintf("s%02se%02s ", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%02se%02s.", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%02se%02s-", s, e)) ||
		// -----
		strings.Contains(lower, fmt.Sprintf("s%se%02s ", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%se%02s.", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%se%02s-", s, e)) ||
		// -----
		strings.Contains(lower, fmt.Sprintf("s%sx%s", s, e)) ||
		// -----
		strings.Contains(lower, fmt.Sprintf("s%02s - e%02s ", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%02s.e%02s ", s, e)) ||
		// -----
		strings.Contains(lower, fmt.Sprintf("s%02se%s ", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%02se%s.", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%02se%s-", s, e)) ||
		// -----
		strings.Contains(lower, fmt.Sprintf("season %s e%s", s, e)) ||
		// ----- abs
		(abs &&
			(strings.Contains(lower, fmt.Sprintf("s%02se%02s", abs_season, abs_episode)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%03s", abs_season, abs_episode)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%04s", abs_season, abs_episode)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%02s", s, abs_episode)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%03s", s, abs_episode)) ||
				strings.Contains(lower, fmt.Sprintf("season %s e%s", s, e)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%s-", s, e)) ||
				strings.Contains(lower, fmt.Sprintf("s%02se%s-", s, e)) ||
				false))
}

func containE_S(name string, s string, e string, abs bool, abs_season string, abs_episode string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, fmt.Sprintf("s%02s - %02s", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("s%s - %02s", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("season %s - %02s", s, e)) ||
		strings.Contains(lower, fmt.Sprintf("season %s - %03s", s, e))
}

func containsAbsoluteE(name string, s string, e string, abs bool, abs_season string, abs_episode string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, fmt.Sprintf("%02s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("%03s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("0%s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("%04s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("e%02s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("e%03s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("e0%s ", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf("e%04s ", abs_episode))
}

func containsAbsoluteE_(name string, s string, e string, abs bool, abs_season string, abs_episode string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, fmt.Sprintf(" %02s.", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" %03s.", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" 0%s.", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" %04s.", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" %02s-", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" %03s-", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" 0%s-", abs_episode)) ||
		strings.Contains(lower, fmt.Sprintf(" %04s-", abs_episode))
}

func isVideo(name string) bool {
	lower := strings.ToLower(name)

	return strings.Contains(lower, ".mp4") ||
		strings.Contains(lower, ".mkv") ||
		strings.Contains(lower, ".avi") ||
		strings.Contains(lower, ".ts") ||
		strings.Contains(lower, ".m3u") ||
		strings.Contains(lower, ".m3u8") ||
		strings.Contains(lower, ".flv")
}

func getQuality(name string) string {
	lower := strings.ToLower(name)

	if slices.ContainsFunc([]string{"2160", "4k", "uhd"}, func(e string) bool {
		return strings.Contains(lower, e)
	}) {
		return " 🌟4k"
	}

	if slices.ContainsFunc([]string{"1080", "fhd"}, func(e string) bool {
		return strings.Contains(lower, e)
	}) {
		return " 🎥FHD"
	}

	if slices.ContainsFunc([]string{"720", "hd"}, func(e string) bool {
		return strings.Contains(lower, e)
	}) {
		return " 📺HD"
	}

	if slices.ContainsFunc([]string{"480p", "380p", "sd"}, func(e string) bool {
		return strings.Contains(lower, e)
	}) {
		return " 📱SD"
	}

	return ""
}

// function getSize(size) {
// 	var gb = 1024 * 1024 * 1024;
// 	var mb = 1024 * 1024;

// 	return (
// 	  "💾 " +
// 	  (size / gb > 1
// 		? `${(size / gb).toFixed(2)} GB`
// 		: `${(size / mb).toFixed(2)} MB`)
// 	);
//   }

func getSize(size int) string {
	gb := 1024 * 1024 * 1024
	mb := 1024 * 1024
	kb := 1024

	size_ := "💾 "

	if size/gb > 1 {
		size_ = size_ + fmt.Sprintf("%.2f GB", float32(size/gb))
	} else if size/mb > 1 {
		size_ = size_ + fmt.Sprintf("%.2f MB", float32(size/mb))
	} else {
		size_ = size_ + fmt.Sprintf("%.2f KB", float32(size/kb))
	}
	return size_

}

func removeDuplicates(strList []types.ItemsParsed) []types.ItemsParsed {
	list := []types.ItemsParsed{}
	for _, item := range strList {
		if !(contains(list, item)) {
			list = append(list, item)
		}
	}
	return list
}

func contains(s []types.ItemsParsed, e types.ItemsParsed) bool {
	for _, a := range s {
		if a.Title == e.Title {
			return true
		}
	}
	return false
}

func filter[T any](slice []T, cb func(T) bool) (ret []T) {
	for _, v := range slice {
		if cb(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func getServers() []types.Server {
	var servers []types.Server
	readFile, errServer := os.Open("./assets/servers.db")

	if errServer != nil {
		fmt.Println("Cant load server file")
		fmt.Println(errServer)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		servers = append(servers, types.Server{
			Host:   strings.Split(fileScanner.Text(), "|")[0],
			ApiKey: strings.Split(fileScanner.Text(), "|")[1],
		})
	}

	return servers
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output
}
