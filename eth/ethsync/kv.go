package ethsync

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	ethlog "github.com/ethereum/go-ethereum/log"
	"github.com/go-redis/redis/v8"
)

// perm store for ankr

type KvSync struct {
	Prefix string
	TempKV redis.UniversalClient
}

func (kvc *KvSync) MergeKey(key string) string {
	return kvc.mergeKey(key)
}

func NewKv(prefix string, tempKVList []string) *KvSync {
	kvc := &KvSync{
		Prefix: prefix,
	}

	if len(tempKVList) == 1 {
		kvc.TempKV = redis.NewClient(&redis.Options{
			Addr: tempKVList[0],
		})
	} else if len(tempKVList) > 1 {
		kvc.TempKV = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:        tempKVList,
			PoolSize:     32,
			MinIdleConns: 32,
		})
	}

	if kvc.TempKV != nil {
		sc := kvc.TempKV.Ping(context.Background())
		if sc.Err() != nil {
			log.Fatal(errors.New("temporary kv storage service ping failed"))
		}
	}
	return kvc
}

func (kvc *KvSync) SetInKV(ctx context.Context, key string, value interface{}) error {
	// if kvc.PermKV != nil {
	// 	return kvc.PermKV.Set(ctx, kvc.mergeKey(key), value, 0).Err()
	// }
	ethlog.Info(fmt.Sprintf("ankr store key is %s", key))
	ethlog.Info(fmt.Sprintf("ankr store value is %s", value))
	return nil
}

func (kvc *KvSync) mergeKey(key string) string {
	return strings.Join([]string{kvc.Prefix, key}, "")
}
