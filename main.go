package main

import (
	"log"
	"os"

	"github.com/lorenzoc25/bittorrent-go/httpdownload"
	"github.com/lorenzoc25/bittorrent-go/torrentfile"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Usage: bittorent-go <mode> <source> <output file name>")
	}
	mode := os.Args[1]
	if mode == "torrent" {
		inPath := os.Args[2]
		outPath := os.Args[3]
		tf, err := torrentfile.Open(inPath)
		if err != nil {
			log.Fatal(err)
		}
		err = tf.DownloadToFile(outPath)
		if err != nil {
			log.Fatal(err)
		}
	} else if mode == "http" {
		url := os.Args[2]
		outPath := os.Args[3]
		err := httpdownload.HTTPDownload(url, outPath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("Unknown mode")
	}
}
