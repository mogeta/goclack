package main

import (
	"io/ioutil"

	"github.com/mogeta/goclack"
	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yaml"
)

type config struct {
	Token string
}

func main() {

	//Load token from yaml
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return
	}

	var m config
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		panic(err)
	}

	//run bot
	bot := goclack.New(m.Token)
	bot.Run()
}
