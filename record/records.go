package record

import (
	"snsGoSDK/spl"
	"snsGoSDK/types"

	"github.com/gagliardetto/solana-go/rpc"
)

func GetArweaveRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.ARWV)
}

func GetBackgroundRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Background)
}

func GetBackpackRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Backpack)
}

func GetBSCRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.BSC)
}
func GetBTCRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.BTC)
}
func GetDiscordRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Discord)
}
func GetDogeRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.DOGE)
}
func GetEmailRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Email)
}
func GetETHRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.ETH)
}
func GetGithubRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Github)
}
func GetInjectiveRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Injective)
}
func GetIPFSRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.IPFS)
}
func GetLTCRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.LTC)
}
func GetPICRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Pic)
}
func GetPointRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.POINT)
}
func GetRedditRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Reddit)
}
func GetSHDWRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.SHDW)
}
func GetSOLRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.SOL)
}
func GetSOLRecordRaw(conn *rpc.Client, domain string) (*spl.NameRegistryState, error) {
	return GetRecordRaw(conn, domain, types.SOL)
}
func GetTelegramRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Telegram)
}
func GetTwitterRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Twitter)
}
func GetURLRecord(conn *rpc.Client, domain string) (string, error) {
	return GetRecordDeserialized(conn, domain, types.Url)
}
