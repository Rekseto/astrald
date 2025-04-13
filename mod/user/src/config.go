package user

import "time"

const (
	keyActiveContract       = "mod.user.active_contract"
	minimalContractLength   = time.Hour
	defaultContractValidity = 365 * 24 * time.Hour
)

type Config struct {
	Identity string `yaml:"identity"`
	Public   bool   `yaml:"public"`
}

var defaultConfig = Config{
	Public: true,
}
