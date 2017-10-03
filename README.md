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
           * client server
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

# Tip
If yoy want to create new feature - just add new file lib/go/apps/[app]/[app].go
and create main function in it. 