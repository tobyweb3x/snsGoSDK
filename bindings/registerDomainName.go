package bindings

import (
	"context"
	"slices"
	"snsGoSDK/instructions"
	spl "snsGoSDK/spl"
	utils "snsGoSDK/utils"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

func RegisterDomainName(
	conn *rpc.Client,
	name string, space uint32,
	buyer,
	buyerTokenAccount,
	mint,
	referrerKey solana.PublicKey,
) ([]*solana.GenericInstruction, error) {

	// basic validation
	if strings.Contains(name, ".") || strings.TrimSpace(name) != name {
		return nil, spl.NewSNSError(spl.InvalidDomain, "The domain name is malformed", nil)
	}

	ixns := make([]*solana.GenericInstruction, 0, 3)
	refIdx := slices.IndexFunc(spl.REFERRERS, func(i solana.PublicKey) bool {
		return i.Equals(referrerKey)
	})

	var (
		refTokenAccount solana.PublicKey
		err             error
	)

	if refIdx != -1 {
		if refTokenAccount, err = spl.GetAssociatedTokenAddressSync(
			referrerKey,
			mint,
			true,
			solana.PublicKey{},
			solana.PublicKey{},
		); err != nil {
			return nil, err
		}

		out, err := conn.GetAccountInfo(
			context.TODO(),
			refTokenAccount,
		)
		if err != nil {
			return nil, err
		}

		if out == nil || out.Value.Data == nil {
			ix := CreateAssociatedTokenAccountIdempotentInstruction(
				buyer,
				refTokenAccount,
				referrerKey,
				mint,
				solana.PublicKey{},
				solana.PublicKey{},
			)

			ixns = append(ixns, ix)
		}
	}

	vault, err := spl.GetAssociatedTokenAddressSync(
		mint,
		spl.VaultOwner,
		true,
		solana.PublicKey{},
		solana.PublicKey{},
	)
	if err != nil {
		return nil, err
	}

	pythFeed, ok := spl.PYTHFeeds[mint.String()]
	if !ok {
		return nil, spl.NewSNSError(spl.PythFeedNotFound, "The Pyth account for the provided mint was not found", nil)
	}

	rf := uint16(refIdx)
	referrerIdOpt := &rf
	if refIdx != -1 {
		referrerIdOpt = nil
	}

	hashed := utils.GetHashedNameSync(name)
	nameAccount, _, err := utils.GetNameAccountKeySync(
		hashed,
		solana.PublicKey{},
		spl.RootDomainAccount)
	if err != nil {
		return nil, err
	}

	hashedReverseLookup := utils.GetHashedNameSync(nameAccount.String())
	reverseLookupAccount, _, err := utils.GetNameAccountKeySync(
		hashedReverseLookup,
		spl.CentralState,
		spl.RootDomainAccount)
	if err != nil {
		return nil, err
	}

	derivedState, _, err := solana.FindProgramAddress(
		[][]byte{nameAccount.Bytes()},
		spl.ResgistryProgramID,
	)
	if err != nil {
		return nil, err
	}

	ixTwo := instructions.NewCreateInstructionV3(
		name,
		space,
		referrerIdOpt,
	).GetInstruction(
		spl.ResgistryProgramID,
		spl.NameProgramID,
		spl.RootDomainAccount,
		nameAccount,
		reverseLookupAccount,
		solana.SystemProgramID,
		spl.CentralState,
		buyer,
		buyerTokenAccount,
		spl.PythMappingAcc,
		solana.MustPublicKeyFromBase58(pythFeed.Product),
		solana.MustPublicKeyFromBase58(pythFeed.Price),
		vault,
		solana.TokenProgramID,
		solana.SysVarRentPubkey,
		derivedState,
		&referrerKey,
	)

	ixns = append(ixns, ixTwo)
	return ixns, nil
}

func CreateAssociatedTokenAccountIdempotentInstruction(
	payer,
	associatedToken,
	owner,
	mint,
	programId,
	associatedTokenProgramId solana.PublicKey,
) *solana.GenericInstruction {

	if programId.IsZero() {
		programId = solana.TokenProgramID
	}
	if associatedTokenProgramId.IsZero() {
		associatedTokenProgramId = solana.SPLAssociatedTokenAccountProgramID
	}

	return BuildAssociatedTokenAccountInstruction(
		payer,
		associatedToken,
		owner,
		mint,
		[]byte{1},
		programId,
		associatedTokenProgramId,
	)
}

func BuildAssociatedTokenAccountInstruction(
	payer solana.PublicKey,
	associatedToken solana.PublicKey,
	owner solana.PublicKey,
	mint solana.PublicKey,
	instructionData []byte,
	programId solana.PublicKey,
	associatedTokenProgramId solana.PublicKey,
) *solana.GenericInstruction {
	keys := solana.AccountMetaSlice{
		{PublicKey: payer, IsSigner: true, IsWritable: true},
		{PublicKey: associatedToken, IsSigner: false, IsWritable: true},
		{PublicKey: owner, IsSigner: false, IsWritable: false},
		{PublicKey: mint, IsSigner: false, IsWritable: false},
		{PublicKey: system.ProgramID, IsSigner: false, IsWritable: false},
		{PublicKey: programId, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		programId,
		keys,
		instructionData,
	)
}
