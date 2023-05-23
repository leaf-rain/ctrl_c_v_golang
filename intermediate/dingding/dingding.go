package dingding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leaf-rain/ctrl_c_v_golang/intermediate/recuperate"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Config 钉钉配置
type Config struct {
	Secret    string   `yaml:"Secret" json:"Secret,omitempty"`       // 密钥
	Urls      string   `yaml:"Urls" json:"Urls,omitempty"`           // 路径
	AtMobiles []string `yaml:"AtMobiles" json:"AtMobiles,omitempty"` // @的手机号（固定）
	AtUserIds []string `yaml:"AtUserIds" json:"AtUserIds,omitempty"` // @的用户id（固定）
	Interval  int64    `yaml:"Interval" json:"Interval,omitempty"`   // 相同信息发送间隔,单位：(秒)
}

// Option 发送小时时可以支持的额外参数
type Option struct {
	msg       string   // 需要发送的信息
	AtMobiles []string // @的手机号（额外添加）
	AtUserIds []string // @的用户id（额外添加）
}

type DingSrv struct {
	opt     *Config
	lock    *sync.Mutex // 锁
	ch      chan Option
	repeat  *sync.Map
	routine *sync.Map
}

func NewDingSrv(opt *Config) *DingSrv {
	svc := &DingSrv{
		opt:  opt,
		lock: new(sync.Mutex),
		ch:   make(chan Option, 1024),
	}
	if opt.Interval > 0 {
		svc.repeat = new(sync.Map)
		svc.routine = new(sync.Map)
	}
	go func() {
		svc.watcher()
	}()
	return svc
}

type ResultForDingDing struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// DingdingSend 发送文本类小时
func (d *DingSrv) DingdingSend(message string) {
	d.ch <- Option{
		msg: message,
	}
}

func (d *DingSrv) Close() {
	close(d.ch)
}

func (d *DingSrv) watcher() {
	for {
		if opt, ok := <-d.ch; ok {
			if d.opt.Interval > 0 {
				var info = GetPanicInfo(opt.msg)
				var now = float64(time.Now().Unix())
				if info != nil && info.Line != "" {
					d.delaySend(info.Line)
					timeStampI, ok2 := d.repeat.Load(info.Line)
					if !ok2 {
						d.repeat.Store(info.Line, now+0.000001)
						_ = d.dingdingSend(opt)
					} else {
						old := timeStampI.(float64)
						if now-old > float64(d.opt.Interval) {
							// 超过间隔，再次打印堆栈
							d.repeat.Store(info.Line, now+0.000001)
							_ = d.dingdingSend(opt)
						} else {
							// 没有超过间隔，计数器+1
							d.repeat.Store(info.Line, old+0.000001)
						}
					}
				}
			} else {
				_ = d.dingdingSend(opt)
			}
		} else {
			break
		}
	}
}

func (d *DingSrv) delaySend(line string) {
	if d.opt.Interval > 0 {
		_, ok := d.routine.Load(line)
		if !ok {
			d.routine.Store(line, struct{}{})
			recuperate.GoSafe(func() {
				defer d.routine.Delete(line)
				for {
					time.Sleep(time.Second * time.Duration(d.opt.Interval))
					timeStampI, ok2 := d.repeat.Load(line)
					if !ok2 {
						break
					} else {
						var now = float64(time.Now().Unix())
						old := timeStampI.(float64)
						if now-old < float64(d.opt.Interval) {
							// 小于直接忽视
							continue
						} else {
							var times = FloatDecimal(old)
							if times > 1 {
								var context = fmt.Sprintf("当前异常%s,在%ds内重复%d次！",
									line, d.opt.Interval, times)
								_ = d.dingdingSend(Option{msg: context})
							}
							break
						}
					}
				}
			})
		}
	}
}

// DingdingDirectSend 直接发送消息
func (d *DingSrv) DingdingDirectSend(message string, opt ...Option) (bool, error) {
	var atMobiles, atUserIds []string
	for _, item := range opt {
		if len(item.AtMobiles) > 0 || len(item.AtUserIds) > 0 {
			atMobiles = append(atMobiles, item.AtMobiles...)
			atUserIds = append(atUserIds, item.AtUserIds...)
		}
	}
	err := d.dingdingSend(Option{
		msg:       message,
		AtMobiles: atMobiles,
		AtUserIds: atUserIds,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *DingSrv) dingdingSend(opt Option) error {
	var reqUrl = d.opt.Urls + sign(d.opt.Secret)
	var atMobiles, atUserIds []string
	if len(opt.AtMobiles) > 0 || len(opt.AtUserIds) > 0 {
		atMobiles = append(d.opt.AtMobiles, opt.AtMobiles...)
		atUserIds = append(d.opt.AtUserIds, opt.AtUserIds...)
	}
	var request, err = newRequestBody(atMobiles, atUserIds, opt.msg, false)
	if err != nil {
		log.Println("[DingdingSend] newRequestBody failed",
			zap.Any("setting", d.opt),
			zap.Any("message", opt.msg),
			zap.Error(err))
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", reqUrl, strings.NewReader(request))
	if err != nil {
		log.Println("[DingdingSend] NewRequest failed",
			zap.Any("setting", d.opt),
			zap.Any("message", opt.msg),
			zap.Error(err))
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		log.Println("[DingdingSend] request Do failed",
			zap.Any("setting", d.opt),
			zap.Any("message", opt.msg),
			zap.Error(err))
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("[DingdingSend] ioutil.ReadAll failed",
			zap.Any("setting", d.opt),
			zap.Any("message", opt.msg),
			zap.Error(err))
		return err
	}
	var sta ResultForDingDing
	err = json.Unmarshal(body, &sta)
	if err != nil {
		log.Println("[DingdingSend] json.Unmarshal failed",
			zap.Any("setting", d.opt),
			zap.Any("message", opt.msg),
			zap.Error(err))
		return err
	}
	if sta.Errcode == 0 || sta.Errmsg == "ok" {
		log.Println("[DingdingSend] success",
			zap.Any("message", opt.msg))
		return nil
	}
	return errors.New("Error response:---" + string(body))
}

func newRequestBody(atm, atu []string, data string, isAtAll bool) (string, error) {
	reqBody := struct {
		At struct {
			AtMobiles []string `json:"atMobiles"`
			AtUserIds []string `json:"atUserIds"`
			IsAtAll   bool     `json:"isAtAll"`
		} `json:"at"`
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
		Msgtype string `json:"msgtype"`
	}{}
	reqBody.Text.Content = data
	reqBody.Msgtype = "text"
	reqBody.At.AtMobiles = atm
	reqBody.At.AtUserIds = atu
	reqBody.At.IsAtAll = isAtAll
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	return string(reqData), nil
}

func sign(secret string) string {
	timestamp := fmt.Sprint(time.Now().UnixNano() / 1000000)
	secStr := timestamp + "\n" + secret
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(secStr))
	sum := h.Sum(nil)
	encode := base64.StdEncoding.EncodeToString(sum)
	urlEncode := url.QueryEscape(encode)
	return "&timestamp=" + timestamp + "&sign=" + urlEncode
}
