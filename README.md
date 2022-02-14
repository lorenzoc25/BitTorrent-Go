# BitTorrent-Go
This is a downloader implemented by Go that supports multiple protocols. Currently, it only supports BitTorrent and Http URL downloads
## Build
After cloning the repo, do 
```
go build
```
then the executable `bittorrent-go` will be ready to use
## Usage
To download a file using this program, do 
```
./bittorent-go <mode> <source> <outpath>
```
where `<mode>` currently can be `torrent` or `http`
- `http`: using http url to download the file, `<source>` should then be a URL
- `torrent`: using a torrent file to create p2p protocal to download the file, `<source>` should be a path for a torrent file 

`<outpath>`: the output path and name of the downloaded file.
