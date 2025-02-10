package recordv2

import "snsGoSDK/types"

// Set of records that are UTF-8 encoded strings
var UTF8Encoded = map[types.Record]struct{}{
	types.IPFS:     {},
	types.ARWV:     {},
	types.LTC:      {},
	types.DOGE:     {},
	types.Email:    {},
	types.Url:      {},
	types.Discord:  {},
	types.Github:   {},
	types.Reddit:   {},
	types.Twitter:  {},
	types.Telegram: {},
	types.Pic:      {},
	types.SHDW:     {},
	types.POINT:    {},
	types.Backpack: {},
	types.TXT:      {},
	types.CNAME:    {},
	types.BTC:      {},
	types.IPNS:     {},
}

var EVMRecords = map[types.Record]struct{}{
	types.ETH:  {},
	types.BSC:  {},
	types.BASE: {},
}
