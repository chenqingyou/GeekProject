package failover

import (
	"GeekProject/homeWork/class2/webook/internal/service/sms"
	"context"
	"sync/atomic"
	"time"
)

type TimeLongFailoverSMSService struct {
	svcs []sms.ServiceSmsInterface
	idx  int32
	cnt  int32 //连续超时的个数
	//阈值，连续超时超过了这个数字，就切换
	threshold        int64
	avgResponseTime  int64 // 新增: 平均响应时间
	lastResponseTime int64 // 新增: 最后一次的响应时间
}

func (t *TimeLongFailoverSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	start := time.Now()             // 开始时间
	idx := atomic.LoadInt32(&t.idx) // 直接加载当前索引
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	responseTime := time.Since(start).Milliseconds() // 响应时间
	// 更新平均响应时间
	newAvg := (t.avgResponseTime + responseTime) / 2
	atomic.StoreInt64(&t.avgResponseTime, newAvg)

	if newAvg > t.lastResponseTime && newAvg > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt64(&t.lastResponseTime, newAvg) // 更新上一次的平均响应时间
		}
		idx = atomic.LoadInt32(&t.idx)
	}
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	default:
		return err
	}
	return err
}
