package blog

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethancdaniel/blog/x/blog/keeper"
	"github.com/ethancdaniel/blog/x/blog/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	// Set all the comment
	for _, elem := range genState.CommentList {
		k.SetComment(ctx, *elem)
	}

	// Set comment count
	k.SetCommentCount(ctx, uint64(len(genState.CommentList)))

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	// this line is used by starport scaffolding # genesis/module/export
	// Get all comment
	commentList := k.GetAllComment(ctx)
	for _, elem := range commentList {
		elem := elem
		genesis.CommentList = append(genesis.CommentList, &elem)
	}

	return genesis
}
