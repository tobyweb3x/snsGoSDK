package types

// Record represents a list of SNS Records
type Record string

const (
	IPFS       Record = "IPFS"
	ARWV       Record = "ARWV"
	SOL        Record = "SOL"
	ETH        Record = "ETH"
	BTC        Record = "BTC"
	LTC        Record = "LTC"
	DOGE       Record = "DOGE"
	Email      Record = "email"
	Url        Record = "url"
	Discord    Record = "discord"
	Github     Record = "github"
	Reddit     Record = "reddit"
	Twitter    Record = "twitter"
	Telegram   Record = "telegram"
	Pic        Record = "pic"
	SHDW       Record = "SHDW"
	POINT      Record = "POINT"
	BSC        Record = "BSC"
	Injective  Record = "INJ"
	Backpack   Record = "backpack"
	A          Record = "A"
	AAAA       Record = "AAAA"
	CNAME      Record = "CNAME"
	TXT        Record = "TXT"
	Background Record = "background"
	BASE       Record = "BASE"
	IPNS       Record = "IPNS"
)

var RecordV1Size map[Record]uint8 = map[Record]uint8{
	SOL:        96,
	ETH:        20,
	BSC:        20,
	Injective:  20,
	A:          4,
	AAAA:       16,
	Background: 32,
}

type RecordVersion uint8

const (
	VersionUnspecified RecordVersion = 0
	V1                 RecordVersion = 1
	V2                 RecordVersion = 2
)
