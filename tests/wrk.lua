wrk.method = "POST"
wrk.body   = "Host: servicer.mgid.com \nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36\nUpgrade-Insecure-Requests: 1 \nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8\nDNT: 1 \n Accept-Encoding: gzip, deflate\nAccept-Language: ru-UA,ru;q=0.9,en-US;q=0.8,en;q=0.7,ru-RU;q=0.6\nCache-Control: no-cache"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
