# go-atifraud-ml

## Apps
Location: ```lib/go/apps```

# Project structure
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

# Integration
1 Start server ./server   
2 Insert content of the method lib/go/apps/client/client.go - main to needed place in you code   
3 Configure URL adress in the code   
4 Configure following lines in you code:   
```
req.Header.Set("Host", "localhost")
req.SetRequestURI("http://localhost:8082")
req.Header.Set("Body-Header-Ip", "62.84.44.222")  -- header ip
req.SetBodyString(`{"Cache-Control":"no-cache","Connection":"Keep-Alive","Pragma":"no-cache","Accept":"*\/*","Accept-Encoding":"gzip, deflate","From":"bingbot(at)microsoft.com","Host":"www.vypekajem.com","User-Agent":"Mozilla\/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X) AppleWebKit\/537.51.1 (KHTML, like Gecko) Version\/7.0 Mobile\/11A465 Safari\/9537.53 (compatible; bingbot\/2.0; +http:\/\/www.bing.com\/bingbot.htm)"}`) - header json as string
```
   
# Tip
If yoy want to create new feature - just add new file lib/go/apps/[app]/[app].go
and create main function in it. 


