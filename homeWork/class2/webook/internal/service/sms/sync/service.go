package sync

import (
	"GeekProject/homeWork/class2/webook/internal/domain"
	"GeekProject/homeWork/class2/webook/internal/repository"
	"GeekProject/homeWork/class2/webook/internal/service/sms"
	"GeekProject/homeWork/class2/webook/internal/service/sms/retryable"
	"GeekProject/homeWork/class2/webook/pkg/ratelimit_win"
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

/*
设计一个新的容错机制，同步转异步的容错机制。当满足以下两个条件中的任何一个时，将请求转储到数据库，后续再另外启动一个 goroutine异步发送出去。
触发了限流。
判定服务商已经崩溃。
要求:
如何判定服务商已经崩溃，不允许使用课程上的判断机制，你需要设计一个新的判断机制，并且解释这种判定机制的决策理由。
控制异步重试次数，转储到数据库之后，可以重试 N 次，重试间隔你可以自由决策。
不允许写死任何参数，即用户必须可以控制控制参数。
保持面向接口和依赖注入风格。
写明这种容错机制适合什么场景，并且有什么优缺点
针对提出的缺点，写出后续的改进方案。
提供单元测试。
*/
var errLimit = fmt.Errorf("触发了限流")

type TimeLongFailoverSyncSMSService struct {
	svcsms               sms.ServiceSmsInterface
	limiter              ratelimit_win.LimitInterface
	repSMS               repository.SMSRepositoryInterface
	retryableService     retryable.SMSRetryableService
	windowSize           int     // 窗口大小
	responseTimes        []int64 // 保存最近的响应时间
	sumResponse          int64   // 保存当前窗口内所有响应时间的总和
	totalRequests        int64   // 总的请求数
	slowRequests         int64   // 超过平均响应时间的请求数
	slowRequestThreshold float64 // 超过平均响应的比例
}

func NewTimeLongFailoverSMSService(ssi sms.ServiceSmsInterface, limiter ratelimit_win.LimitInterface,
	repSMS repository.SMSRepositoryInterface, retryableService retryable.SMSRetryableService, slowRequestThreshold float64) sms.ServiceSmsInterface {
	return &TimeLongFailoverSyncSMSService{
		svcsms:               ssi,
		limiter:              limiter,
		windowSize:           1000,
		slowRequestThreshold: slowRequestThreshold,
	}
}

func (t *TimeLongFailoverSyncSMSService) Send(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	limited, err := t.limiter.Limited(ctx, "sms:tencent")
	if err != nil || limited {
		return t.handleAsyncFallback(ctx, tpl, args, numbers...)
	}
	err = t.performSend(ctx, tpl, args, numbers...)
	// 如果出现非context.DeadlineExceeded错误，则认为服务商可能已崩溃，使用异步发送。
	if err != nil && err != context.DeadlineExceeded {
		return t.handleAsyncFallback(ctx, tpl, args, numbers...)
	}
	return err
}

func (t *TimeLongFailoverSyncSMSService) performSend(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	start := time.Now()
	err := t.svcsms.Send(ctx, tpl, args, numbers...)
	responseTime := time.Since(start).Milliseconds()
	//计算每次响应的时间。
	avgResponseTime := t.updateAverageResponseTime(responseTime)
	atomic.AddInt64(&t.totalRequests, 1)
	// 如果此次响应时间超过平均响应时间，增加慢请求计数
	if responseTime > avgResponseTime {
		atomic.AddInt64(&t.slowRequests, 1)
	}
	// 计算慢请求的比率
	slowRatio := float64(t.slowRequests) / float64(t.totalRequests)
	// 如果慢请求的比率超过了20%，触发相应的处理逻辑
	if slowRatio > t.slowRequestThreshold {
		return t.handleAsyncFallback(ctx, tpl, args, numbers...)
	}
	return err
}

func (t *TimeLongFailoverSyncSMSService) handleAsyncFallback(ctx context.Context, tpl string, args []sms.NameArg, numbers ...string) error {
	if err := t.repSMS.InsertSMS(ctx, domain.SmsDomain{
		Tpl:     tpl,
		NameArg: args,
		Numbers: numbers,
	}); err != nil {
		return err
	}
	go func() {
		newCtx := context.Background()
		if err := t.syncSend(newCtx, tpl); err != nil {
			fmt.Println("Failed to send asynchronously:", err)
		}
	}()
	return errLimit
}

func (t *TimeLongFailoverSyncSMSService) updateAverageResponseTime(responseTime int64) int64 {
	// 如果当前窗口的大小已经达到了预设的窗口大小
	if len(t.responseTimes) >= t.windowSize {
		// 从总和中移除最旧的响应时间
		oldest := t.responseTimes[0]
		t.sumResponse -= oldest
		// 从窗口中移除最旧的响应时间，以保持窗口大小不变
		t.responseTimes = t.responseTimes[1:]
	}
	// 将新的响应时间添加到窗口和总和中
	t.responseTimes = append(t.responseTimes, responseTime)
	t.sumResponse += responseTime
	// 返回当前窗口的平均响应时间
	return t.sumResponse / int64(len(t.responseTimes))
}

func (t *TimeLongFailoverSyncSMSService) syncSend(ctx context.Context, tpl string) error {
	//先根据id去数据库里面查询这个tpl
	sendDomain, err := t.repSMS.SendSMS(ctx, tpl)
	if err != nil {
		return err
	}
	//然后进行重试
	return t.retryableService.Send(ctx, sendDomain.Tpl, sendDomain.NameArg, sendDomain.Numbers...)
}

/*
缺点：
对瞬时变化反应不足：只看平均响应时间可能会掩盖瞬时的变化。例如，如果一段时间的响应很快，突然之间有几次响应很慢，简单的平均值可能不会反映这种情况。
不考虑长期趋势：只考虑最近的响应时间，而不是更长期的趋势，可能导致过早或过晚地触发故障转移。
对异常值敏感：单个的长时间响应（可能是一个异常值）可能导致故障转移，即使系统总体上还是健康的。
反复切换：如果两个服务提供商的响应时间相差不大，可能会导致反复的故障转移，这可能是不必要的，并可能导致更大的延迟和不稳定性。

改进方法：
加权平均：最近的响应可以被赋予更高的权重，这样系统可以更快地对近期的变化作出反应。
设置恢复阈值：在进行故障转移后，不要立即在响应时间变好时切回，而是设置一个恢复阈值，只有当响应时间比这个阈值好得多时才考虑切回。
考虑其他指标：除了平均响应时间外，还可以考虑其他指标，如成功率、错误率等。
引入退避策略：如果发生切换，可以暂时不使用该服务一段时间，给它时间恢复。
异常值检测：可以使用算法来检测和过滤异常值，确保它们不会对故障转移决策产生不良影响。
设置最小切换间隔：即使检测到需要进行故障转移，也确保两次切换之间有一个最小的时间间隔
*/
