package redisNotificationDispatcher

import (
	"path"
	"runtime"
	"testing"

	"github.com/go-redis/redis/v7"
)

type _TestSubscriber struct {
	t     *testing.T
	msgs  []*redis.Message
	errs  []error
	fulls []bool

	onMessage func(*redis.Message, int64)
}

var testSubscriber _TestSubscriber

func (me *_TestSubscriber) OnNotifyMessage(msg *redis.Message, skipped int64) {
	me.onMessage(msg, skipped)
}

func (me *_TestSubscriber) OnNotifyError(err error) {
	me.errs = append(me.errs, err)
	me.t.Logf(`OnError %v`, err)
}

func (me *_TestSubscriber) GetReady(t *testing.T) {
	me.msgs = nil
	me.errs = nil
	me.fulls = nil
	me.t = t
}

func TestInit(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dir := path.Dir(file)
	os.Setenv(`CONFIG_FILE`, path.Join(dir, `test.yaml`))
	config.Init()

}
