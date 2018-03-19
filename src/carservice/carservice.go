package carservice

import (
	"fmt"
	//"gitlab.dev.daikuan.com/platform/utils/go/cache/redis"
	"time"
	"gitlab.dev.daikuan.com/platform/servicecenter/lib/go/servicecenter/pkg/cache"
	"carservice/model"
	"carservice/cachekey"
	"carservice/data"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"errors"
	"strings"
)

// test region
type ttt struct {
	Name string
	Year int
}

func (car *CarService) Test() {
	//// 在redis-cli没获取到数据、有点晕
	//rds := redis.NewRedisCache()
	//rds.StartAndGC(`{"conn":"192.168.145.3:6379"}`)
	//expire := 10 * time.Minute
	//if err := rds.Put("test", []byte("shine"), expire); err != nil {
	//	fmt.Errorf("redis set data err:", err)
	//}
	//data := rds.Get("test")
	//switch s := data.(type) {
	//case []byte:
	//	fmt.Println("redis key data :", string(s))
	//	break
	//default:
	//	fmt.Println(s)
	//}

	expire := 10 * time.Minute
	reader := new(cache.MultipleCacheReader)
	cacheSetting := cache.CacheSettingInfo{ExpiredSeconds: expire, CacheKey: "reader"}
	var v ttt
	if err := reader.Get(cacheSetting, &v, getScourceData); err != nil {
		fmt.Errorf("cache read error ")
	} else {
		fmt.Println("cache reader get data :", v)
	}
	fmt.Println(v.Name, v.Year)
}
func getScourceData() interface{} {
	return &ttt{Name: "迈腾", Year: 2018}
}

// test region end

type CarService struct {
}

// 内部变量
var (
	cacheReader *cache.MultipleCacheReader
)

func init() {
	cacheReader = new(cache.MultipleCacheReader)
}

// 获取车款基础信息
func (src *CarService) GetCarBaseInfo(carId int) (car model.CarBaseInfo, err error) {
	cacheSetting := cache.CacheSettingInfo{
		ExpiredSeconds: cachekey.CarBaseInfo.ExpiredSeconds,
		CacheKey:       fmt.Sprintf(cachekey.CarBaseInfo.CacheKey, carId)}

	err = cacheReader.Get(
		cacheSetting,
		&car,
		func() interface{} {
			t, e := data.GetCarBaseInfo(carId)
			if e != nil {
				return nil
			}
			return t
		})
	return car, err
}

//获取所有车系的白底图
func (srv *CarService) GetAllCarSerialWhiteBgImgUrlDictionary() (map[int]string, error) {
	cacheSetting := cache.CacheSettingInfo{
		ExpiredSeconds: cachekey.DicCarSerialWhiteBgImgUrls.ExpiredSeconds,
		CacheKey:       cachekey.DicCarSerialWhiteBgImgUrls.CacheKey}

	ret := make(map[int]string)
	err := cacheReader.Get(
		cacheSetting,
		&ret,
		getAllCarSerialWhiteBgImgUrlDictionary)
	return ret, err
}

func getAllCarSerialWhiteBgImgUrlDictionary() interface{} {
	ret := make(map[int]string)

	if resp, err := http.Get("http://webapi.photo.bitauto.com/photoApi/capi/yx/v1/model/getcoverlist?Cache=true"); err != nil {
		fmt.Println(err)
		return nil
	} else {
		defer resp.Body.Close()

		if data, err := ioutil.ReadAll(resp.Body); err != nil {
			fmt.Println(err)
			return nil
		} else {
			mp := make(map[string]interface{})
			if err := json.Unmarshal(data, &mp); err == nil {
				datalist := make([]map[string]interface{}, 0)
				tempdata, _ := json.Marshal(mp["Data"])
				json.Unmarshal(tempdata, &datalist)
				//fmt.Println(datalist)
				for _, v := range datalist {
					tid, _ := v["ModelId"]
					turl, _ := v["WhiteCoverUrl"]
					serialId, _ := strconv.Atoi(fmt.Sprintf("%v", tid))
					url := fmt.Sprintf("%s", turl)
					//fmt.Println(id, url)
					if serialId > 0 && len(url) > 0 && ret[serialId] == "" {
						ret[serialId] = url
					}
				}
			}
			//fmt.Println("获取到车系图片数量",mp)
		}
	}
	return ret
}

// 获取某车系的白底图


