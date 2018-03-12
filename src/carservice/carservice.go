package carservice

import (
	"fmt"
	//"gitlab.dev.daikuan.com/platform/utils/go/cache/redis"
	"time"
	"gitlab.dev.daikuan.com/platform/servicecenter/lib/go/servicecenter/pkg/cache"
	"carservice/model"
	"carservice/cachekey"
	"carservice/data"
	"go/types"
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
				return types.Nil{}
			}
			return t
		})
	return car, err
}

func (srv *CarService) GetAllCarSerialWhiteBgImgUrlDictionary() (map[int]string,error) {
	cacheSetting := cache.CacheSettingInfo{
		ExpiredSeconds: cachekey.DicCarSerialWhiteBgImgUrls.ExpiredSeconds,
		CacheKey:       cachekey.DicCarSerialWhiteBgImgUrls.CacheKey}

	ret := make(map[int]string)
	err := cacheReader.Get(
		cacheSetting,
		&ret,
		getAllCarSerialWhiteBgImgUrlDictionary)
	return ret,err
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

func (srv *CarService) GetCarSerialWhiteBgImgUrl(carSerialId, imgSize int, dicCarSerialWhiteBgImgUrls map[int]string) (string,error) {
	var ret string
	if carSerialId <= 0 {
		return ret ,errors.New("车系id错误")
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
			}else{
				return ""
			}
		})
	if err != nil{
		return "",err
	}
	ret = strings.Replace(ret,"{0}",strconv.Itoa(imgSize),-1)
	return ret,err
}

func (srv* CarService) GetCarBaseInfoWithImgSize(carId, imgSize int) (model.CarBaseInfo,error){
	var err error
	car,err := srv.GetCarBaseInfo(carId)
	if err!=nil{
		return car,err
	}
	imgUrl,err := srv.GetCarSerialWhiteBgImgUrl(car.CarSerialId,imgSize,nil)
	if err!=nil{
		fmt.Println("获取车系图片失败","serival:",car.CarSerialId)
	}
	car.CarSerialImgUrl = imgUrl
	return car,nil
}

