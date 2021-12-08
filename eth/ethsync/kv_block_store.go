package ethsync

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

var (
	blkFormat     = "/bk/%s"
	blkHashFormat = "/bk//hash/%s"
	txFormat      = "/tx/%s"
	rtFormat      = "/rt/%s"
)

type blockStore struct {
	kvc *kv
}

func NewBlockStore(kvc *kv) *blockStore {
	return &blockStore{
		kvc: kvc,
	}
}

func (b *blockStore) StoreBlock(ctx context.Context, block *types.Block, receipts types.Receipts, td *big.Int) error {
	var (
		blkHeightHex string
		body         []byte
		err          error
		number       uint64
	)
	number = block.NumberU64()
	blkHeightHex = fmt.Sprintf("0x%x", number)
	log.Info("ankr start save to kv, block is : ", blkHeightHex)

	// store txs & receipts
	if len(receipts) != 0 {
		if err = b.receiptHandler(ctx, receipts); err != nil {
			log.Error("ankr save receipt failed fail, block is ", blkHeightHex)
			return err
		}
	}
	if len(block.Body().Transactions) != 0 {
		if err = b.txsHandler(ctx, block.Body().Transactions); err != nil {
			log.Error("ankr save transactions fail, block is ", blkHeightHex)
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
	if err = b.kvc.SetInPerm(ctx, blkPath, body); err != nil {
		log.Error("ankr save block body fail, block is ", blkHeightHex)
		return err
	}
	if err = b.kvc.SetInPerm(ctx, hashPath, b.kvc.MergeKey(blkPath)); err != nil {
		log.Error("ankr save block hashPath fail, block is ", blkHeightHex)
		return err
	}
	return nil
}

// func (b *blockStore) blockHander(block *types.Block) error{
// 	blkPath := fmt.Sprintf(blkFormat, blkHeightHex)
// 	hashPath := fmt.Sprintf(blkHashFormat, newBlk.Hash)
// 	// make body with header and body
// 	body := json.Marshal(KvBlock(block))
// 	if err = b.kvc.SetInPerm(ctx, blkPath, body); err != nil {
// 		log.Error().Uint64("number", number).Str("numberHex", blkHeightHex).Err(err).Msg("store block failed")
// 		return err
// 	}
// 	return nil
// }

func (b *blockStore) receiptHandler(ctx context.Context, rts types.Receipts) error {
	var (
		rtbs []byte
		err  error
	)

	for _, rt := range rts {
		if rt.TxHash.String() == "" {
			continue
		}
		rtbs, _ = rt.MarshalJSON()
		if err = b.kvc.SetInPerm(ctx, fmt.Sprintf(rtFormat, rt.TxHash), rtbs); err != nil {
			return err
		}
	}
	return nil
}

func (b *blockStore) txsHandler(ctx context.Context, rts types.Transactions) error {
	var (
		txb []byte
		err error
	)
	for _, tx := range rts {
		if tx.Hash().String() == "" {
			continue
		}
		txb, _ = tx.MarshalJSON()
		if err = b.kvc.SetInPerm(ctx, fmt.Sprintf(txFormat, tx.Hash()), txb); err != nil {
			return err
		}
	}
	return nil
}

func (b *blockStore) blocksHander(cxt context.Context)
