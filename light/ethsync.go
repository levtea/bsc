package light

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func syncByHeader(lc *LightChain, chain []*types.Header) error {
	for _, header := range chain {
		log.Info(fmt.Sprintf("ankr header is %s", header.Number.String()))
		// block
		// block, err := lc.GetBlockByHash(context.Background(), header.Hash())
		// if err != nil {
		// 	return errors.New(fmt.Sprintf("ankr GetBlockByHash error is %v", err))
		// }
		// body, err := json.Marshal(ethsync.KvBlock(block, lc.GetTd(block.Hash(), block.NumberU64())))
		// if err != nil {
		// 	log.Error("ankr json Marshal block fail, err is %v", err)
		// }
		// log.Info(fmt.Sprintf("block=%s", string(body)))

		// tx
		// for _, tx := range block.Transactions() {
		// 	txJson, _ := tx.MarshalJSON()
		// 	log.Info(fmt.Sprintf("ankrTx is %s", string(txJson)))
		// }

		// receipt
		// receipts, err := GetBlockReceipts(context.Background(), lc.odr, block.Hash(), block.NumberU64())
		// if err != nil {
		// 	return errors.New(fmt.Sprintf("ankr GetBlockReceipts error is %v", err.Error()))
		// }
		// for _, receipt := range receipts {
		// 	// receiptJson, _ := receipt.MarshalJSON()
		// 	// log.Info(fmt.Sprintf("ankrReceipt is %s", string(receiptJson)))

		// 	log.Info(fmt.Sprintf("ankrReceipt contract is %s", receipt.ContractAddress))

		// }

		receipts, err := GetBlockReceipts(context.Background(), lc.odr, header.Hash(), header.Number.Uint64())
		if err != nil {
			return errors.New(fmt.Sprintf("ankr GetBlockReceipts error is %v", err.Error()))
		}
		for _, receipt := range receipts {
			receiptJson, _ := receipt.MarshalJSON()
			log.Info(fmt.Sprintf("ankr receipt is %s", string(receiptJson)))

			log.Info(fmt.Sprintf("ankr contract is %s", receipt.ContractAddress))

		}
	}
	return nil
}
