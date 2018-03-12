package data

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	//"time"
	"carservice/model"
	//"net/url"
	"gitlab.dev.daikuan.com/platform/servicecenter/lib/go/servicecenter/pkg/config"
	"github.com/tidwall/gjson"
)

var (
	engine *xorm.Engine
)

func init() {
	//query := url.Values{}
	//query.Add("connection timeout", fmt.Sprintf("%d", 180))
	//query.Add("database", "Finance")

	//u := &url.URL{
	//	Scheme:   "sqlserver",
	//	User:     url.UserPassword("sa", "1234.abcd"),
	//	Host:     fmt.Sprintf("%s", "192.168.151.56"),
	//	Path:     "SERVER14", // if connecting to an instance instead of a port
	//	RawQuery: query.Encode(),
	//}
	//
	//connectionString := u.String()
	//
	//if Eg, err := xorm.NewEngine("mssql", connectionString); err != nil {
	//	fmt.Println(err)
	//	panic("数据库链接失败")
	//} else {
	//	engine = Eg
	//}

	schema := gjson.Get(config.AppConfig.Json, "db.schema").String()
	connStr := gjson.Get(config.AppConfig.Json, "db.connectionString").String()

	// 看源码发现 只支持 mssql配ODBC 格式的链接字符串
	//if Eg, err := xorm.NewEngine("mssql", "odbc:server=192.168.151.56\\SERVER14;user id=sa;password={1234.abcd};database=Finance;connection timeout=30"); err != nil {
	if Eg, err := xorm.NewEngine(schema, connStr); err != nil {
		fmt.Println(err)
		panic("数据库链接失败")
	} else {
		engine = Eg
	}

	engine.SetMapper(core.SameMapper{}) //与字段、表名一致  不区分大小写
}

func GetCarBaseInfo(carId int) (car model.CarBaseInfo, err error) {
	has, err := engine.Table("Car_relation").Alias("c").
		Select(`c.Car_Id AS CarId, LTRIM(RTRIM(c.Car_Name)) AS CarName, c.cs_Id AS CarSerialId, LTRIM(RTRIM(s.csName)) AS CarSerialName,
		s.cb_Id AS CarBrandId, LTRIM(RTRIM(b.cb_Name)) AS CarBrandName, m.bs_Id AS CarMasterBrandId, LTRIM(RTRIM(m.bs_Name)) AS CarMasterBrandName,
		c.Car_YearType AS CarYear, c.car_ReferPrice AS CarReferPrice, s.allSpell AS carserialallspell, s.csShowName AS carserialshowNAME `).
		Join("inner", []string{"Car_Serial", "s"}, "c.Cs_Id = s.cs_Id").
		Join("inner", []string{"Car_Brand", "b"}, "s.cb_Id = b.cb_Id").
		Join("inner", []string{"Car_MasterBrand_Rel", "r"}, "b.cb_Id = r.cb_Id AND r.IsState = 0").
		Join("inner", []string{"Car_MasterBrand", "m"}, "r.bs_Id = m.bs_Id").
		Where("c.Car_Id = ?", carId).Get(&car)

	if has {
		fmt.Println(car)
	}
	return
}
