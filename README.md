# Torrent site middleware

A flask app to help deal with and make useful torrent site RSS feeds

## Usage

Build with `go build middleware.go`

the app takes 4 get params 

* url: The URL encode url of the RSS feed

* use_regex: to replace certain lines in the rss feed (Set to 1 for true)

* to_replace: The string to be replaced

* replacement: The string to replace to_replace

Example of a full request http://127.0.0.1:5000/?url=https%3A%2F%2Fworldwidetorrents.unblocked.vet%2Frss.php%3Fcat%3D132%26dllink%3D1&use_regex=1&to_replace=https%3A%2F%2Fworldwidetorrents.me&replacement=https%3A%2F%2Fworldwidetorrents.unblocked.vet

This example will take the rss feed of the proxy site worldwidetorrents.unblocked.vet and point the torrent download links to worldwidetorrents.unblocked.vet (They would normally be pointed at worldwidetorrents.me)

## Extractors 

Some site rss feed don't include the download link/magnet links but instead links to the download page. 
To download from these feeds extractors are used. 

Extractors are functions that download a webpage, extract a torrent file and return it

The format for extract urls is /extractor/<site>/<url_id>/ where <site> is the site name and <url_id> is
a string that the extract can use to get the url of the torrent file

Extractors currenlt exist for the following sites

* 1337x.to
* demonoid.pw

Example extractor url: 127.0.0.1:5000/extractor/demonoid/000000000000000

## Scrapers

Some sites don't provide rss feeds at all, this is where scrapers come in.

Scrapers take a page from a site and turn it into a rss feed

The format for extract urls is /rss/<site>?url=<url> where <site> is the name of the site (without TLD) and <url> is the url encoded url

Extractors currenlt exist for the following sites

* 1337x.to

Example scraper url: 127.0.0.1:5000/rss/1337x?url=https%3A%2F%2F1337x.to%2Fuser%USERNAME%2F
