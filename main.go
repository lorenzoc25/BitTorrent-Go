package main

import (
	"log"
	"os"

	"github.com/lorenzoc25/bittorrent-go/torrentfile"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: bittorent-go torrentfile <output file name>")
	}
	inPath := os.Args[1]
	outPath := os.Args[2]
	tf, err := torrentfile.Open(inPath)
	if err != nil {
		log.Fatal(err)
	}
	err = tf.DownloadToFile(outPath)
	if err != nil {
		log.Fatal(err)
	}
}
