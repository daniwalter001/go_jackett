package main

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/daniwalter001/jackett_fiber/types"
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

func defaultTracker() []string {
	return []string{
		"udp://open.stealth.si:80/announce",
		"udp://exodus.desync.com:6969/announce",
		"udp://tracker.cyberia.is:6969/announce",
		"udp://tracker.opentrackr.org:1337/announce",
		"udp://tracker.torrent.eu.org:451/announce",
		"udp://explodie.org:6969/announce",
		"udp://tracker.birkenwald.de:6969/announce",
		"udp://tracker.moeking.me:6969/announce",
		"udp://ipv4.tracker.harry.lu:80/announce",
		"udp://odd-hd.fr:6969/announce",
		"udp://tracker.jamesthebard.net:6969/announce",
		"udp://tracker.picotorrent.one:6969/announce",
		"udp://tracker.0x7c0.com:6969/announce",
		"udp://oh.fuuuuuck.com:6969/announce",
		"udp://ttk2.nbaonlineservice.com:6969/announce",
		"udp://open.demonii.com:1337/announce",
		"udp://tracker.tiny-vps.com:6969/announce",
		"udp://p4p.arenabg.com:1337/announce",
		"udp://tracker.dler.org:6969/announce",
		"udp://movies.zsw.ca:6969/announce",
		"udp://tracker.openbittorrent.com:6969/announce",
		"udp://uploads.gamecoast.net:6969/announce",
		"udp://ipv6.tracker.harry.lu:80/announce",
		"udp://tracker1.bt.moack.co.kr:80/announce",
		"udp://opentracker.i2p.rocks:6969/announce",
		"udp://eddie4.nl:6969/announce",
		"udp://bt1.archive.org:6969/announce",
		"udp://tracker.swateam.org.uk:2710/announce",
		"http://tracker.openbittorrent.com:80/announce",
		"http://tracker.opentrackr.org:1337/announce",
		"https://tracker1.520.jp:443/announce",
		"https://tracker.tamersunion.org:443/announce",
		"https://tracker.imgoingto.icu:443/announce",
		"http://nyaa.tracker.wf:7777/announce",
		"udp://tracker2.dler.org:80/announce",
		"udp://tracker.theoks.net:6969/announce",
		"udp://tracker.dump.cl:6969/announce",
		"udp://tracker.bittor.pw:1337/announce",
		"udp://tracker.4.babico.name.tr:3131/announce",
		"udp://sanincode.com:6969/announce",
		"udp://retracker01-msk-virt.corbina.net:80/announce",
		"udp://private.anonseed.com:6969/announce",
		"udp://open.free-tracker.ga:6969/announce",
		"udp://isk.richardsw.club:6969/announce",
		"udp://htz3.noho.st:6969/announce",
		"udp://epider.me:6969/announce",
		"udp://bt.ktrackers.com:6666/announce",
		"udp://acxx.de:6969/announce",
		"udp://aarsen.me:6969/announce",
		"udp://6ahddutb1ucc3cp.ru:6969/announce",
		"udp://yahor.of.by:6969/announce",
		"udp://v2.iperson.xyz:6969/announce",
		"udp://tracker1.myporn.club:9337/announce",
		"udp://tracker.therarbg.com:6969/announce",
		"udp://tracker.qu.ax:6969/announce",
		"udp://tracker.publictracker.xyz:6969/announce",
		"udp://tracker.netmap.top:6969/announce",
		"udp://tracker.farted.net:6969/announce",
		"udp://tracker.cubonegro.lol:6969/announce",
		"udp://tracker.ccp.ovh:6969/announce",
		"udp://thouvenin.cloud:6969/announce",
		"udp://thinking.duckdns.org:6969/announce",
		"udp://tamas3.ynh.fr:6969/announce",
		"udp://ryjer.com:6969/announce",
		"udp://run.publictracker.xyz:6969/announce",
		"udp://run-2.publictracker.xyz:6969/announce",
		"udp://public.tracker.vraphim.com:6969/announce",
		"udp://public.publictracker.xyz:6969/announce",
		"udp://public-tracker.cf:6969/announce",
		"udp://opentracker.io:6969/announce",
		"udp://open.u-p.pw:6969/announce",
		"udp://open.dstud.io:6969/announce",
		"udp://new-line.net:6969/announce",
		"udp://moonburrow.club:6969/announce",
		"udp://mail.segso.net:6969/announce",
		"udp://free.publictracker.xyz:6969/announce",
		"udp://carr.codes:6969/announce",
		"udp://bt2.archive.org:6969/announce",
		"udp://6.pocketnet.app:6969/announce",
		"udp://1c.premierzal.ru:6969/announce",
		"udp://tracker.t-rb.org:6969/announce",
		"udp://tracker.srv00.com:6969/announce",
		"udp://tracker.artixlinux.org:6969/announce",
		"udp://tracker-udp.gbitt.info:80/announce",
		"udp://torrents.artixlinux.org:6969/announce",
		"udp://psyco.fr:6969/announce",
		"udp://mail.artixlinux.org:6969/announce",
		"udp://lloria.fr:6969/announce",
		"udp://fh2.cmp-gaming.com:6969/announce",
		"udp://concen.org:6969/announce",
		"udp://boysbitte.be:6969/announce",
		"udp://aegir.sexy:6969/announce"}
}

// fmt.Println("------------start----------------")
// fmt.Printf("Name:%s\n", el.DisplayPath())
// fmt.Printf("containEandS:%t\n", containEandS(el.DisplayPath(), strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(abs_season), strconv.Itoa(abs_episode)))
// fmt.Printf("containE_S:%t\n", containE_S(el.DisplayPath(), strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(abs_season), strconv.Itoa(abs_episode)))
// fmt.Printf("containsAbsoluteE:%t\n", containsAbsoluteE(el.DisplayPath(), strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(s), strconv.Itoa(e)))
// fmt.Printf("containsAbsoluteE_:%t\n", containsAbsoluteE_(el.DisplayPath(), strconv.Itoa(s), strconv.Itoa(e), abs == "true", strconv.Itoa(s), strconv.Itoa(e)))
// fmt.Println("------------end----------------")
