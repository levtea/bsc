package ethsync

import (
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"
)

type Header struct {
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
	Uncles           []*Header   `json:"uncles"`
}

func KvBlock(block *types.Block) *BlockOnlyTx {
	var transactionsRes []string
	if len(block.Transactions()) != 0 {
		for _, t := range block.Transactions() {
			transactionsRes = append(transactionsRes, t.Hash().String())
		}
	}
	var unclesRes []*Header
	if len(block.Uncles()) != 0 {
		for _, t := range block.Uncles() {
			unclesRes = append(unclesRes, kvHeader(t))
		}
	}
	return &BlockOnlyTx{
		Header: Header{
			Difficulty:  block.Header().Difficulty.String(),
			Extra:       string(block.Header().Extra),
			GasLimit:    strconv.FormatUint(block.Header().GasLimit, 10),
			GasUsed:     strconv.FormatUint(block.Header().GasUsed, 10),
			Bloom:       string(block.Header().Bloom.Bytes()),
			Coinbase:    block.Header().Coinbase.String(),
			MixDigest:   block.Header().MixDigest.String(),
			Nonce:       strconv.FormatUint(block.Header().Nonce.Uint64(), 10),
			Number:      block.Header().Number.String(),
			ParentHash:  block.Header().ParentHash.String(),
			ReceiptHash: block.Header().ReceiptHash.String(),
			UncleHash:   block.Header().UncleHash.String(),
			Size:        block.Header().Size().String(),
			Root:        block.Header().Root.String(),
			Time:        strconv.FormatUint(block.Header().Time, 10),
		},
		TotalDifficulty:  block.Difficulty(),
		TransactionsRoot: block.Header().TxHash.String(),
		Transactions:     transactionsRes,
		Hash:             block.Hash().String(),
		Uncles:           unclesRes,
	}
}

func kvHeader(header *types.Header) *Header {
	return &Header{
		Difficulty:  header.Difficulty.String(),
		Extra:       string(header.Extra),
		GasLimit:    strconv.FormatUint(header.GasLimit, 10),
		GasUsed:     strconv.FormatUint(header.GasUsed, 10),
		Bloom:       string(header.Bloom.Bytes()),
		Coinbase:    header.Coinbase.String(),
		MixDigest:   header.MixDigest.String(),
		Nonce:       strconv.FormatUint(header.Nonce.Uint64(), 10),
		Number:      header.Number.String(),
		ParentHash:  header.ParentHash.String(),
		ReceiptHash: header.ReceiptHash.String(),
		UncleHash:   header.UncleHash.String(),
		Size:        header.Size().String(),
		Root:        header.Root.String(),
		Time:        strconv.FormatUint(header.Time, 10),
	}
}
