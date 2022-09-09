package ethsync

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

const (
	bufferSize = 8192
)

var (
	syncInfo *SyncInfo
)

type Header struct {
	Hash        string `json:"hash"`
	Difficulty  string `json:"difficulty"`
	Extra       string `json:"extraData"`
	GasLimit    string `json:"gasLimit"`
	GasUsed     string `json:"gasUsed"`
	Bloom       string `json:"logsBloom"`
	Coinbase    string `json:"miner"`
	MixDigest   string `json:"mixHash"`
	Nonce       string `json:"nonce"`
	Number      string `json:"number"`
	ParentHash  string `json:"parentHash"`
	ReceiptHash string `json:"receiptsRoot"`
	UncleHash   string `json:"sha3Uncles"`
	Size        string `json:"size"`
	Root        string `json:"stateRoot"`
	Time        string `json:"timestamp"`
}

type BlockOnlyTx struct {
	Header
	TotalDifficulty  interface{} `json:"totalDifficulty"`
	TransactionsRoot string      `json:"transactionsRoot"`
	Transactions     []string    `json:"transactions"`
	Hash             string      `json:"hash"`
	Uncles           []string    `json:"uncles"`
	Size             string      `json:"size"`
}

func BytesToInt64String(buf []byte) string {
	res := "0x" + hex.EncodeToString(buf)
	return res
}

func KvBlock(block *types.Block) *BlockOnlyTx {
	var transactionsRes []string
	if len(block.Transactions()) != 0 {
		for _, t := range block.Transactions() {
			transactionsRes = append(transactionsRes, t.Hash().Hex())
		}
	}
	unclesRes := make([]string, 0, len(block.Uncles()))
	if len(block.Uncles()) != 0 {
		for _, t := range block.Uncles() {
			unclesRes = append(unclesRes, t.Hash().Hex())
		}
	}
	return &BlockOnlyTx{
		Header: Header{
			Hash:        block.Hash().Hex(),
			Difficulty:  "0x" + strconv.FormatUint(block.Header().Difficulty.Uint64(), 16),
			Extra:       BytesToInt64String(block.Header().Extra),
			GasLimit:    "0x" + strconv.FormatUint(block.Header().GasLimit, 16),
			GasUsed:     "0x" + strconv.FormatUint(block.Header().GasUsed, 16),
			Bloom:       BytesToInt64String(block.Header().Bloom.Bytes()),
			Coinbase:    block.Header().Coinbase.String(),
			MixDigest:   block.Header().MixDigest.String(),
			Nonce:       "0x" + strconv.FormatUint(block.Header().Nonce.Uint64(), 16),
			Number:      "0x" + strconv.FormatUint(block.Header().Number.Uint64(), 16),
			ParentHash:  block.Header().ParentHash.String(),
			ReceiptHash: block.Header().ReceiptHash.String(),
			UncleHash:   block.Header().UncleHash.String(),
			Size:        "0x" + strconv.FormatUint(uint64(block.Header().Size()), 16),
			Root:        block.Header().Root.String(),
			Time:        "0x" + strconv.FormatUint(block.Header().Time, 16),
		},
		TotalDifficulty:  "0x" + block.Difficulty().Text(16),
		TransactionsRoot: block.Header().TxHash.String(),
		Transactions:     transactionsRes,
		Hash:             block.Hash().Hex(),
		Uncles:           unclesRes,
		Size:             "0x" + strconv.FormatUint(uint64(block.Size()), 16),
	}
}

type SyncInfo struct {
	HeaderChan chan []*types.Header
}

func NewSyncInfo() *SyncInfo {
	headersChan := make(chan []*types.Header, bufferSize)
	return &SyncInfo{
		HeaderChan: headersChan,
	}
}

func init() {
	log.Info("ankr start get headers")
	syncInfo = NewSyncInfo()
	go func() {
		for headers := range syncInfo.HeaderChan {

			for _, header := range headers {
				log.Info(fmt.Sprintf("ankr header is %s", header.Number.String()))
			}

			// headerJson, _ := header.MarshalJSON()
			// log.Info("ankr get headers", header.Number.Uint64(), headerJson)

		}
	}()
}

func Extract(headers []*types.Header) {
	syncInfo.HeaderChan <- headers
}
