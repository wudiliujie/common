package redis

import (
	"errors"
	_redis "github.com/go-redis/redis"
	"github.com/wudiliujie/common/convert"
	"github.com/wudiliujie/common/log"
	"time"
)

const RedisNil = "redis: nil"

var Redis _redis.Cmdable

//func Init(redisAddr string) {
//	log.Release("初始化redis")
//	cfg := &_redis.ClusterOptions{}
//	cfg.Addrs = []string{redisAddr}
//	cfg.ReadTimeout = 500 * time.Millisecond
//	cfg.WriteTimeout = 500 * time.Millisecond
//	cfg.MinIdleConns = 100
//	cfg.PoolSize = 1024
//	Redis = _redis.NewClusterClient(cfg)
//
//}
func Init(redisdbAddr, pass string, isJq bool) {
	log.Release("redis init ")
	if isJq {
		cfg := &_redis.ClusterOptions{}
		cfg.Addrs = []string{redisdbAddr}
		cfg.ReadTimeout = 500 * time.Millisecond
		cfg.WriteTimeout = 500 * time.Millisecond
		//cfg.MinIdleConns=10
		cfg.MaxConnAge = 30 * time.Minute
		cfg.PoolSize = 64
		Redis = _redis.NewClusterClient(cfg)
	} else {
		cfg := &_redis.Options{
			Addr:     redisdbAddr,
			Password: pass, // no password set
			DB:       0,    // use default DB
		}
		cfg.ReadTimeout = 500 * time.Millisecond
		cfg.WriteTimeout = 500 * time.Millisecond
		//cfg.MinIdleConns=10
		cfg.MaxConnAge = 30 * time.Minute
		cfg.PoolSize = 64
		Redis = _redis.NewClient(cfg)
	}
}

func GetKeyExists(key string) (bool, error) {
	ret := Redis.Exists(key)
	if ret.Err() != nil {
		log.Error("GetKeyExists:%v", ret.Err())
		return false, ret.Err()
	}
	if ret.Val() == 0 {
		return false, nil
	}
	return true, nil
}

// 执行增加操作不判断正负值
func IncrbyFieldNum(key string, val int64) error {
	ret := Redis.IncrBy(key, val)
	return ret.Err()
}

func SetIncrbyHashFiledVal(key string, field interface{}, val interface{}) (int64, error) {
	ret := Redis.HIncrBy(key, convert.ToString(field), convert.ToInt64(val))
	if ret.Err() != nil {
		log.Error("SetHashFiledVal:%v>>>%v>>>%v>>>%v", key, field, val, ret.Err())
	}
	return ret.Val(), ret.Err()
}

//获取key 结果
func GetIntVal64(key string) int64 {
	ret, err := Redis.Get(key).Int64()

	if err != nil && err.Error() != RedisNil {
		log.Error("GetIntVal64:%v", err)
		return 0
	}
	return ret
}

//获取key 结果
func GetIntVal32(key string) int32 {
	return int32(GetIntVal64(key))
}
func GetStringVal(key string) string {
	ret := Redis.Get(key)

	if ret.Err() != nil && ret.Err().Error() != RedisNil {
		log.Error("GetStringVal:%v", ret.Err())
		return ""
	}
	return ret.Val()
}

//获取key 结果
func SetValue(key string, val interface{}) error {
	ret := Redis.Set(key, val, 0)
	if ret.Err() != nil {
		log.Error("SetValue:%v", ret.Err())
		return ret.Err()
	}
	return nil
}

func GetHashKeyExists(key string, field interface{}) (bool, error) {
	ret := Redis.HExists(key, convert.ToString(field))
	if ret.Err() != nil {
		log.Error("GetKeyExists:%v", ret.Err())
		return false, ret.Err()
	}
	return ret.Val(), nil
}

//通用设置hash字段值
func SetHashFiledVal(key string, field interface{}, val interface{}) error {
	ret := Redis.HSet(key, convert.ToString(field), val)
	if ret.Err() != nil {
		log.Error("SetHashFiledVal:%v>>>%v>>>%v>>>%v", key, field, val, ret.Err())
	}
	return ret.Err()

}

//通用设置hash字段值
func HsetNXFiledVal(key string, field interface{}, val interface{}) error {
	ret := Redis.HSetNX(key, convert.ToString(field), convert.ToString(val))
	if ret.Err() != nil {
		log.Error("HsetNXFiledVal:%v>>>%v>>>%v>>>%v", key, field, val, ret.Err())
	}
	return ret.Err()
}

