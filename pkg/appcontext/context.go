package appcontext

import (
	"golang.org/x/net/context"
	"github.com/spf13/viper"
)

type Context struct {
	Viper *viper.Viper
	GoContext context.Context
}
