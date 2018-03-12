package cachekey

import (
	"time"
	"gitlab.dev.daikuan.com/platform/servicecenter/lib/go/servicecenter/pkg/cache"
)

var (
	CarBaseInfo cache.CacheSettingInfo
	DicCarSerialWhiteBgImgUrls cache.CacheSettingInfo
	CarSerialWhiteBgImgUrl cache.CacheSettingInfo
)

func init() {
	CarBaseInfo = cache.CacheSettingInfo{CacheKey: "Go:CarBaseInfo:%d", ExpiredSeconds: 2 * time.Hour}
	DicCarSerialWhiteBgImgUrls = cache.CacheSettingInfo{CacheKey: "Go:DicCarSerialWhiteBgImgUrls", ExpiredSeconds: 2 * time.Hour}
	CarSerialWhiteBgImgUrl = cache.CacheSettingInfo{CacheKey:"Go:CarSerialWhiteBgImgUrl:SerialId:%d",ExpiredSeconds:2*time.Hour}
}
