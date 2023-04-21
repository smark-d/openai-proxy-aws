package filter

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/smark-d/openai-proxy-aws/server/comm"
	"sync"
)

var (
	client *redis.Client = nil
	once   sync.Once
)

const (
	OPENAIKEYS         = "OPENAIKEYS"          // store OpenAI's API Key
	CURR_COUNT_PREFIX  = "CURR_COUNT_PREFIX_"  // prefix for custom key's current count
	TOTAL_COUNT_PREFIX = "TOTAL_COUNT_PREFIX_" // prefix for custom key's total count
)

func init() {
	InitRedis()
}

func InitRedis() {
	comm.InitConfig()
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", comm.Conf.Redis.Host, comm.Conf.Redis.Port),
		Password: comm.Conf.Redis.Password,
		DB:       comm.Conf.Redis.DB,
	})
}

func getRedisClient() *redis.Client {
	return client
}

// AddOpenAIkey Add a key to the set of OpenAI keys
func AddOpenAIkey(value string) error {
	return getRedisClient().SAdd(OPENAIKEYS, value).Err()
}

// GetOpenAIkey Get a random key from the set of OpenAI keys
func GetOpenAIkey() (string, error) {
	return getRedisClient().SRandMember(OPENAIKEYS).Result()
}

// RemoveOpenAIkey Remove a key from the set of OpenAI keys
func RemoveOpenAIkey(value string) error {
	return getRedisClient().SRem(OPENAIKEYS, value).Err()
}

func SetTotalCount(customKey string, count int64) error {
	// set the total count, expire after 1 day
	return getRedisClient().Set(TOTAL_COUNT_PREFIX+customKey, count, 0).Err()
}

func GetTotalCount(customKey string) (int64, error) {
	// get the total count
	return getRedisClient().Get(TOTAL_COUNT_PREFIX + customKey).Int64()
}

func IncrCurrCount(customKey string) error {
	// add the current count
	return getRedisClient().IncrBy(CURR_COUNT_PREFIX+customKey, 1).Err()
}

func GetCurrCount(customKey string) (int64, error) {
	// get the current count
	return getRedisClient().Get(CURR_COUNT_PREFIX + customKey).Int64()
}
