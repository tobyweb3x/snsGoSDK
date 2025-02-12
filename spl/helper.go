package spl

import (
	"errors"

	"github.com/gagliardetto/solana-go"
)

func GetAssociatedTokenAddressSync(
	mint, owner solana.PublicKey,
	allowOwnerOffCurve bool,
	programId, associatedTokenProgramId solana.PublicKey,
) (solana.PublicKey, error) {

	if !allowOwnerOffCurve && !solana.IsOnCurve(owner.Bytes()) {
		return solana.PublicKey{}, errors.New("token owner is off-curve")
	}

	if associatedTokenProgramId.IsZero() {
		associatedTokenProgramId = solana.SPLAssociatedTokenAccountProgramID
	}

	if programId.IsZero() {
		programId = solana.TokenProgramID
	}

	seed := [][]byte{
		owner.Bytes(),
		programId.Bytes(),
		mint.Bytes(),
	}
	address, _, err := solana.FindProgramAddress(seed, associatedTokenProgramId)
	if err != nil {
		return solana.PublicKey{}, err
	}

	return address, nil
}
