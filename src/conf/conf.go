package conf

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type MgoCfg struct {
	Server    string `xml:"server"`
	Port      string `xml:"port"`
	DB        string `xml:"db"`
	PageLimit int    `xml:"page_limit"`
}

type RdsCfg struct {
	Server    string `xml:"server"`
	Port      string `xml:"port"`
	MaxActive int    `xml:"max_active"`
	MaxIdle   int    `xml:"max_idle"`
}

type ServerCfg struct {
	Host string `xml:"host"`
	Port string `xml:"port"`
}

type Config struct {
	XMLName            xml.Name  `xml:"config"`
	MgoCfg             MgoCfg    `xml:"mongo"`
	RdsCfg             RdsCfg    `xml:"redis"`
	ServerCfg          ServerCfg `xml:"server"`
	LogCfgName         string    `xml:"log_cfg_name"`
	TokenLastTime      uint32    `xml:"token_last_time"`
	AdminUser          AdminUser `xml:"admin_user"`
	VideoCheckUrl      string    `xml:"video_check_url"`
	SensitiveWordsFile string    `xml:"sensitive_words_file"`
}

type AdminUser struct {
	AccID int64 `xml:"accid"`
}

var (
	Cfg Config
)

func GetCfg() Config {
	return Cfg
}

func init() {
	content, err := ioutil.ReadFile("conf/conf.xml")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	err = xml.Unmarshal(content, &Cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
