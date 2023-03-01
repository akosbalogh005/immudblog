package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"immudblog/model"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
)

var ServerFlags = struct {
	Scheme   string `long:"scheme" short:"s" description:"Server scheme." choice:"http" choice:"https" default:"http" env:"SERVER_SCHEME"`
	Port     int    `long:"port" short:"p" description:"Server Port." default:"8080" env:"SERVER_PORT"`
	Host     string `long:"host" short:"h" description:"Server Host (for swagger UI access)." default:"localhost" env:"SERVER_HOST"`
	User     string `long:"user" short:"u" description:"Users for simple basic authentication" default:"user" env:"SERVER_USER"`
	Password string `long:"password" short:"w" description:"Users password for basic authentication" default:"user" env:"SERVER_PASSWORD"`
}{}

var GetLogsFlags = struct {
	IsGetLogs   bool   `long:"getlogs" description:"Use GetLogs API method." env:"GETLOGS"`
	Count       int    `long:"getlogs-count" description:"Max number of logs to get." default:"100" env:"GETLOGS_COUNT"`
	Application string `long:"getlogs-app" description:"Application filter." env:"GETLOGS_APP"`
}{}

var CountLogsFlags = struct {
	IsCountLogs bool `long:"countlogs" description:"Use CountLogs API method." env:"COUNTLOGS"`
}{}

var AddLogsFlags = struct {
	IsAddLogs bool   `long:"addlogs" description:"Use AddLogs API method." env:"ADDLOGS"`
	BatchSize int    `long:"addlogs-batchsize" description:"Batch size of AddLogs" env:"ADDLOGS_BATCHSIZE" default:"1"`
	FileName  string `long:"addlogs-filename" description:"The CSV file name to be loaded into the ImmuDB via ImmuDBLog Server. The input file should be unix file and its format is: VERSION,HOSTNAME,APPLICATION,PID,PRI,TS,MESSAGEID,MESSAGE" env:"ADDLOGS_FILENAME" default:"input.csv"`
}{}

func InitFlags() {
	parser := flags.NewParser(nil, flags.Default)
	parser.AddGroup("Server Options", "Server Options", &ServerFlags)
	parser.AddGroup("GetLogs Options", "GetLogs Options", &GetLogsFlags)
	parser.AddGroup("CountLogs Options", "GetLogs Options", &CountLogsFlags)
	parser.AddGroup("AddLogs Options", "AddLogs Options", &AddLogsFlags)

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}
}

// Immudb Logstore Server Client
func main() {
	fmt.Println("Starting")

	InitFlags()

	if AddLogsFlags.IsAddLogs {
		RunAddLogs()
	}

	if GetLogsFlags.IsGetLogs {
		RunGetLog()
	}

	if CountLogsFlags.IsCountLogs {
		RunCountLogs()
	}

}

func getUrl(path string, queryParam string) string {
	ret := fmt.Sprintf("%s://%s:%d/api/v1/%s", ServerFlags.Scheme, ServerFlags.Host, ServerFlags.Port, path)
	if queryParam != "" {
		ret = fmt.Sprintf("%s?%s", ret, queryParam)
	}
	fmt.Printf("Using API URL: %s\n", ret)
	return ret
}

func doAndPrintResult(req *http.Request) {
	h := http.Client{}
	req.SetBasicAuth(ServerFlags.User, ServerFlags.Password)
	resp, err := h.Do(req)
	fmt.Printf("HTTP Result: %v (%d)\n", http.StatusText(resp.StatusCode), resp.StatusCode)
	if err != nil {
		fmt.Printf("Error (HTTP): %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error (reading): %v\n", err)
		return
	}
	fmt.Println("Result:")
	fmt.Printf("%v\n", string(body))

}

func RunGetLog() {
	fmt.Println("------- GetLog START ------")
	queryParam := fmt.Sprintf("count=%d&application=%s", GetLogsFlags.Count, GetLogsFlags.Application)
	req, _ := http.NewRequest("GET", getUrl("logs", queryParam), nil)
	doAndPrintResult(req)
	fmt.Println("------- GetLog END ------")

}

func RunCountLogs() {
	fmt.Println("------- CountLogs START ------")
	req, _ := http.NewRequest("GET", getUrl("logs/count", ""), nil)
	doAndPrintResult(req)
	fmt.Println("------- CountLogs END ------")
}

func readCSVFile(fn string) (ret []model.Log, err error) {
	f, err := os.Open(AddLogsFlags.FileName)
	if err != nil {
		return
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		fullLine := fileScanner.Text()
		line := strings.SplitN(fullLine, ",", 8)
		if len(line) != 8 {
			fmt.Printf("Error. Not enough values. Skipping line:'%v'", fullLine)
			continue
		}
		ver, err := strconv.ParseInt(line[0], 10, 0)
		if err != nil {
			fmt.Printf("Error: %v, Skipping line:'%v'", err, line)
			continue
		}
		pri, err := strconv.ParseInt(line[4], 10, 0)
		if err != nil {
			fmt.Printf("Error: %v, Skipping line:'%v'", err, line)
			continue
		}
		t, err := time.Parse(time.RFC3339, line[5])
		if err != nil {
			fmt.Printf("Error: %v, Skipping line:'%v'", err, line)
			continue
		}
		mid, err := strconv.ParseInt(line[6], 10, 0)
		if err != nil {
			fmt.Printf("Error: %v, Skipping line:'%v'", err, line)
			continue
		}
		l := model.Log{Version: int32(ver), Hostname: line[1], Application: line[2], Pid: line[3], Pri: int32(pri), Timestamp: t, Messageid: mid, Message: line[7]}
		ret = append(ret, l)
	}
	f.Close()
	return
}

func RunAddLogs() {
	fmt.Println("------- AddLogs START ------")
	ret, _ := readCSVFile(AddLogsFlags.FileName)
	fmt.Printf("Parsed %d record from file: %v\n", len(ret), AddLogsFlags.FileName)

	var toSend []model.Log
	for i := 0; i < len(ret); i++ {
		toSend = append(toSend, ret[i])

		if len(toSend) >= AddLogsFlags.BatchSize || i == len(ret)-1 {
			fmt.Printf("Sending with batch size:%d\n", len(toSend))
			body, err := json.Marshal(toSend)
			if err != nil {
				fmt.Printf("Error (JSON): %v\n", err)
			} else {
				bodyReader := bytes.NewReader(body)
				req, _ := http.NewRequest("POST", getUrl("logs", ""), bodyReader)
				doAndPrintResult(req)
			}
			toSend = toSend[:0]
		}
	}
	fmt.Println("------- AddLogs END ------")

}
