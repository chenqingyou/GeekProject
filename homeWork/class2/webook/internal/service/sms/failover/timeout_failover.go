package failover

import (
	"GeekProject/homeWork/class2/webook/internal/service/sms"
	"context"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.ServiceSmsInterface
	idx  int32
	cnt  int32 //连续超时的个数
	//阈值，连续超时超过了这个数字，就切换
	threshold int32
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	idx := atomic.AddInt32(&t.idx, 1)
	cnt := atomic.AddInt32(&t.cnt, 1)
	if cnt > t.threshold {
		//这里需要切换,新的下标 取余
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//我成功往后挪了一位
			atomic.StoreInt32(&t.cnt, 0)
		}
		//else 就是出现了并发，别人换成功了
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	case nil:
		//你的连续状态被打断
		atomic.StoreInt32(&t.cnt, 0)
	default:
		//-超时偶发
		//其他错误，直接下一个
		return err
	}
	return err
}
