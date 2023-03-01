package config

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

const (
	LOG_TYPE_PLAIN = "plain"
	LOG_TYPE_JSON  = "json"
)

var ServerFlags = struct {
	Port      int    `long:"port" short:"p" description:"Server Port." default:"8080" env:"PORT"`
	Host      string `long:"host" short:"h" description:"Server Host (for swagger UI access)." default:"localhost" env:"HOST"`
	AuthUsers string `long:"auth" short:"a" description:"Users for simple authentication. 2 level CSV () (user1:password1:role1,user2:password2:role2)" default:"admin:admin:write,user:user:read" env:"AUTH_USERS"`
}{}

var ImmudbFlags = struct {
	Port     int    `long:"dbport" description:"Immudb Port." default:"3322" env:"DB_PORT"`
	Host     string `long:"dbhost" description:"Immudb Host." default:"localhost" env:"DB_HOST"`
	Username string `long:"dbuser" description:"Username for Immudb." default:"immudb" env:"DB_USER"`
	Password string `long:"dbpassword" description:"Password for Immudb." default:"immudb" env:"DB_PASSWORD"`
	Database string `long:"dbdatabase" description:"Database for Immudb." default:"defaultdb" env:"DB_DATABASE"`
}{}

var LogFlags = struct {
	Type  string `long:"log_type" description:"Type of log output" env:"LOG_TYPE" choice:"plain" choice:"json" default:"plain"`
	Debug bool   `long:"debug" description:"Enable debug log" env:"DEBUG"`
}{}

func Init() {

	parser := flags.NewParser(nil, flags.Default)
	parser.AddGroup("Server Options", "Server Options", &ServerFlags)
	parser.AddGroup("Logger Options", "Logger Options", &LogFlags)
	parser.AddGroup("Immudb connection Options", "Immudb connection Options", &ImmudbFlags)

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
