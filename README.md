### Simple IMDB Web Scraper

> Built using `go1.19.2`

#### What is does ?

This program scrapes `www.imdb.com/name` for celebreties of a given birthday where you provide a Date and Month and stores all data in a file `output.json` and this program stores cache in `.imdb_cache` and uses from here.

#### How it does ?

This Simple IMDB Web Scaper uses [`colly`](https://github.com/gocolly/colly) ans `os` package to perform the scraping and writing the json into a file.
