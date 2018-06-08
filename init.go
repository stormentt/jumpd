package main

import (
	"git.secrecy.rocks/stormentt/jumpd/groups"
	"git.secrecy.rocks/stormentt/jumpd/servers"
	"git.secrecy.rocks/stormentt/jumpd/users"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	initLogging()

	loadServers()
	loadGroups()
	loadUsers()

	servers.Defaults()
	groups.Defaults()
	users.Defaults()
}

func loadServers() {
	err := viper.UnmarshalKey("servers", &servers.Servers)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to load servers from config")
	}
}

func loadGroups() {
	err := viper.UnmarshalKey("groups", &groups.Groups)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to load groups from config")
	}
}

func loadUsers() {
	err := viper.UnmarshalKey("users", &users.Users)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("unable to load users from config")
	}
}
