package nodes

import (
	"github.com/cryptopunkscc/astrald/astral"
	"github.com/cryptopunkscc/astrald/mod/nodes"
	"time"
)

type dbEndpoint struct {
	Identity  *astral.Identity `gorm:"primaryKey"`
	Network   string           `gorm:"primaryKey"`
	Address   string       `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (dbEndpoint) TableName() string {
	return nodes.DBPrefix + "endpoints"
}
