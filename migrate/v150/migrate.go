package v150

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	v040 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v040"
	v043 "github.com/cosmos/cosmos-sdk/x/genutil/legacy/v043"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"

	communitytypes "github.com/Decentr-net/decentr/x/community/types"
	operationstypes "github.com/Decentr-net/decentr/x/operations/types"
	tokentypes "github.com/Decentr-net/decentr/x/token/types"
)

func Migrate(appState types.AppMap, clientCtx client.Context) types.AppMap {
	appState = v040.Migrate(appState, clientCtx)
	appState = v043.Migrate(appState, clientCtx)

	v039codec, v040codec := clientCtx.LegacyAmino, clientCtx.Codec
	appState = migrateCommunity(appState, v039codec, v040codec)
	appState = migrateOperations(appState, v039codec, v040codec)
	appState = migrateToken(appState, v039codec, v040codec)

	return appState
}

func migrateCommunity(appState types.AppMap, v039Codec *codec.LegacyAmino, v040Codec codec.Codec) types.AppMap {
	if appState["community"] != nil {
		var oldState CommunityState
		v039Codec.MustUnmarshalJSON(appState["community"], &oldState)

		var newState communitytypes.GenesisState

		for _, v := range oldState.Params.Moderators {
			addr, err := sdk.AccAddressFromBech32(v)
			if err != nil {
				panic(fmt.Errorf("failed to parse community/moderators: %w", err))
			}
			newState.Params.Moderators = append(newState.Params.Moderators, addr)
		}
		newState.Params.FixedGas = communitytypes.FixedGasParams(oldState.Params.FixedGas)

		for _, v := range oldState.Posts {
			newState.Posts = append(newState.Posts, communitytypes.Post{
				Owner:        v.Owner,
				Uuid:         v.UUID,
				Title:        v.Title,
				PreviewImage: v.PreviewImage,
				Category:     communitytypes.Category(v.Category),
				Text:         v.Text,
			})
		}

		for _, v := range oldState.Likes {
			newState.Likes = append(newState.Likes, communitytypes.Like{
				Owner:     v.Owner,
				PostOwner: v.PostOwner,
				PostUuid:  v.PostUUID,
				Weight:    communitytypes.LikeWeight(v.Weight),
			})
		}

		for owner, followed := range oldState.Followers {
			var aa []sdk.AccAddress
			for _, v := range followed {
				addr, err := sdk.AccAddressFromBech32(v)
				if err != nil {
					panic(fmt.Errorf("failed to parse follower for %s: %w", owner, addr))
				}
				aa = append(aa, addr)
			}
			newState.Following[owner] = communitytypes.GenesisState_AddressList{Address: aa}
		}

		appState["community"] = v040Codec.MustMarshalJSON(&newState)
	}

	return appState
}

func migrateOperations(appState types.AppMap, v039Codec *codec.LegacyAmino, v040Codec codec.Codec) types.AppMap {
	if appState["operations"] != nil {
		var oldState OperationsState
		v039Codec.MustUnmarshalJSON(appState["operations"], &oldState)

		var newState operationstypes.GenesisState

		for _, v := range oldState.Params.Supervisors {
			addr, err := sdk.AccAddressFromBech32(v)
			if err != nil {
				panic(fmt.Errorf("failed to parse operations/supervisors: %w", err))
			}
			newState.Params.Supervisors = append(newState.Params.Supervisors, addr)
		}
		newState.Params.FixedGas = operationstypes.FixedGasParams(oldState.Params.FixedGas)
		newState.Params.MinGasPrice = oldState.Params.MinGasPrice

		appState["operations"] = v040Codec.MustMarshalJSON(&newState)
	}

	return appState
}

func migrateToken(appState types.AppMap, v039Codec *codec.LegacyAmino, v040Codec codec.Codec) types.AppMap {
	if appState["token"] != nil {
		var oldState TokenState
		v039Codec.MustUnmarshalJSON(appState["token"], &oldState)

		var newState tokentypes.GenesisState

		newState.Params.RewardsBlockInterval = uint64(oldState.Params.RewardsBlockInterval)
		for k, v := range oldState.Balances {
			newState.Balances[k] = sdk.DecProto{Dec: sdk.NewDecFromInt(v)}
		}
		for k, v := range oldState.Deltas {
			newState.Deltas[k] = sdk.DecProto{Dec: sdk.NewDecFromInt(v)}
		}

		appState["token"] = v040Codec.MustMarshalJSON(&newState)
	}

	return appState
}
