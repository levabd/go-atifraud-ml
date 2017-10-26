package main

import (
	"github.com/valyala/fasthttp"
	"time"
	"log"
)
func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer timeTrack(time.Now(), "main")

	c := &fasthttp.Client{}
	c.MaxIdleConnDuration =  100 * time.Second

	var req fasthttp.Request
	req.Header.SetMethod("POST")
	req.SetRequestURI("http://localhost:8082")
	req.Header.Set("Host", "localhost")
	req.Header.Set("Body-Header-Ip", "91.203.170.233")
	req.SetBodyString(`{"Host": "www.rosimperija.info", "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8", "Referer": "https://yandex.ru/clck/jsredir?from=yandex.ru%3Bsearch%2F%3Bweb%3B%3B&text=&etext=1572.l9A_fF9N3s6zTLPv7_VlKAyxSRpofP75td_9hB_tUK12A75Ihj4ySfZC4ckGerzh-ISFNKD6DoCgJwPCmywLbpPRANfsatzdOF-QR_2xcis.e37960a9453d8a6d820b41077410c065a9d972d1&uuid=&state=PEtFfuTeVD5kpHnK9lio9aJ2gf1Q1OEQHP1rbfzHEMvZEAs4QuMnSA,,&&cst=AiuY0DBWFJ5Hyx_fyvalFLfT1J2GZbwEDz_Aa37CIawAt5D7-Sz_SujPk_wpBGJzMLydPHAySPGyR8i7ousOJOtndL1We7iIPbx070iuipbxxwnlNsdwsSMXOvlQinM1WSMtMa-Pj29m-_JzKLergVl8toV_EBawOj7HiGmGaoc4r-3sushQa_itF1VUaskitltb2Pf2Lh6FMEqn5RdGy_h9gB5D3HJAmJShH57MBl1Eri_wrvoBVGMhDknvgOas1w2kQ3G-S6D7Dx1hB1mr4qyTLq8eRzFUdFLA97gIqmEMEj3zhzDyX253_PPZpWIQKmj7s5kGSS3X0yuFiuE4wg,,&data=UlNrNmk5WktYejR0eWJFYk1LdmtxblFwR0NVTC1fWWV6dTU3bjkyY3dXV3R4WTNHdEUtcDlCQ1hLSkJtRnZIVVBkbkJSeWExLXdKSi1PZGo1dmt5UmZWcEZTdi0zSjljd3lyejgtc0tTUWh1VUNheWNxd2c1ZVJxZHhxd3ZOUzc,&sign=c7d606e97567c07a833100be59e1e0a8&keyno=0&b64e=2&ref=orjY4mGPRjk5boDnW0uvlrrd71vZw9kpxv4OtDXjGxXQ74AnQVptYk8pPF0X0V6Xm-mT5cCQVJ9cq1VOuOpCWIJn1H4U9a-3O-I6kmDsa6BfsjJ3Yz-Ncdt9oMlxDkvsoXUyNyVtDImRpYB63aJKvAd3lhgEu0RmgzFMMXOMKMm94lRrXi31Lf6vUsVH39U8Qbrq32wBcmzr6f7L3uu2iL7c66S_jyAXuqu7nmR6nKxZMl4t23p6xuumZ2NWxFqoRmvDY0CUMjw2qkt0mgWoeP4MpkZkr3VqlXQXFOSifiD0TKry6kkTxg,,&l10n=ru&cts=1507840878165&mc=4.7750726752065225&bu=uniq150784397399974980", "Connection": "keep-alive", "User-Agent": "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 YaBrowser/17.9.1.768 Yowser/2.5 Safari/537.36", "Accept-Encoding": "gzip, deflate", "Accept-Language": "ru,en;q=0.8", "Upgrade-Insecure-Requests": "1"}`)
	var resp fasthttp.Response

	err := c.DoTimeout(&req, &resp, time.Second)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode() == 200 { // OK
		useResponseBody(resp.Body())
	} else{
		println(resp.StatusCode())
		println(resp.Body())
	}
}

func useResponseBody(body []byte) {
	log.Println("resp.Body", string(body))
}