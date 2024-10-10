package subinfo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	units = "KMGTPE"
	unit  = 1024
)

var (
	invalidSubscribeUserInfoError = errors.New("invalid subscribe user info ")
)

// SubscriptionInfo 结构体用于存储解析后的信息
type SubscriptionInfo struct {
	Upload   int64
	Download int64
	Total    int64
	Expire   time.Time
}

// ParseSubscriptionInfo parse Subscription-Userinfo
func ParseSubscriptionInfo(info string) (*SubscriptionInfo, error) {
	result := &SubscriptionInfo{}
	pairs := strings.Split(info, ";")
	if len(pairs) == 0 {
		return nil, invalidSubscribeUserInfoError
	}
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])

		switch key {
		case "upload":
			result.Upload, _ = strconv.ParseInt(value, 10, 64)
		case "download":
			result.Download, _ = strconv.ParseInt(value, 10, 64)
		case "total":
			result.Total, _ = strconv.ParseInt(value, 10, 64)
		case "expire":
			expireTime, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				result.Expire = time.Unix(expireTime, 0)
			}
		}
	}
	return result, nil
}

func formatBytes(bytes int64) string {
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), units[exp])
}

// String impl Stringer to display
func (si *SubscriptionInfo) String() string {
	return fmt.Sprintf("ℹ Upload: %s Download: %s Total: %s ExpireAt: %s",
		formatBytes(si.Upload),
		formatBytes(si.Download),
		formatBytes(si.Total),
		si.Expire.Format("2006-01-02 15:04"))
}
