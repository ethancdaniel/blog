package keeper

import (
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethancdaniel/blog/x/blog/types"
	"strconv"
	"time"
)

// GetCommentCount get the total number of comment
func (k Keeper) GetCommentCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	count, err := strconv.ParseUint(string(bz), 10, 64)
	if err != nil {
		// Panic because the count should be always formattable to iint64
		panic("cannot decode count")
	}

	return count
}

// SetCommentCount set the total number of comment
func (k Keeper) SetCommentCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentCountKey))
	byteKey := types.KeyPrefix(types.CommentCountKey)
	bz := []byte(strconv.FormatUint(count, 10))
	store.Set(byteKey, bz)
}

// AppendComment appends a comment in the store with a new id and update the count
func (k Keeper) AppendComment(
	ctx sdk.Context,
	creator string,
	body string,
	postID string,
	argTime string,
) uint64 {
	// Create the comment
	count := k.GetCommentCount(ctx)
	var comment = types.Comment{
		Creator: creator,
		Id:      count,
		Body:    body,
		PostID:  postID,
		Time: argTime,
	}

	// Checks if commenter is author of post, deny comment request if so. Fails silently
	postOwner := k.GetPostOwner(ctx, postID)
	if postOwner == creator {
		return count;
	}

	// Time is stored as string. Parse to time.Time object to compare to ctx.BlockTime()
	const timeFormat = "2006-01-02 15:04:05 -0700"
	parsedTime, _ := time.Parse(timeFormat, comment.Time)
	diffDuration := parsedTime.Sub(ctx.BlockTime())

	// Check if parsedTime is outside of interval (ctx.Blocktime +/- 60 seconds). If so, fail silently
	const acceptableDuration = time.Duration(60) * time.Second
	if diffDuration > acceptableDuration ||
	   diffDuration < -acceptableDuration {
		return count
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	value := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), value)
	store = prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentPostIDKey))

	

	postID_uint, err := strconv.ParseUint(comment.PostID, 10, 64)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, comment.Id)
	b := buf[:n]

	if store.Has(GetCommentIDBytes(postID_uint)) {
		result := store.Get(GetCommentIDBytes(postID_uint))
		result = append(result, b...)
		store.Set(GetCommentIDBytes(postID_uint), result)
	} else {
		store.Set(GetCommentIDBytes(postID_uint), b)
	}
	// Update comment count
	k.SetCommentCount(ctx, count+1)

	return count
}

// SetComment set a specific comment in the store
func (k Keeper) SetComment(ctx sdk.Context, comment types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	b := k.cdc.MustMarshalBinaryBare(&comment)
	store.Set(GetCommentIDBytes(comment.Id), b)
}

// GetComment returns a comment from its id
func (k Keeper) GetComment(ctx sdk.Context, id uint64) types.Comment {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	var comment types.Comment
	k.cdc.MustUnmarshalBinaryBare(store.Get(GetCommentIDBytes(id)), &comment)
	return comment
}

// HasComment checks if the comment exists in the store
func (k Keeper) HasComment(ctx sdk.Context, id uint64) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	return store.Has(GetCommentIDBytes(id))
}

// GetCommentOwner returns the creator of the comment
func (k Keeper) GetCommentOwner(ctx sdk.Context, id uint64) string {
	return k.GetComment(ctx, id).Creator
}

// RemoveComment removes a comment from the store
func (k Keeper) RemoveComment(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	store.Delete(GetCommentIDBytes(id))
}

// GetAllComment returns all comment
func (k Keeper) GetAllComment(ctx sdk.Context) (list []types.Comment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommentKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Comment
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetCommentIDBytes returns the byte representation of the ID
func GetCommentIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetCommentIDFromBytes returns ID in uint64 format from a byte array
func GetCommentIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
