package filter

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"./impl"
	"github.com/gomodule/redigo/redis"
	"github.com/willf/bloom"
	"strconv"
)

/*
filter package provides data structures and methods for filter function.

Filter function provide three methods: Init() initialize BloomFilter, Update() set BloomFilter, Filter() filter.

Init() methods required three parameters, and return the object BloomFilterTool.
	server	-	storage address
	iterm	-	elements numbers
	fp		-	the false positive rate of a particular Bloom filter for

Update() pass key and values, then set the bloom filter, and return success numbers.

Filter() pass key and values, return the values not in bloom filter.
*/

type BloomFilterTool struct {
	n    uint
	fp   float64
	pool *redis.Pool
}

// kv：
// key: prefix + uid
// value: BloomFilterStorage
//type BloomFilterStorage struct {
//	version int `json:"version"`
//	preBloomFilter *bloom.BloomFilter `json:"pre_bloom_filter"`
//	curBloomFilter *bloom.BloomFilter `json:"cur_bloom_filter"`
//}
type BloomFilterStorage struct {
	Version        int                `json:"version"`
	PreBloomFilter *bloom.BloomFilter `json:"preBloomFilter"`
	CurBloomFilter *bloom.BloomFilter `json:"curBloomFilter"`
	PreItems       uint               `json:"preItems"`
	CurItems       uint               `json:"curItems"`
}

var (
	threshold float64 = 0.7
)

func (bft *BloomFilterTool) SetTheshold(param float64) {
	threshold = param
}
func (bft *BloomFilterTool) GetThreshold() float64 {
	return threshold
}

// initialize BloomFilterTool
func Init(server string, n uint, fp float64) *BloomFilterTool {
	return &BloomFilterTool{n, fp, impl.PoolInit(server)}
}

func (bft *BloomFilterTool) New(version int, PreBloomFilter *bloom.BloomFilter, CurBloomFilter *bloom.BloomFilter) *BloomFilterStorage {
	return &BloomFilterStorage{version, PreBloomFilter, CurBloomFilter, 0, 0}
}

// set BloomFilter
func (bft *BloomFilterTool) Update(key string, values []uint64) (int, error) {
	var cnt int
	if len(values) < 1 {
		return cnt, errors.New("<key, values>, need values")
	}
	// Before set BloomFilter, we get the object from storage, there are three situations:
	// 1. if an error occur, means no object in the storage, we need build it, and then set it;
	// 2. check the threshold, if the items reach the threshold,
	// we need switch perVersion and curVersion, then modify curVersion and set it;
	// 3. if we get the right version, we only need modify curVersion, then set it;
	version := 0
	var bloomFilterStorage BloomFilterStorage
	json_BloomFilterStorage, err := impl.GetFunc(bft.pool, key)
	if err != nil {
		PreBloomFilter := bloom.NewWithEstimates(bft.n, bft.fp)
		CurBloomFilter := bloom.NewWithEstimates(bft.n, bft.fp)
		bloomFilterStorage = *bft.New(version, PreBloomFilter, CurBloomFilter)
	} else {
		err = json.Unmarshal(json_BloomFilterStorage, &bloomFilterStorage)
		if err != nil {
			return cnt, errors.New("Update json unmarshal error")
		}
		totalItems, err := strconv.ParseFloat(strconv.FormatUint(uint64(bloomFilterStorage.CurItems+bloomFilterStorage.PreItems), 10), 64)
		totalCap, err := strconv.ParseFloat(strconv.FormatUint(uint64(bft.n), 10), 64)
		if err != nil {
			return cnt, errors.New("Update calculate version error")
		}
		if totalItems/totalCap > threshold { // 大于阈值时切换版本
			bloomFilterStorage.Version = 1 - version
			tmpBFS := bloomFilterStorage.PreBloomFilter
			bloomFilterStorage.PreBloomFilter = bloomFilterStorage.CurBloomFilter
			bloomFilterStorage.CurBloomFilter = tmpBFS
			bloomFilterStorage.CurBloomFilter.ClearAll()
			bloomFilterStorage.PreItems = bloomFilterStorage.CurItems
			bloomFilterStorage.CurItems = 0
		}
	}
	for _, ele := range values {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, ele)
		if !bloomFilterStorage.CurBloomFilter.Test(buf) &&
			!bloomFilterStorage.PreBloomFilter.Test(buf) {
			bloomFilterStorage.CurBloomFilter.Add(buf)
			cnt++
			bloomFilterStorage.CurItems++
		}
	}
	json_BloomFilterStorage, err = json.Marshal(bloomFilterStorage)
	if err != nil {
		return 0, errors.New(" Update json marshal error")
	}
	err = impl.SetFunc(bft.pool, key, []byte(json_BloomFilterStorage))
	if err != nil {
		return 0, errors.New("Update SetFunc error")
	}
	return cnt, err
}

