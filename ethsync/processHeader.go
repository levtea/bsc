package ethsync

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

const (
	bufferSize = 8192
)

var (
	syncInfo *SyncInfo
)

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
