package main

import "github.com/levabd/go-atifraud-ml/lib/go/services"

func main() {
	println(services.IsCrawler("62.84.44.222", `Mozilla\/5.0(compatible;Googlebot\/2.1;+http:\/\/www.google.com\/bot.html)`))
}
