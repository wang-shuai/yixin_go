package carservice

import (
	"testing"
	"fmt"
)

func TestCarService_Test(t *testing.T) {
	service := new(CarService)
	service.Test()
}

func TestCarService_GetCarBaseInfo(t *testing.T) {
	service := new(CarService)
	car, _ := service.GetCarBaseInfo(122123)
	fmt.Println(car)
}

func TestCarService_GetAllCarSerialWhiteBgImgUrlDictionary(t *testing.T) {
	service := new(CarService)
	car,_ := service.GetAllCarSerialWhiteBgImgUrlDictionary()
	fmt.Println(car)
}

func TestCarService_GetCarSerialWhiteBgImgUrl(t *testing.T) {
	service := new(CarService)
	imgurl, _ := service.GetCarSerialWhiteBgImgUrl(4610,7,nil)
	fmt.Println(imgurl)
}

func TestCarService_GetCarBaseInfoWithImgSize(t *testing.T) {
	srv:= & CarService{}
	car,err:= srv.GetCarBaseInfoWithImgSize(122123,7)
	fmt.Println(car,err)
}