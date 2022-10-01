package configs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml"
)

func ReadConfig(configPath string) Config {
	configFile := fmt.Sprintf("%s/config.toml", configPath)
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("failed to read config: %v\n", err)
	}

	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to parse config: %v\n", err)
	}
	return cfg
}

func (db *Database) GetDSN() string {
	var buf bytes.Buffer

	// [username[:password]@]
	if len(db.User) > 0 {
		buf.WriteString(db.User)
		if len(db.Password) > 0 {
			buf.WriteByte(':')
			buf.WriteString(db.Password)
		}
		buf.WriteByte('@')
	}

	// [(address)]
	if len(db.Host) > 0 {
		buf.WriteByte('(')
		buf.WriteString(db.Host)
		if len(db.Port) > 0 {
			buf.WriteByte(':')
			buf.WriteString(db.Port)
		}
		buf.WriteByte(')')
	}

	// /dbname
	buf.WriteByte('/')
	buf.WriteString(db.Name)
	return buf.String()
}
