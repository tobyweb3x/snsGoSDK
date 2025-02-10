package snsRecord

import "github.com/gagliardetto/solana-go"

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
		panic(err)
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