func (srv *CarService) GetCarSerialWhiteBgImgUrl(carSerialId, imgSize int, dicCarSerialWhiteBgImgUrls map[int]string) (string, error) {
	var ret string
	if carSerialId <= 0 {
		return ret, errors.New("车系id错误")
	}
	cacheSetting := cache.CacheSettingInfo{
		ExpiredSeconds: cachekey.CarSerialWhiteBgImgUrl.ExpiredSeconds,
		CacheKey:       fmt.Sprintf(cachekey.CarSerialWhiteBgImgUrl.CacheKey, carSerialId)}
	err := cacheReader.Get(
		cacheSetting,
		&ret,
		func() interface{} {
			if len(dicCarSerialWhiteBgImgUrls) == 0 {
				dicCarSerialWhiteBgImgUrls,_ = srv.GetAllCarSerialWhiteBgImgUrlDictionary()
			}
			imgUrl, ok := dicCarSerialWhiteBgImgUrls[carSerialId]
			if ok {
				return imgUrl
			} else {
				return ""
			}
		})
	if err != nil {
		return "", err
	}
	ret = strings.Replace(ret, "{0}", strconv.Itoa(imgSize), -1)
	return ret, err

}

// 获取带白底图的车款信息
func (srv *CarService) GetCarBaseInfoWithImgSize(carId, imgSize int) (model.CarBaseInfo, error) {
	var err error
	car, err := srv.GetCarBaseInfo(carId)
	if err != nil {
		return car, err
	}
	imgUrl, err := srv.GetCarSerialWhiteBgImgUrl(car.CarSerialId, imgSize, nil)
	if err != nil {
		fmt.Println("获取车系图片失败", "serival:", car.CarSerialId)
	}
	car.CarSerialImgUrl = imgUrl
	return car, nil
}

// 获取主品牌logo
func (srv *CarService) GetCarMasterBrandLogoUrl(masterBrandId, imgSize int) string {
	return getCarMasterBrandLogoUrl(masterBrandId, imgSize)
}
func getCarMasterBrandLogoUrl(masterBrandId, imgSize int) string {
	return fmt.Sprintf("http://image.bitautoimg.com/bt/car/default/images/logo/masterbrand/png/%[2]d/m_%[1]d_%[2]d.png", masterBrandId, imgSize)
}

// 获取按字母分组品牌信息 --  CarService.GetCarSelectorMasterBrandList 可以覆盖该方法
func (srv *CarService) GetCarMBrandGroupByInitials(sale bool) ([]model.ViewCarSelectorInfo, error) {
	cacheSetting := cache.CacheSettingInfo{
		CacheKey:       fmt.Sprintf(cachekey.GetCarMBrandAndInitials.CacheKey, sale),
		ExpiredSeconds: cachekey.GetCarMBrandAndInitials.ExpiredSeconds}

	var result []model.ViewCarSelectorInfo
	err := cacheReader.Get(
		cacheSetting,
		&result,
		func() interface{} {
			t, e := getCarMBrandGroupByInitials(sale)
			if e != nil {
				return nil
			}
			return t
		})
	if err != nil {
		return nil, err
	}
	return result, err
}

// 获取选车控件 按字母分组 带皮肤标识的主品牌信息
func (srv *CarService) GetCarSelectorMasterBrandList(onlyOnSale bool, needTag bool, from string) ([]model.ViewCarSelectorInfo, error) {
	mbGroup, err := srv.GetCarMBrandGroupByInitials(onlyOnSale)
	if err != nil {
		return nil, err
	}
	if needTag {
		handleCarSelectorSkinTag(1,mbGroup,from)
	}
	return mbGroup, err
}

func getCarMBrandGroupByInitials(sale bool) ([]*model.ViewCarSelectorInfo, error) {
	tData, err := data.GetCarMBrandAndInitials(sale)
	if err != nil {
		return nil, err
	}
	result := make([]*model.ViewCarSelectorInfo, 0)
	exist := make(map[string]int)
	initials := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i, key := range initials {
		var selector *model.ViewCarSelectorInfo
		if idx, ok := exist[key]; ok {
			selector = result[idx]
		} else {
			selector = &model.ViewCarSelectorInfo{}
			selector.CategoryName = key
			selector.CategoryCollection = make([]model.ViewCarSelectorItemInfo, 0)
			exist[key] = i
			result = append(result, selector)
		}
		//go func(model.ViewCarSelectorInfo,[]model.ViewCarSelectorItemInfo) {
		//
		//}(selector,tData)
		for _, val := range tData {
			if strings.ToUpper(val.CategoryBasis) == key {
				item := model.ViewCarSelectorItemInfo{
					Id:     val.Id,
					Name:   val.Name,
					ImgUrl: getCarMasterBrandLogoUrl(val.Id, model.W100H100),
					Spell:  val.Spell,
				}
				selector.CategoryCollection = append(selector.CategoryCollection, item)
			}
		}
	}

	return result, nil
}

