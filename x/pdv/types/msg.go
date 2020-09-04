package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgCreatePDV defines a CreatePDV message
type MsgCreatePDV struct {
	Timestamp time.Time      `json:"timestamp"`
	Address   string         `json:"address"`
	Owner     sdk.AccAddress `json:"owner"`
	DataType  PDVType        `json:"type"`
}

// NewMsgSetName is a constructor function for MsgCreatePDV
func NewMsgCreatePDV(timestamp time.Time, value string, dataType PDVType, owner sdk.AccAddress) MsgCreatePDV {
	return MsgCreatePDV{
		Timestamp: timestamp,
		Address:   value,
		Owner:     owner,
		DataType:  dataType,
	}
}

// Route should return the name of the module
func (msg MsgCreatePDV) Route() string { return RouterKey }

// Type should return the action
func (msg MsgCreatePDV) Type() string { return "create_pdv" }

// ValidateBasic runs stateless checks on the message
func (msg MsgCreatePDV) ValidateBasic() error {
	if msg.Owner.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Owner.String())
	}
	if len(msg.Address) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Address cannot be empty")
	}
	if time.Now().Before(msg.Timestamp) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Timestamp can't be in the future")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgCreatePDV) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgCreatePDV) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