func MapInsertHash(key string, data map[string]interface{}) (bool, error) {
	for k, v := range data {
		if v != nil {
			err := HsetNXFiledVal(key, k, v)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}
func StringMapInsertHash(key string, data map[string]string) (bool, error) {
	for k, v := range data {
		err := HsetNXFiledVal(key, k, v)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func GetHashStringVal(key string, field interface{}) string {
	ret := Redis.HGet(key, convert.ToString(field))
	if ret.Err() != nil && ret.Err().Error() != RedisNil {
		log.Error("GetHashStringVal:%v>>>%v>>>%v", key, field, ret.Err())
		return ""
	}
	return ret.Val()
}

func GetHashAllVal(key string) map[string]string {
	ret := Redis.HGetAll(key)
	if ret.Err() != nil {
		log.Error("GetHashAllVal:%v>>>%v", key, ret.Err())
		return map[string]string{}
	}
	return ret.Val()
}

func GetHashIntVal(key string, field interface{}) (int32, error) {

	ret, err := Redis.HGet(key, convert.ToString(field)).Int()
	if err != nil && err.Error() != RedisNil {
		log.Error("GetHashIntVal:%v>>%v>>%v", key, field, err)
		return 0, err
	}
	return int32(ret), nil
}
func GetHashInt32Val(key string, field interface{}) int32 {
	ret, err := Redis.HGet(key, convert.ToString(field)).Int()
	if err != nil && err.Error() != RedisNil {
		log.Error("GetHashIntVal:%v>>%v>>%v", key, field, err)

		return 0
	}
	return int32(ret)
}
func GetHashIntVal64(key string, field interface{}) int64 {
	ret, err := Redis.HGet(key, convert.ToString(field)).Int64()
	if err != nil && err.Error() != RedisNil {
		log.Error("GetHashIntVal64:%v>>%v>>%v", key, field, err)
		return 0
	}
	return ret
}

//在道具当前数量上增加或减少
func IncrbyHashIntFieldNum(key string, field interface{}, val int32) error {
	ret := Redis.HIncrBy(key, convert.ToString(field), convert.ToInt64(val))

	if ret.Err() != nil {
		log.Error("IncrbyHashIntFieldNum:%v>>%v>>>%v>>>%v", key, field, val, ret.Err())
	}
	n := ret.Val()
	if n < 0 {
		log.Debug("%s操作失败,%v>>%v", key, field, val)
		Redis.HIncrBy(key, convert.ToString(field), -convert.ToInt64(val))
		return errors.New("数量不足")
	}
	return ret.Err()
}

//在道具当前数量上增加或减少
func IncrbyHashInt64FieldNum(key string, field string, val int64) error {
	ret := Redis.HIncrBy(key, convert.ToString(field), val)
	if ret.Err() != nil {
		log.Error("IncrbyHashInt64FieldNum:%v>>%v>>>%v>>%v", key, field, val, ret.Err())
	}
	n := ret.Val()
	if n < 0 {
		log.Error("%s操作失败:%v>>%v", key, field, val)
		Redis.HIncrBy(key, convert.ToString(field), -val)
		return errors.New("数量不足")
	}
	return ret.Err()
}
func DelHashKey(key string, field interface{}) error {
	ret := Redis.HDel(key, convert.ToString(field))
	if ret.Err() != nil {
		log.Error("DelKey：%v", ret.Err())
	}
	return ret.Err()
}

func DelKey(key string) {
	ret := Redis.Del(key)
	if ret.Err() != nil {
		log.Error("DelKey：%v", ret.Err())
	}
}

//在道具当前数量上增加或减少
func IncrbyHashIntFieldNumBack(key string, field interface{}, val int64) (int64, error) {
	ret := Redis.HIncrBy(key, convert.ToString(field), val)
	if ret.Err() != nil {
		log.Error("IncrbyHashIntFieldNum:%v>>%v>>>%v>>>%v", key, field, val, ret.Err())
	}
	if ret.Val() < 0 {
		//log.Debug("%s操作失败,%v>>%v", key, field, val)
		Redis.HIncrBy(key, convert.ToString(field), -val)
		return 0, errors.New("数量不足")
	}
	return ret.Val(), ret.Err()
}

//zset根据分数大小排序
//start:起始排名，0代表第一名
//end： 结束排名，-1代表全部排名
//func ZsetGetRankWithByScore(key string, start string, end string) map[string]string {
//	ret := Redis.ZRangeByScoreWithScores(key, &_redis.ZRangeBy{Min: start, Max: end})
//	if ret.Err() != nil {
//		log.Error("ZsetGetRankWithByScore:%v", ret.Err())
//		return nil
//	}
//	result := make(map[string]string)
//	for _, v := range ret.Val() {
//		result[convert.ToString(v.Member)] = convert.ToString(v.Score)
//	}
//	return result
//}

func ZsetDelMember(key string, member interface{}) {

	ret := Redis.ZRem(key, member)

	if ret.Err() != nil {
		log.Error("ZsetDelMemberByScore：%v", ret.Err())
	}
}

//通用设置hash字段值
func SaddHashFiledVal(key string, val interface{}) error {
	ret := Redis.SAdd(key, val)

	if ret.Err() != nil {
		log.Error("SaddHashFiledVal:%v>>>%v>>>%v", key, val, ret.Err())
	}
	return ret.Err()
}

func SetEx(key string, second int32, value string) {
	ret := Redis.Set(key, value, time.Duration(second)*time.Second)

	if ret.Err() != nil {
		log.Error("SetEx:%v>>>%v>>>%v", key, value, ret.Err())
	}
}
func EXPIRE(key string, second int32) {
	ret := Redis.Expire(key, time.Duration(second)*time.Second)

	if ret.Err() != nil {
		log.Error("EXPIRE:%v>>>%v>>>%v", key, second, ret.Err())
	}
}

//移除过期时间
func PERSIST(key string) {
	ret := Redis.Persist(key)
	if ret.Err() != nil {
		log.Error("EXPIRE:%v>>>%v", key, ret.Err())
	}
}

func INCR(key string) int64 {
	ret := Redis.Incr(key)
	if ret.Err() != nil {
		log.Error("INCR:%v>>>%v", key, ret.Err())
	}
	return ret.Val()
}
func LPush(key string, values ...interface{}) {
	ret := Redis.LPush(key, values...)
	if ret.Err() != nil {
		log.Error("LPush:%v>>>%v", key, ret.Err())
	}
}
func RPop(key string) (int, error) {
	ret := Redis.RPop(key)
	if ret.Err() != nil {
		log.Error("RPop:%v>>>%v", key, ret.Err())
	}
	return ret.Int()
}
func Time() int64 {
	t, err := Redis.Time().Result()
	if err != nil {
		return time.Now().Unix()
	}
	return t.Unix()
}