func (srv *CarService) GetCarSerialListWithMBrand(mbid int, sale bool,needTag bool,from string)([]model.ViewCarSelectorInfo, error){
	mbGroup, err := srv.GetCarSerialListGroupByMBrand(mbid,sale)
	if err != nil {
		return nil, err
	}
	if needTag {
		handleCarSelectorSkinTag(3,mbGroup,from)
	}
	return mbGroup, err
}

func (srv *CarService) GetCarSerialListGroupByMBrand(mbid int, sale bool)([]model.ViewCarSelectorInfo, error){
	cacheSetting := cache.CacheSettingInfo{
		CacheKey:       fmt.Sprintf(cachekey.CarSerialListByMBrand.CacheKey, mbid, sale),
		ExpiredSeconds: cachekey.CarSerialListByMBrand.ExpiredSeconds}

	var result []model.ViewCarSelectorInfo
	err := cacheReader.Get(
		cacheSetting,
		&result,
		func() interface{} {
			t, e := getCarSerialListWithMBrand(mbid, sale)
			if e != nil {
				return nil
			}
			return t
		})
	if err != nil {
		return nil, err
	}
	return result, err
}

func getCarSerialListWithMBrand(mbid int, sale bool)([]*model.ViewCarSelectorInfo, error){
	result := make([]*model.ViewCarSelectorInfo,0)

	allserials, err := data.CarSerialListByMBrand(mbid,sale)
	if err != nil {
		return nil, err
	}

	// 存放到map中
	cateMap:= make(map[int][]*model.ViewCarSelectorItemInfo)
	for i := range allserials{
		serial := allserials[i]
		var items []*model.ViewCarSelectorItemInfo
		var ok bool
		if items,ok = cateMap[serial.CategoryBasisId];!ok{
			items = make([]*model.ViewCarSelectorItemInfo,0)
			cateMap[serial.CategoryBasisId] = items
		}
		serial.ImgUrl,_ = new(CarService).GetCarSerialWhiteBgImgUrl(serial.Id,1,nil)
		if len(serial.ImgUrl)==0{
			serial.ImgUrl = "http://img4.yixinfinance.com/taoche/common/images/taoche_default.html.png"
		}
		items = append(items,&serial)
		cateMap[serial.CategoryBasisId] = items //append(items,&serial)
	}
	//map 转成 实体集合列表
	for key := range cateMap{
		sltor := model.ViewCarSelectorInfo{}
		sltor.CategoryId = key
		for i,val := range cateMap[key]{
			sltor.CategoryCollection = append(sltor.CategoryCollection,*val)
			if i==0 {
				sltor.CategoryName = val.CategoryBasis
				if val.CountryId != 90{
					sltor.CategoryName = "进口" + val.CategoryBasis
				}
			}
		}
		result = append(result,&sltor)
	}

	return result,nil
}
func getMatchingSkinRuleInfo(from string) (model.C2BSkinRuleInfo, error) {
	cacheSetting := cache.CacheSettingInfo{
		ExpiredSeconds: cachekey.MatchingSkinRuleInfo.ExpiredSeconds,
		CacheKey:       fmt.Sprintf(cachekey.MatchingSkinRuleInfo.CacheKey, from)}
	var entity model.C2BSkinRuleInfo
	err := cacheReader.Get(cacheSetting, &entity, func() interface{} {
		t, e := data.GetMatchingSkinRuleInfo(from)
		if e != nil {
			return nil
		}
		return t
	})
	return entity, err
}

func getMatchSkinRuleCarRelationIdList(level, ruleId int) ([]int, error) {
	cacheSetting := cache.CacheSettingInfo{
		CacheKey:       fmt.Sprintf(cachekey.MatchingSkinRuleCarRelationIdList.CacheKey, level, ruleId),
		ExpiredSeconds: cachekey.MatchingSkinRuleCarRelationIdList.ExpiredSeconds}
	result := make([]int, 0)
	err := cacheReader.Get(cacheSetting,
		&result,
		func() interface{} {
			t, e := data.GetMatchSkinRuleCarRelationIdList(level, ruleId)
			if e != nil {
				return nil
			}
			return t
		})
	if err != nil {
		fmt.Println(err)
	}
	return result, err
}

func handleCarSelectorSkinTag(level int,selector []model.ViewCarSelectorInfo,from string){
	skin, er := getMatchingSkinRuleInfo(from)
	if er != nil {
		fmt.Println("获取皮肤错误，", er)
	} else {
		matchIds, e := getMatchSkinRuleCarRelationIdList(level, skin.Id)
		if e == nil && len(matchIds) > 0 {
			for _, id := range matchIds {
				for idx := range selector {
					for i := range selector[idx].CategoryCollection {
						item := &selector[idx].CategoryCollection[i]
						if item.Id == id {
							item.TagText = skin.HighlightText
						}
					}
				}
			}
		}
	}
}



