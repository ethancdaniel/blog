package keeper

import (
	"github.com/ethancdaniel/blog/x/blog/types"
)

var _ types.QueryServer = Keeper{}
