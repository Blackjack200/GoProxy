package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type config struct {
	Bind         string `json:"local"`
	Remote       string `json:"remote"`
	LocalStatus  bool   `json:"local-status"`
	LocalMotd    string `json:"local-motd"`
	ProxySideXBL bool   `json:"proxy-side-xbl"`
	RemoteXBL    bool   `json:"remote-xbl"`
	RemoteStatus string `json:"remote-status-addr"`
	ResourcePack bool   `json:"bypass-resource-pack"`
	SafeConnect  bool   `json:"safe-connect"`
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
	Cfg.LocalMotd = "GoProxy"
	Cfg.LocalStatus = true
	Cfg.RemoteStatus = Cfg.Remote
	Cfg.ProxySideXBL = true
	Cfg.RemoteXBL = false
	Cfg.SafeConnect = true
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

func RemoteStatus() string {
	return Cfg.RemoteStatus
}

func ProxySideXBL() bool {
	return Cfg.ProxySideXBL
}

func RemoteXBL() bool {
	return Cfg.RemoteXBL
}

func LocalStatus() bool {
	return Cfg.LocalStatus
}

func SafeConnect() bool {
	return Cfg.SafeConnect
}

func BypassResourcePack() bool {
	return Cfg.ResourcePack
}

func Motd() string {
	return Cfg.LocalMotd
}
