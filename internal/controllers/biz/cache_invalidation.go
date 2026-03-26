package biz

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/wangyahua6688-maker/tk-common/utils/ctxx"
)

// invalidatePublicLotteryCaches 清除首页/开奖看板/开奖现场的 Redis 缓存。
// 说明：
// 1. 首页彩种按钮时间来自 /public/home，因此要清首页概览缓存；
// 2. 开奖区与直播区来自 /public/special-lotteries/{id}/dashboard，因此要清单彩种看板缓存；
// 3. 开奖现场页使用独立缓存，因此同样需要按彩种清理。
func invalidatePublicLotteryCaches(ctx context.Context, specialLotteryID uint) error {
	// 管理后台没有 Redis 时直接跳过，不把缓存能力当成硬依赖。
	redisClient, ok := ctxx.Get[*redis.Client](ctx, ctxx.RedisKey)
	if !ok || redisClient == nil {
		return nil
	}

	// 统一删除本次彩种配置变更会影响到的公开缓存键。
	keys := []string{
		"tk:home:overview:v1",
		fmt.Sprintf("tk:business:dashboard:%d", specialLotteryID),
		fmt.Sprintf("tk:live_scene:page:v1:%d", specialLotteryID),
	}

	return redisClient.Del(ctx, keys...).Err()
}
