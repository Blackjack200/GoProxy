package config

import (
	"encoding/json"
	"github.com/sandertv/gophertunnel/minecraft/auth"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"time"
)

var (
	Token     *oauth2.Token
	TokenSrc  oauth2.TokenSource
	TokenFile string
)

type jsonToken struct {
	Access  string `json:"access_token"`
	Type    string `json:"token_type"`
	Refresh string `json:"refresh_token"`
}

func CacheTokenNotExists() bool {
	_, s := os.Stat(TokenFile)
	return os.IsNotExist(s)
}

func InitializeToken() error {
	TokenFile = WorkingDir + "/token.json"

	if CacheTokenNotExists() {
		logrus.Info("Generate new Token")
		var err error
		Token, err = auth.RequestLiveToken()
		if err != nil {
			panic(err)
		}
		_ = WriteToken(Token)
	} else {
		//TODO improve idiot code
		con, _ := ioutil.ReadFile(TokenFile)
		data := &jsonToken{}
		err := json.Unmarshal(con, data)
		Token = &oauth2.Token{}

		Token.AccessToken = data.Access
		Token.RefreshToken = data.Refresh
		Token.TokenType = data.Type
		Token.Expiry = time.Now().AddDate(100, 0, 0)

		TokenSrc = oauth2.StaticTokenSource(Token)
		logrus.Info("Use cached Token")
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteToken(token *oauth2.Token) error {
	bytes, err := json.MarshalIndent(*token, "", "	")
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile(TokenFile, bytes, 0777)
	return nil
}
