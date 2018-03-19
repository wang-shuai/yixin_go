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

func TestCarService_GetCarMBrandAndInitials(t *testing.T) {
	src := new(CarService)
	data,err:= src.GetCarMBrandGroupByInitials(true)
	fmt.Println(data,err)
}

func TestCarService_GetCarSelectorMasterBrandList(t *testing.T) {
	srv := new (CarService)
	data,err:= srv.GetCarSelectorMasterBrandList(true,true,"1008")
	fmt.Println(data,"\n\r",err)
}

func TestCarService_GetCarSerialListWithMBrand(t *testing.T) {
	srv := new (CarService)
	data,err:= srv.GetCarSerialListWithMBrand(5,true,true,"1008")
	fmt.Println(data,"\n\r",err)
}