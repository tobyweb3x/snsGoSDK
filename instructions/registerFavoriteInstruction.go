package instructions

import (
	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type RegisterFavoriteInstruction struct {
	Tag uint8
}

func (rf *RegisterFavoriteInstruction) serialize() ([]byte, error) {
	return borsh.Serialize(*rf)
}

func NewRegisterFavoriteInstruction() *RegisterFavoriteInstruction {
	return &RegisterFavoriteInstruction{
		Tag: 6,
	}
}

func (rf *RegisterFavoriteInstruction) getInstruction(
	programId,
	nameAccount,
	favouriteAccount,
	owner,
	systemProgram,
	optParent solana.PublicKey) *solana.GenericInstruction {

	data, err := rf.serialize()
	if err != nil {
		panic(err)
	}

	keys := solana.AccountMetaSlice{
		{PublicKey: nameAccount, IsSigner: false, IsWritable: false},
		{PublicKey: favouriteAccount, IsSigner: false, IsWritable: true},
		{PublicKey: owner, IsSigner: true, IsWritable: true},
		{PublicKey: systemProgram, IsSigner: false, IsWritable: false},
	}

	if !optParent.IsZero() {
		keys = append(keys, &solana.AccountMeta{PublicKey: optParent, IsSigner: false, IsWritable: false})
	}

	return solana.NewInstruction(
		programId,
		keys,
		data,
	)
}