func (bft *BloomFilterTool) UpdateWithNotConfirmExist(key string, values []uint64) (int, error) {
	var cnt int
	if len(values) < 1 {
		return cnt, errors.New("<key, values>, need values")
	}
	version := 0
	var bloomFilterStorage BloomFilterStorage
	json_BloomFilterStorage, err := impl.GetFunc(bft.pool, key)
	if err != nil {
		PreBloomFilter := bloom.NewWithEstimates(bft.n, bft.fp)
		CurBloomFilter := bloom.NewWithEstimates(bft.n, bft.fp)
		bloomFilterStorage = *bft.New(version, PreBloomFilter, CurBloomFilter)
	} else {
		err = json.Unmarshal(json_BloomFilterStorage, &bloomFilterStorage)
		if err != nil {
			return cnt, errors.New("Update json unmarshal error")
		}
		if bloomFilterStorage.Version != version {
			bloomFilterStorage.Version = version
			tmpBFS := bloomFilterStorage.PreBloomFilter
			bloomFilterStorage.PreBloomFilter = bloomFilterStorage.CurBloomFilter
			bloomFilterStorage.CurBloomFilter = tmpBFS
			bloomFilterStorage.CurBloomFilter.ClearAll()
			bloomFilterStorage.PreItems = bloomFilterStorage.CurItems
			bloomFilterStorage.CurItems = 0
		}
	}
	for _, ele := range values {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, ele)
		if !bloomFilterStorage.CurBloomFilter.Test(buf) &&
			!bloomFilterStorage.PreBloomFilter.Test(buf) {
			bloomFilterStorage.CurBloomFilter.Add(buf)
			cnt++
			bloomFilterStorage.CurItems++
		}
	}
	json_BloomFilterStorage, err = json.Marshal(bloomFilterStorage)
	if err != nil {
		return 0, errors.New(" Update json marshal error")
	}
	err = impl.SetFunc(bft.pool, key, []byte(json_BloomFilterStorage))
	if err != nil {
		return 0, errors.New("Update SetFunc error")
	}
	return cnt, err
}

// BloomFilter filter
func (bft *BloomFilterTool) Filter(key string, values []uint64) ([]uint64, []uint64, error) {
	var res []uint64
	var inres []uint64
	if len(values) < 1 {
		return res, inres, nil
	}
	json_BloomFilterStorage, err := impl.GetFunc(bft.pool, key)
	// if an error occur, means no object in the storage
	if err != nil {
		return values, inres, nil
	}
	var bloomFilterStorage BloomFilterStorage
	err = json.Unmarshal(json_BloomFilterStorage, &bloomFilterStorage)
	if err != nil {
		return values, inres, errors.New("Filter json unmarshal error")
	}
	for _, ele := range values {
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, ele)
		if !bloomFilterStorage.CurBloomFilter.Test(buf) &&
			!bloomFilterStorage.PreBloomFilter.Test(buf) {
			res = append(res, ele)
		} else {
			inres = append(inres, ele)
		}
	}
	return res, inres, err
}

// Get used space
func (bft *BloomFilterTool) GetKeyItems(key string) (uint, error) {
	json_BloomFilterStorage, err := impl.GetFunc(bft.pool, key)
	if err != nil {
		return 0, err
	}
	var bloomFilterStorage BloomFilterStorage
	err = json.Unmarshal(json_BloomFilterStorage, &bloomFilterStorage)
	if err != nil {
		return 0, errors.New("GetUsedSpace json unmarshal error")
	}
	return bloomFilterStorage.PreItems + bloomFilterStorage.CurItems, err
}

// BloomFilter version
func (bft *BloomFilterTool) GetVersion(key string) (int, error) {
	json_BloomFilterStorage, err := impl.GetFunc(bft.pool, key)
	if err != nil {
		return -1, err
	}
	var bloomFilterStorage BloomFilterStorage
	err = json.Unmarshal(json_BloomFilterStorage, &bloomFilterStorage)
	if err != nil {
		return -1, errors.New("GetVersion json unmarshal error")
	}
	return bloomFilterStorage.Version, err
}
