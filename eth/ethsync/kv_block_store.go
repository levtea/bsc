package ethsync

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

var (
	blkFormat     = "/bk/%s"
	blkHashFormat = "/bk//hash/%s"
	txFormat      = "/tx/%s"
	rtFormat      = "/rt/%s"
)

type BlockStore struct {
	kvc *KvSync
}

func NewBlockStore(kvc *KvSync) *BlockStore {
	return &BlockStore{
		kvc: kvc,
	}
}

func (b *BlockStore) StoreBlock(ctx context.Context, block *types.Block, receipts types.Receipts) error {
	var (
		blkHeightHex string
		body         []byte
		err          error
		number       uint64
	)
	number = block.NumberU64()
	blkHeightHex = fmt.Sprintf("0x%x", number)
	log.Info(fmt.Sprintf("ankr start save to kv, block is : ", blkHeightHex))

	// store txs & receipts
	if len(receipts) != 0 {
		if err = b.receiptHandler(ctx, receipts); err != nil {
			log.Error(fmt.Sprintf("ankr save receipt failed fail, block is : %s", blkHeightHex))
			return err
		}
	}
	if len(block.Body().Transactions) != 0 {
		if err = b.txsHandler(ctx, block.Body().Transactions); err != nil {
			log.Error(fmt.Sprintf("ankr save transactions failed fail, block is : %s", blkHeightHex))
			return err
		}
	}

	// save new block into kv
	body, err = json.Marshal(KvBlock(block))
	if err != nil {
		log.Error("ankr json Marshal block fail, block is ", blkHeightHex)
		return err
	}
	blkPath := fmt.Sprintf(blkFormat, blkHeightHex)
	hashPath := fmt.Sprintf(blkHashFormat, block.Hash().String())
	// make body with header and body
	if err = b.kvc.SetInKV(ctx, blkPath, body); err != nil {
		log.Error("ankr save block body fail, block is ", blkHeightHex)
		return err
	}
	if err = b.kvc.SetInKV(ctx, hashPath, b.kvc.MergeKey(blkPath)); err != nil {
		log.Error("ankr save block hashPath fail, block is ", blkHeightHex)
		return err
	}
	return nil
}

func (b *BlockStore) receiptHandler(ctx context.Context, rts types.Receipts) error {
	var (
		rtbs []byte
		err  error
	)

	for _, rt := range rts {
		if rt.TxHash.String() == "" {
			continue
		}
		rtbs, _ = rt.MarshalJSON()
		if err = b.kvc.SetInKV(ctx, fmt.Sprintf(rtFormat, rt.TxHash), rtbs); err != nil {
			return err
		}
	}
	return nil
}

func (b *BlockStore) txsHandler(ctx context.Context, rts types.Transactions) error {
	var (
		txb []byte
		err error
	)
	for _, tx := range rts {
		if tx.Hash().String() == "" {
			continue
		}
		txb, _ = tx.MarshalJSON()
		if err = b.kvc.SetInKV(ctx, fmt.Sprintf(txFormat, tx.Hash()), txb); err != nil {
			return err
		}
	}
	return nil
}
