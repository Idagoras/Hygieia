package pm

import (
	"fmt"
	"github.com/casbin/casbin"
	xormadapter "github.com/casbin/xorm-adapter"
	_ "github.com/go-sql-driver/mysql"
)

type Enforcer struct {
	E *casbin.Enforcer
}

func NewEnforcer(database string, dbAddress string, configAddress string) (*Enforcer, error) {
	a := xormadapter.NewAdapter(database, dbAddress, false)
	e := casbin.NewEnforcer(configAddress, a)
	err := e.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load policy : %v", err)
	}
	return &Enforcer{E: e}, nil
}
