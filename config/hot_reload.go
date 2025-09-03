// hot_reload.go
// 配置热更新实现，支持动态调整参数
// 通过监听config.yaml变化，实时更新配置

package config

import (
	"fmt"
	"time"
)

var ConfigHandlers = make(map[string]func(string) error)

// RegisterConfigHandler 注册配置处理函数
// 参数：path 配置路径，handler 处理函数
// 返回：无
func RegisterConfigHandler(path string, handler func(string) error) {
	ConfigHandlers[path] = handler
}

// InitConfigHandlers 初始化配置处理逻辑
// 参数：无
// 返回：错误信息
func InitConfigHandlers() error {
	RegisterConfigHandler("monitor_interval", func(value string) error {
		newInterval, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		// 直接使用全局配置的最新值
		GlobalConfig.MonitorInterval = newInterval.String()
		return nil
	})
	return nil
}

// ReloadConfig 热重载配置核心实现
// 参数：configPath 配置文件路径
// 返回：错误信息
func ReloadConfig(configPath string) error {
	// 加载新配置
	newCfg := &AppConfig{}
	err := newCfg.LoadConfig(configPath)
	if err != nil {
		return err
	}
	// 收集需要触发的配置项值
	values := make(map[string]string)
	values["monitor_interval"] = newCfg.MonitorInterval

	// 原子替换全局配置
	GlobalConfig = newCfg

	// 触发所有注册的处理函数
	for path, handler := range ConfigHandlers {
		value, ok := values[path]
		if !ok {
			continue // 忽略未收集的配置项
		}
		if err := handler(value); err != nil {
			return fmt.Errorf("处理配置项 %s 失败: %v", path, err)
		}
	}

	return nil
}
