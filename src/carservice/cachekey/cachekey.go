package cachekey

import (
	"time"
	"gitlab.dev.daikuan.com/platform/servicecenter/lib/go/servicecenter/pkg/cache"
)

var (
	CarBaseInfo                       cache.CacheSettingInfo
	DicCarSerialWhiteBgImgUrls        cache.CacheSettingInfo
	CarSerialWhiteBgImgUrl            cache.CacheSettingInfo
	MatchingSkinRuleInfo              cache.CacheSettingInfo
	MatchingSkinRuleCarRelationIdList cache.CacheSettingInfo
	GetCarMBrandAndInitials           cache.CacheSettingInfo
	CarSerialListByMBrand cache.CacheSettingInfo
)

func init() {
	CarBaseInfo = cache.CacheSettingInfo{CacheKey: "Go:CarBaseInfo:%d", ExpiredSeconds: 2 * time.Hour}
	DicCarSerialWhiteBgImgUrls = cache.CacheSettingInfo{CacheKey: "Go:DicCarSerialWhiteBgImgUrls", ExpiredSeconds: 2 * time.Hour}
	CarSerialWhiteBgImgUrl = cache.CacheSettingInfo{CacheKey: "Go:CarSerialWhiteBgImgUrl:SerialId:%d", ExpiredSeconds: 2 * time.Hour}
	MatchingSkinRuleInfo = cache.CacheSettingInfo{CacheKey: "Go:MatchingSkinRuleInfo:From:%s", ExpiredSeconds: 2 * time.Hour}
	MatchingSkinRuleCarRelationIdList = cache.CacheSettingInfo{CacheKey: "Go:MatchingSkinRuleCarRelationIdList:Level_%d:skinRuleId_%d", ExpiredSeconds: 2 * time.Hour}
	GetCarMBrandAndInitials = cache.CacheSettingInfo{CacheKey: "Go:GetCarMBrandAndInitials:Sale:%t", ExpiredSeconds: 2 * time.Hour}
	CarSerialListByMBrand = cache.CacheSettingInfo{CacheKey: "Go:CarSerialListByMBrand:MbId_%d:Sale_%t", ExpiredSeconds: 2 * time.Hour}
}
