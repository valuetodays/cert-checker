package main

import (
	"bytes"
	"cert-checker/internal/checker"
	"cert-checker/internal/config"
	"cert-checker/internal/notifier"
	"encoding/json"
	"flag"
	"github.com/robfig/cron/v3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	configFile := flag.String("config", "config/config.yaml", "Path to config file")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 日志输出到文件 + 控制台
	logFile, err := os.OpenFile("cert-checker.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// 初始化 cron
	c := cron.New(cron.WithSeconds())

	ntf := notifier.NewNotifier(&cfg.Notifiers.Email, &cfg.Notifiers.DingTalk, &cfg.Notifiers.WeCom, &cfg.Notifiers.Bark)
	// 定义检查函数
	checkDomains := func() {
		domains := getDomains(cfg)
		for _, domain := range domains {
			info, err := checker.CheckCert(domain, cfg.Alert.Threshold)
			if err != nil {
				log.Printf("检查域名 %s 失败: %v", domain, err)
				continue
			}

			log.Printf("域名: %-30s 过期时间: %s 剩余天数: %d\n",
				domain,
				info.ExpiryDate.Format("2006-01-02"),
				info.ExpiresIn)

			if info.IsExpired || info.IsWarning {
				msg := notifier.AlertMessage{
					Domain:     domain,
					ExpiryDate: info.ExpiryDate.Format("2006-01-02 15:04:05"),
					DaysLeft:   info.ExpiresIn,
				}
				if err := ntf.Send(msg); err != nil {
					log.Printf("发送通知失败: %v", err)
				}
			}
		}
	}

	// 启动 cron 定时任务
	_, _ = c.AddFunc("0 0 10 * * ?", checkDomains)
	c.Start()
	log.Println("证书监控服务已启动...")

	// 程序启动时立即执行一次
	go checkDomains()

	// 等待退出信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigs
		log.Printf("接收到信号 %v，准备退出...", sig)
		c.Stop() // 停止 cron
	}()

	wg.Wait()
	log.Println("证书监控服务已退出。")
}

func getDomains(cfg *config.Config) []string {
	domainConfig := cfg.DomainConfig
	enabledDomainUrl := domainConfig.EnabledDomainUrl
	if enabledDomainUrl {
		return domainConfig.List
	}
	url := domainConfig.DomainUrl
	return getDomainsByUrl(url)
}

func getDomainsByUrl(url string) []string {
	// POST 参数
	payload := map[string]interface{}{
		"key1": "value1",
	}

	// 将 payload 转为 JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// 创建 POST 请求
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 读取响应 body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("body: %v\n", body)

	// 解析 JSON
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	return result.Data
}

type Response struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data []string `json:"data"`
}
