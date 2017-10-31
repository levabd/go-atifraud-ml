# go-atifraud-ml

### Apps
Location: ```lib/go/apps```

### Project structure
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

### For the dev process run in the project root
```
go run [running app in lib/go/apps/*]
```

### To build app run in the project root
```
go build [running app in lib/go/apps/*]
```

### To test app run in the project root
```
go test -v ./... 
```

### Integration   
1 Insert content of the method lib/go/apps/client/client.go - main to needed place in you code   
2 Configure URL adress in the code   
3 Configure following lines in you code:   
```
req.Header.Set("Host", "localhost")
req.SetRequestURI("http://localhost:8082")
req.Header.Set("Body-Header-Ip", "62.84.44.222")  -- header ip
req.SetBodyString(`{"Cache-Control":"no-cache","Connection":"Keep-Alive","Pragma":"no-cache","Accept":"*\/*","Accept-Encoding":"gzip, deflate","From":"bingbot(at)microsoft.com","Host":"www.vypekajem.com","User-Agent":"Mozilla\/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit\/537.51.1 (KHTML, like Gecko) Version\/7.0 Mobile\/11A465 Safari\/9537.53 (compatible; bingbot\/2.0; +http:\/\/www.bing.com\/bingbot.htm)"}`) - header json as string
```
   
### Tip
If yoy want to create new feature - just add new file lib/go/apps/[app]/[app].go
and create main function in it. 

### Preparation before running  
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
```
This will create prediction_server bin file in the <b>lib/python</b> directory 

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
@daily cd /home/vmuser/gocode/src/github.com/levabd/go-atifraud-ml && go run lib/go/apps/parser/parse_gz_logs.go > /home/vmuser/gocode/src/github.com/levabd/go-atifraud-ml/data/logs_go/parse_gz_logs.log 2>&1
```
Save changes

Start <b>antifraud_server</b> and <b>prediction_server</b> as deamon
```
nohup ./antifraud_server > data/logs_go/antifraud_server.log 2>&1 &
nohup ./lib/python/prediction_server > data/logs_go/prediction_server.log 2>&1 &
```

