from flask import Flask, send_file
import requests
import cfscrape
from flask import request
import urllib.parse
import re
from io import BytesIO

app = Flask(__name__)


def get_scraper(use_cfs):
    if use_cfs == "1":
        return cfscrape.create_scraper()
    else:
        return requests


def get_torrent_from_demonoid(id, use_cfs):
    scraper = get_scraper(use_cfs)
    content = scraper.get("https://www.demonoid.pw/genlb.php?genid=" + id).content.decode("utf-8")
    return scraper.get(re.findall("https?://(?:www.)hypercache.pw/metadata/[a-zA-Z0-9/\\?]+", content)[0]).content


def get_extractor(extractor, url, use_cfs):
    if extractor == "demonoid":
        return get_torrent_from_demonoid(url, use_cfs)
#     TODO: throw error if no extractor is found


@app.route('/')
def show_item():
    url = request.args.get("url", type=str)
    url = urllib.parse.unquote(url)
    use_cfs = request.args.get('use_cfs', type=str)
    use_regex = request.args.get('use_regex', type=str)

    scraper = get_scraper(use_cfs)

    content = scraper.get(url).content
    content = content.decode("utf-8")
    if use_regex == "1":
        tr = urllib.parse.unquote(request.args.get('to_replace', type=str))
        replacement = urllib.parse.unquote(request.args.get('replacement', type=str))
        content = content.replace(tr, replacement)
    return content


@app.route('/extractor/<site>/<url_id>/')
def extract(site, url_id):
    use_cfs = "0"
    return send_file(BytesIO(get_extractor(site, url_id, use_cfs)))


if __name__ == "__main__":
    app.run()

