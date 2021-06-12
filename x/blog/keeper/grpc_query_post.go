// x/block/keeper/grpc_query_post.go
package keeper

import (
	"context"
	"strconv"
	"strings"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethancdaniel/blog/x/blog/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PostAll(c context.Context, req *types.QueryAllPostRequest) (*types.QueryAllPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var posts []*types.Post
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	postStore := prefix.NewStore(store, types.KeyPrefix(types.PostKey))
	postIDStore := prefix.NewStore(store, types.KeyPrefix(types.CommentPostIDKey))

	pageRes, err := query.Paginate(postStore, req.Pagination, func(key []byte, value []byte) error {
		var post types.Post
		if err := k.cdc.UnmarshalBinaryBare(value, &post); err != nil {
			return err
		}

		postID_uint, int_err := strconv.ParseUint(post.Id, 10, 64)
		if int_err != nil {
			panic(int_err)
		}
		b := postIDStore.Get(GetCommentIDBytes(postID_uint))

		post.Comments = formatCommentIDs(b)
		posts = append(posts, &post)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPostResponse{Post: posts, Pagination: pageRes}, nil
}

// Converts byte array to string of comma-separated base10 numbers
func formatCommentIDs (arr []byte) string {
	result := fmt.Sprintf("%d", arr)
	result = result[1:len(result) - 1]
	result = strings.ReplaceAll(result, " ", ", ")

	return result 
}

func (k Keeper) Post(c context.Context, req *types.QueryGetPostRequest) (*types.QueryGetPostResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var post types.Post
	ctx := sdk.UnwrapSDKContext(c)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PostKey))
	k.cdc.MustUnmarshalBinaryBare(store.Get(types.KeyPrefix(types.PostKey+req.Id)), &post)

	return &types.QueryGetPostResponse{Post: &post}, nil
}
