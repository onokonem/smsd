package cfg

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

type Esme struct {
	Enabled         bool
	Name            string
	Operator        string
	Host            string
	Port            int
	EnquireInterval int    `toml:"enquire_interval"`
	SystemId        string `toml:"system_id"`
	Password        string
	DstNumberPrefix int    `toml:"dst_number_prefix"`
	SrcAddr         string `toml:"src_addr"`
	AddrTon         int    `toml:"addr_ton"`
	AddrNpi         int    `toml:"addr_npi"`
	SrcTon          int    `toml:"src_ton"`
	SrcNpi          int    `toml:"src_npi"`
	TzShift         int    `toml:"tz_shift"`
}

type Database struct {
	User       string
	Password   string
	Host       string
	Name       string
	StorableDb string `toml:"storable_db"`
}

type Cfg struct {
	Host  string
	Port  string
	Db    Database `toml:"database"`
	Esmes map[string]Esme
}

func New(path *string, logger *zap.Logger) Cfg {
	var c Cfg
	if _, err := toml.DecodeFile(*path, &c); err != nil {
		panic(err.Error())
	}
	if c.Db.User == "" || c.Db.Password == "" || c.Db.Host == "" || c.Db.StorableDb == "" || c.Port == "" || c.Host == "" {
		panic(fmt.Errorf("missing mandatory config parameters"))
	}
	for name, options := range c.Esmes {
		if options.Host == "" || options.Port == 0 || options.SystemId == "" || options.Password == "" || options.SrcAddr == "" {
			logger.Warn("Missing one or many parameters for oper", zap.String("name", name), zap.Any("time", time.Now().Format(time.RFC3339)))
			delete(c.Esmes, name)
			continue
		}
	}
	return c
}
