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

type kv struct {
	Prefix string
	PermKV redis.UniversalClient
}

func (kvc *kv) MergeKey(key string) string {
	return kvc.mergeKey(key)
}

func NewKv(prefix string, permKVList []string) *kv {
	kvc := &kv{
		Prefix: prefix,
	}

	if len(permKVList) == 1 {
		kvc.PermKV = redis.NewClient(&redis.Options{
			Addr: permKVList[0],
		})
	} else if len(permKVList) > 1 {
		kvc.PermKV = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    "mymaster",
			SentinelAddrs: permKVList,
		})
		//kvc.PermKV = redis.NewClusterClient(&redis.ClusterOptions{
		//	Addrs:        permKVList,
		//	PoolSize:     32,
		//	MinIdleConns: 32,
		//})
	}
	if kvc.PermKV != nil {
		sc := kvc.PermKV.Ping(context.Background())
		if sc.Err() != nil {
			log.Fatal(errors.New("permanent kv storage service ping failed"))
		}
	}
	return kvc
}

func (kvc *kv) SetInPerm(ctx context.Context, key string, value interface{}) error {
	// if kvc.PermKV != nil {
	// 	return kvc.PermKV.Set(ctx, kvc.mergeKey(key), value, 0).Err()
	// }
	ethlog.Info(fmt.Sprintf("ankr store key is %s", key))
	ethlog.Info(fmt.Sprintf("ankr store value is %s", value))
	return nil
}

func (kvc *kv) mergeKey(key string) string {
	return strings.Join([]string{kvc.Prefix, key}, "")
}
