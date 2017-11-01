# go-atifraud-ml

> Bot recognizer

## Table of Contents

- [Apps](#apps)
- [Project structure](#project-structure)
- [Integration](#integration)
- [Education process](#education-process)
- [Configure and Running](#configure-and-running)
- [Benchmark](#benchmark)
  * [Apache Bench](#apache-bench)
  * [wrk](#wrk)
- [Tips](#tips)

## Apps
```
lib/go/apps
```

## Project structure
```
* data - App data
    * db - Udget db 
    * logs - Header logs
    * logs_go - Go app logs. See services/logging.go
    * unit_tests_files - Logs for unit testing. 
* lib
    * go 
       * apps - Applications 
           * parser
           * server
           * client
       * helpers - Helpers
       * services - Logger, parser, paid_generator, and go udger methods implementation
       * models - Log, GzLog models (Log model contain methods for trimming value data and order data fields) 
       * .env
* tests - some files for benchmark tools
* initiate_prediction_cluster.sh - bash script for running prediction servers (current ports - 5000-5007)
```

## Integration   
1 Insert content of the method lib/go/apps/client/client.go - main to needed place in you code   
2 Configure URL address in the code   
3 Configure following lines in you code:   
```
req.Header.Set("Host", "localhost")
req.SetRequestURI("http://localhost:8082")
req.Header.Set("X-Real-IP", "62.84.44.222") - client ip(will by check in udger DB0
req.SetBodyString(``Host: servicer.mgid.com
                   User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.75 Safari/537.36
                   Upgrade-Insecure-Requests: 1
                   Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8
                   DNT: 1
                   Accept-Encoding: gzip, deflate
                   Accept-Language: ru-UA,ru;q=0.9,en-US;q=0.8,en;q=0.7,ru-RU;q=0.6
                   Cache-Control: no-cache``) - header json as string
```
   

## Education process

To start education run in root of the project
```
go run lib/go/apps/parser/parse_gz_logs.go -sample 80000 -educate true
```
This will:
 - prepare data for the train
 - train
 - send requests to all the prediction servers to reload model (server addresses are hardcoded at the moment)
  
There are another command flags in <b>lib/go/apps/parser/parse_gz_logs.go</b> look briefly this file.

## Configure and Running 
Set postgresql login and password in files
* lib/go/.env
* lib/python/train.pyx

Build antifraud server by running
```
go build -o antifraud_server lib/go/apps/server/server.go 
```
This will create antifraud_server bin file in the root of the project

Build prediction server by running
```
cd lib/python
cython --embed -o server.c server.pyx && gcc -Os -I /usr/include/python3.6m -o prediction_server server.c -lpython3.6m -lpthread -lm -lutil -ldl
cd - 
```
This will create prediction_server bin file in the <b>lib/python</b> directory

Start prediction servers on ports 5000-5007
```
./initiate_prediction_cluster.sh
```
This will start 8 (for CPU) prediction server instances (Load balancing mechanism between this servers is implemented in lib/go/apps/server/server.go) 

Make log files
```
touch data/logs_go/antifraud_server.log
touch data/logs_go/prediction_server.log
touch data/logs_go/parse_gz_logs.log.log
```

Configure <b>crontab</b> for scheduling education
Open crontab: 
```
crontab -e
```
Add following content:
```
@daily cd /home/vmuser/gocode/src/github.com/levabd/go-atifraud-ml && go run lib/go/apps/parser/parse_gz_logs.go -sample 80000 -educate true > /home/vmuser/gocode/src/github.com/levabd/go-atifraud-ml/data/logs_go/parse_gz_logs.log 2>&1
```
Save changes

Start <b>antifraud_server</b> as deamon
```
nohup ./antifraud_server -predictionServers 0.0.0.0:5000,0.0.0.0:5001,0.0.0.0:5002,0.0.0.0:5003,0.0.0.0:5004,0.0.0.0:5005,0.0.0.0:5006,0.0.0.0:5007 > data/logs_go/antifraud_server.log 2>&1 &
```

## Benchmark

### [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html)
```
ab -p tests/postbody.txt -H "X-Real-IP: 62.84.44.222"  -c 100 -n 2000 -k http://127.0.0.1:8082/
This is ApacheBench, Version 2.3 <$Revision: 1706008 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking 127.0.0.1 (be patient)
Completed 200 requests
Completed 400 requests
Completed 600 requests
Completed 800 requests
Completed 1000 requests
Completed 1200 requests
Completed 1400 requests
Completed 1600 requests
Completed 1800 requests
Completed 2000 requests
Finished 2000 requests


Server Software:        fasthttp
Server Hostname:        127.0.0.1
Server Port:            8082

Document Path:          /
Document Length:        5 bytes

Concurrency Level:      100
Time taken for tests:   2.215 seconds
Complete requests:      2000
Failed requests:        50
   (Connect: 0, Receive: 0, Length: 50, Exceptions: 0)
Non-2xx responses:      47
Keep-Alive requests:    2000
Total transferred:      330976 bytes
Total body sent:        1172000
HTML transferred:       14224 bytes
Requests per second:    902.96 [#/sec] (mean)
Time per request:       110.747 [ms] (mean)
Time per request:       1.107 [ms] (mean, across all concurrent requests)
Transfer rate:          145.93 [Kbytes/sec] received
                        516.73 kb/s sent
                        662.66 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       2
Processing:     2  108  53.0    103     286
Waiting:        2  108  53.0    103     286
Total:          2  108  53.0    103     286

Percentage of the requests served within a certain time (ms)
  50%    103
  66%    130
  75%    144
  80%    154
  90%    179
  95%    200
  98%    230
  99%    245
 100%    286 (longest request)

```

### [wrk](https://github.com/wg/wrk)
```
wrk -t12 -c400 -d30s -s tests/wrk.lua ttp://localhost:8082/

----
Running 30s test @ ttp://localhost:8082/
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   377.50ms  241.04ms   1.86s    72.97%
    Req/Sec    80.10     38.68   282.00     62.61%
  28643 requests in 30.10s, 4.48MB read
  Non-2xx or 3xx responses: 405
Requests/sec:    951.65
Transfer/sec:    152.55KB
```

## Tips
If yoy want to create new feature - just add new file lib/go/apps/[app]/[app].go
and create main function in it. 
