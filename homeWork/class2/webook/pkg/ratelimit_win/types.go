package ratelimit_win

import (
	"context"
)

type LimitInterface interface {
	//有没有触发限流，key就是限流对象
	//bool代表是否限流，true就是要限流
	Limited(ctx context.Context, key string) (bool, error)
}
