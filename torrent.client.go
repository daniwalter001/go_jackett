package main

import "github.com/anacrolix/torrent"

func TorrentClient() *torrent.Client {
	config := torrent.NewDefaultClientConfig()
	config.DataDir = "./temp"
	config.ListenHost = func(network string) string { return "localhost" }
	client, err := torrent.NewClient(config)

	if err != nil {
		_c, _ := torrent.NewClient(nil)
		return _c
	}

	return client
}
