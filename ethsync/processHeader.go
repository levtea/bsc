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
	HeaderChan chan *types.Header
}

func NewSyncInfo() *SyncInfo {
	headerChan := make(chan *types.Header, bufferSize)
	return &SyncInfo{
		HeaderChan: headerChan,
	}
}

func init() {
	log.Info("ankr start get headers")
	syncInfo = NewSyncInfo()
	go func() {
		for header := range syncInfo.HeaderChan {

			// headerJson, _ := header.MarshalJSON()
			// log.Info("ankr get headers", header.Number.Uint64(), headerJson)
			log.Info(fmt.Sprintf("ankr header is %s", header.Number.String()))

		}
	}()
}

func Extract(header *types.Header) {
	syncInfo.HeaderChan <- header
}
