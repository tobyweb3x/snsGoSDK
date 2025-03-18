package spl

import (
	"errors"

	"github.com/gagliardetto/solana-go"
)

func GetAssociatedTokenAddressSync(
	mint, owner solana.PublicKey,
	allowOwnerOffCurve bool,

) (solana.PublicKey, error) {

	if !allowOwnerOffCurve && !solana.IsOnCurve(owner.Bytes()) {
		return solana.PublicKey{}, errors.New("token owner is off-curve")
	}

	address, _, err := solana.FindAssociatedTokenAddress(
		owner,
		mint,
	)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return address, nil
}
