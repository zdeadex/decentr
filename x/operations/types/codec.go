package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgDistributeRewards{}, "operations/MsgDistributeRewards", nil)
	cdc.RegisterConcrete(MsgResetAccount{}, "operations/MsgResetAccount", nil)
	cdc.RegisterConcrete(MsgBanAccount{}, "operations/MsgBanAccount", nil)
	cdc.RegisterConcrete(MsgMint{}, "operations/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "operations/MsgBurn", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDistributeRewards{},
		&MsgResetAccount{},
		&MsgBanAccount{},
		&MsgMint{},
		&MsgBurn{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
