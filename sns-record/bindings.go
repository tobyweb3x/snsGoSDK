package snsRecord

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
)

var (
	// SNSRecordsID is the mainnet program ID
	//  SNSRecordsID = solana.MustPublicKeyFromBase58("HP3D4D1ZCmohQGFVms2SS4LCANgJyksBf5s1F77FuFjZ")
	SNSRecordsID = solana.MustPublicKeyFromBase58("HP3D4D1ZCmohQGFVms2SS4LCANgJyksBf5s1F77FuFjZ")

	// CentralStateSNSRecords central state account
	//  CentralStateSNSRecords, _, _ := solana.FindProgramAddress(
	//  [][]byte{SNSRecordsID.Bytes()}, SNSRecordsID)
	CentralStateSNSRecords solana.PublicKey
)

func init() {
	out, _, err := solana.FindProgramAddress(
		[][]byte{SNSRecordsID.Bytes()}, SNSRecordsID)
	if err != nil {
		panic(fmt.Errorf("failed to init package snsRecord: %w", err))
	}
	CentralStateSNSRecords = out
}

func AllocateAndPostRecord(
	feePayer,
	recordKey,
	domainKey,
	domainOwner,
	nameProgramID,
	programId solana.PublicKey,
	record string, content []byte,
) (*solana.GenericInstruction, error) {

	ix, err := NewAllocateAndPostRecordInstruction(record, content).GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}

func DeleteRecord(
	feePayer,
	domainKey,
	domainOwner,
	recordKey,
	nameProgramID,
	programId solana.PublicKey,
) (*solana.GenericInstruction, error) {

	ix, err := NewDeleteRecordInstruction().GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}

func ValidateEthSignature(
	feePayer,
	recordKey,
	domainKey,
	domainOwner,
	nameProgramID,
	programId solana.PublicKey,
	validation Validation,
	signature, expectedPubkey []byte,
) (*solana.GenericInstruction, error) {

	ix, err := NewValidateEthereumSignatureInstruction(
		uint8(validation),
		signature,
		expectedPubkey,
	).GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}

func EditRecord(
	feePayer,
	recordKey,
	domainKey,
	domainOwner,
	nameProgramID,
	programId solana.PublicKey,
	record string, content []byte,
) (*solana.GenericInstruction, error) {

	ix, err := NewEditRecordInstruction(record, content).GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}

func ValidateSolanaSignature(
	feePayer,
	recordKey,
	domainKey,
	domainOwner,
	verifier,
	nameProgramID,
	programId solana.PublicKey,
	staleness bool,
) (*solana.GenericInstruction, error) {

	ix, err := NewValidateSolanaSignatureInstruction(
		staleness,
	).GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
		verifier,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}

func WriteRoa(
	feePayer,
	nameProgramID,
	recordKey,
	domainKey,
	domainOwner,
	roaId,
	programId solana.PublicKey,
) (*solana.GenericInstruction, error) {

	ix, err := NewWriteRoaInstruction(roaId.Bytes()).GetInstruction(
		programId,
		solana.SystemProgramID,
		nameProgramID,
		feePayer,
		recordKey,
		domainKey,
		domainOwner,
		CentralStateSNSRecords,
	)
	if err != nil {
		return nil, err
	}

	return ix, nil
}
