package dingding

import (
	"fmt"
	"testing"
	"time"
)

func TestDingdingSend(t *testing.T) {
	var svc = NewDingSrv(&Config{
		Secret:    "SECd7fa376f9f96a319485c7891ee95b96f90a6b16aa0305dda9745b093a2179ca2",
		Urls:      "https://oapi.dingtalk.com/robot/send?access_token=f2a22a6ef7f36d1caa825ec1b95949f8d1d0a6af33037084caf801db3a1321f2",
		AtMobiles: nil,
		AtUserIds: nil,
		Interval:  10,
	})
	svc.DingdingSend(content)
	for i := 0; i < 10; i++ {
		svc.DingdingSend(content)
	}
	time.Sleep(time.Minute * 10)
}

func TestEx(t *testing.T) {
	fmt.Printf("%+v\n", GetPanicInfo(content))
}

var content = `
runtime error: (测试发送) invalid memory address or nil pointer dereference
goroutine 609 [running]:
runtime/debug.Stack()
	/Users/dartou/.gvm/gos/go1.17/src/runtime/debug/stack.go:24 +0x7c
pm.dartou.com/dartou/haju_ddz_go/common/tools/recuperate.(*PanicGroup).Go.func1.1(0x140006e0fe0)
	/Users/dartou/workspace/dartou/doudizhu_server_go/common/tools/recuperate/recover.go:93 +0x50
panic({0x104c84160, 0x10636abf0})
	/Users/dartou/.gvm/gos/go1.17/src/runtime/panic.go:1052 +0x2b0
pm.dartou.com/dartou/haju_ddz_go/common/apps/activity/activityclient.(*defaultAgent).LoginInitInfo(0x140008e9db0, {0x1050024b8, 0x140003a90b0}, 0x140008e9da0, {0x0, 0x0, 0x0})
	/Users/dartou/workspace/dartou/doudizhu_server_go/common/apps/activity/activityclient/acitivity.go:52 +0x30
pm.dartou.com/dartou/haju_ddz_go/center/logic.(*Logic).RequestInitInfo(0x140003a91d0, 0x5f5e8c1)
	/Users/dartou/workspace/dartou/doudizhu_server_go/center/logic/activity.go:54 +0x424
pm.dartou.com/dartou/haju_ddz_go/center/logic.(*Logic).Online(0x140003a91d0, 0x1, {0x14000279000, 0x78, 0x80})
	/Users/dartou/workspace/dartou/doudizhu_server_go/center/logic/online.go:88 +0x105c
pm.dartou.com/dartou/haju_ddz_go/center/application.(*CenterServer).Invoke.func1()
	/Users/dartou/workspace/dartou/doudizhu_server_go/center/application/method.go:35 +0x248
pm.dartou.com/dartou/haju_ddz_go/common/tools/recuperate.(*PanicGroup).Go.func1(0x140006e0fe0, 0x140003a90e0)
	/Users/dartou/workspace/dartou/doudizhu_server_go/common/tools/recuperate/recover.go:97 +0x5c
created by pm.dartou.com/dartou/haju_ddz_go/common/tools/recuperate.(*PanicGroup).Go
	/Users/dartou/workspace/dartou/doudizhu_server_go/common/tools/recuperate/recover.go:90 +0x94
`
