package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Bind         string `json:"bind"`
	Remote       string `json:"remote"`
	ProxyXbox    bool   `json:"proxy-xbox-auth"`
	RemoteXbox   bool   `json:"xbox-auth"`
	Motd         string `json:"motd"`
	ResourcePack bool   `json:"bypass-resource-pack"`
}

var (
	Cfg        *config
	File       string
	WorkingDir string
)

func Initialize() error {
	WorkingDir, _ = os.Getwd()
	WorkingDir, _ = filepath.Abs(WorkingDir)
	File = WorkingDir + "/config.json"

	Cfg = &config{}
	Cfg.Bind = "0.0.0.0:19132"
	Cfg.Remote = "127.0.0.1:19133"
	Cfg.Motd = "GoProxy"
	Cfg.ProxyXbox = true
	Cfg.RemoteXbox = false
	Cfg.ResourcePack = true

	if _, s := os.Stat(File); os.IsNotExist(s) {
		bytes, err := json.MarshalIndent(Cfg, "", "	")
		if err != nil {
			return err
		}
		_ = ioutil.WriteFile(File, bytes, 0777)
	} else {
		con, _ := ioutil.ReadFile(File)
		err := json.Unmarshal(con, Cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func Bind() string {
	return Cfg.Bind
}

func Remote() string {
	return Cfg.Remote
}

func XBL() bool {
	return Cfg.ProxyXbox
}

func RemoteXBL() bool {
	return Cfg.RemoteXbox
}

func BypassResourcePack() bool {
	return Cfg.ResourcePack
}

func Motd() string {
	return Cfg.Motd
}
