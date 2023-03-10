# README file for immudblog 

## Exercise

The goal is to write a simple REST or GRPC service in Golang that uses immudb as a database to store lines of logs. Service and immudb should be easily deployable using docker, docker-compose or similar. 

There should be another simple testing tool that allows you to easily create log entries. 


The service should allow to:

 - Store single log line
 
 - Store batch of log lines

 - Print history of stored logs (all, last x)

 - Print number of stored logs

 - (optional) Simple authentication mechanism to restrict read/write access.

 - (optional) Support for log buckets, so logs from different applications can be separated i.e. depending on source or some token used. 


## Design decisions

### Data Model 

As the task is storeing logs I decided to store the logs a format compatible with standard syslog format (RFC5424). The model
represents one logline in IETF syslog.

Immudb supports KV and SQL interfaces. I choosed SQL format storing the logs. 


### API 

I choosed REST API with swagger 2.0 format. The swagger test UI and documentation is generated by swaggo (https://github.com/swaggo/swag) tool. The API requires basic authentication.

### Configuration
 
The application support flags and can be overwritten with env variables.


## Structure of the project

Directories:
- cmd: contains server and test client main file
- immudb: package for persit layer. Using gorm and sql for database operations. Singleton pattern is used for connection.
- restapi: contains the HTTP hanlers, routes
- model: data models 
- docs: generated by swaggo. Base of swagger UI
- service: service layer is between restapi and persist layer. Place of the complex services (in one transaction). In this project it is not relevant as there is no complex service implemented.

Makefile is used for tooling.

Some unittests are implemented. Used github.com/stretchr/testify/mock for mock necessary interfaces. These are in mock directory.

Ready for build inside docker container as well.

## Comments for requirements

*Requirements:*
 - Store single log line 
 - Store batch of log lines

*Implementation:*

One API endpoint for store single and batch : /api/v1/logs POST


*Requirements:*
 - Print history of stored logs (all, last x)

*Implementation:*

With API endpoint /api/v1/logs GET stored logs can be queried. With query param 'count' the number of max returned logs can be set. 

*Requirements:*
 - Print number of stored logs
*Implementation:*

With API endpoint GET /api/v1/logs/count the number of stored logs can be queried.


*Requirements:*
 - (optional) Simple authentication mechanism to restrict read/write access.

*Implementation:*

basic authentication is necessary. There are 2 roles: read and write. For Adding logs write role is necessary.


*Requirements:*
 - (optional) Support for log buckets, so logs from different applications can be separated i.e. depending on source or some token used. 

*Implementation:*

As the model contains 'application' field the stored logs can be separated with it. API GET /api/v1/logs supports application filter by application query parameter. Logs from separate applcation can be queried separately.



### Tooling / building

Pre requirements:

installed go 1.18 or installed docker and docker-compose


Makefile jobs:

- Build project local
   ```sh
   make build
   ```

- Run test 
   ```sh
   make test
   ```

- Run test and open coverage in browser
   ```sh
   make testwithresults
   ```

- run and build in docker (with compose)
   ```sh
   docker-compose up --build
   ```
   It runs 2 containers: immudb and immudblog-server. 
   Exposed ports: 8080 for server and 8088 for immudb WEBUI.
   immudb:  http://127.0.0.1:8088/

   server swagger API: http://localhost:8080/swagger/index.html


## Usage

### immudblog-server

By default uses the default configuration for immudb.

```
$ bin/immudblog-server --help
Usage:
  immudblog-server [OPTIONS]

Server Options:
  -p, --port=                 Server Port. (default: 8080) [$PORT]
  -h, --host=                 Server Host (for swagger UI access). (default: localhost) [$HOST]
  -a, --auth=                 Users for simple authentication. 2 level CSV () (user1:password1:role1,user2:password2:role2) (default:
                              admin:admin:write,user:user:read) [$AUTH_USERS]

Logger Options:
      --log_type=[plain|json] Type of log output (default: plain) [$LOG_TPYE]
      --debug                 Enable debug log [$DEBUG]

Immudb connection Options:
      --dbport=               Immudb Port. (default: 3322) [$DB_PORT]
      --dbhost=               Immudb Host. (default: localhost) [$DB_HOST]
      --dbuser=               Username for Immudb. (default: immudb) [$DB_USER]
      --dbpassword=           Password for Immudb. (default: immudb) [$DB_PASSWORD]
      --dbdatabase=           Database for Immudb. (default: defaultdb) [$DB_DATABASE]

Help Options:
  -h, --help                  Show this help message

```

### immudblog-cli

Test tool for testing the server. 

```
$ bin/immudblog-cli --help
Starting
Usage:
  immudblog-cli [OPTIONS]

Server Options:
  -s, --scheme=[http|https] Server scheme. (default: http) [$SERVER_SCHEME]
  -p, --port=               Server Port. (default: 8080) [$SERVER_PORT]
  -h, --host=               Server Host (for swagger UI access). (default: localhost) [$SERVER_HOST]
  -u, --user=               Users for simple basic authentication (default: user) [$SERVER_USER]
  -w, --password=           Users password for basic authentication (default: user) [$SERVER_PASSWORD]

GetLogs Options:
      --getlogs             Use GetLogs API method. [$GETLOGS]
      --getlogs-count=      Max number of logs to get. (default: 100) [$GETLOGS_COUNT]
      --getlogs-app=        Application filter. [$GETLOGS_APP]

CountLogs Options:
      --countlogs           Use CountLogs API method. [$COUNTLOGS]

AddLogs Options:
      --addlogs             Use AddLogs API method. [$ADDLOGS]
      --addlogs-batchsize=  Batch size of AddLogs (default: 1) [$ADDLOGS_BATCHSIZE]
      --addlogs-filename=   The CSV file name to be loaded into the ImmuDB via ImmuDBLog Server. The input file should be unix file and its format is:
                            VERSION,HOSTNAME,APPLICATION,PID,PRI,TS,MESSAGEID,MESSAGE (default: input.csv) [$ADDLOGS_FILENAME]

Help Options:
  -h, --help                Show this help message
```


Adding logs can be performed with CSV file.

Sample execution for testing 3 API endpont:
```
$ bin/immudblog-cli --addlogs --addlogs-filename=samples/input.csv --user=admin --password=admin --addlogs-batchsize=3 --countlogs --getlogs --getlogs-count=3
go mod tidy
go build -o bin/immudblog-server cmd/immudblog-server/main.go
go build -o bin/immudblog-cli cmd/immudblog-cli/main.go
Starting
------- AddLogs START ------
Parsed 5 record from file: samples/input.csv
Sending with batch size:3
Using API URL: http://localhost:8080/api/v1/logs
HTTP Result: OK (200)
Result:
{"code":200,"message":"stored 3 records"}
Sending with batch size:2
Using API URL: http://localhost:8080/api/v1/logs
HTTP Result: OK (200)
Result:
{"code":200,"message":"stored 2 records"}
------- AddLogs END ------
------- GetLog START ------
Using API URL: http://localhost:8080/api/v1/logs?count=3&application=
HTTP Result: OK (200)
Result:
[{"id":53,"hostname":"host1","application":"app1","pid":"1234","pri":1,"timestamp":"2023-02-28T22:49:34.013Z","messageid":1111,"meaasge":"This is a message5"},{"id":52,"hostname":"host1","application":"app1","pid":"1234","pri":1,"timestamp":"2023-02-28T22:48:34.013Z","messageid":1111,"meaasge":"This is a message4"},{"id":51,"hostname":"host1","application":"app1","pid":"1234","pri":1,"timestamp":"2023-02-28T22:47:34.013Z","messageid":1111,"meaasge":"This is a message3"}]
------- GetLog END ------
------- CountLogs START ------
Using API URL: http://localhost:8080/api/v1/logs/count
HTTP Result: OK (200)
Result:
{"count":53}
------- CountLogs END ------
```

