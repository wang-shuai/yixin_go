package data

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"
	"carservice/model"
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
	engine.ShowSQL(true)                //展示每次执行的sql
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

func GetCarMBrandAndInitials(sale bool) ([]model.ViewCarSelectorItemInfo, error) {
	cb_cond := "b.CbSaleState = '在销'"
	cs_cond := "s.CsSaleState = '在销'"
	c_cond := "c.car_SaleState = 95"
	if !sale {
		cb_cond = "(b.CbSaleState = '在销' OR b.CbSaleState = '停销')"
		cs_cond = "(s.CsSaleState = '在销' OR s.CsSaleState = '停销')"
		c_cond = "( c.car_SaleState = 95 OR c.car_SaleState = 96 )"
	}

	result := make([]model.ViewCarSelectorItemInfo, 0)
	err := engine.SQL(fmt.Sprintf(`SELECT m.bs_Id AS Id, m.bs_Name AS Name, SUBSTRING(m.spell,1,1) AS CategoryBasis, m.urlspell AS Spell
                            FROM Car_MasterBrand AS m WITH (NOLOCK)
                            WHERE m.IsState = 0 AND m.IsLock = 0 AND m.bs_Id IN
                            (
	                            SELECT DISTINCT r.bs_Id
	                            FROM Car_MasterBrand_Rel AS r WITH (NOLOCK)
	                            INNER JOIN Car_Brand AS b WITH (NOLOCK) ON r.cb_Id = b.cb_Id
	                            INNER JOIN Car_Serial AS s WITH (NOLOCK) ON b.cb_Id = s.cb_Id
	                            INNER JOIN Car_relation AS c WITH (NOLOCK) ON s.cs_Id = c.Cs_Id
	                            WHERE r.IsState = 0
	                            AND b.IsState = 0 AND b.IsLock = 0
	                            AND %[1]s
	                            AND s.IsState = 0 AND s.IsLock = 0
	                            AND %[2]s
	                            AND c.IsState = 0 AND c.IsLock = 0 AND c.car_ReferPrice IS NOT NULL AND c.Car_YearType IS NOT NULL
	                            AND %[3]s
                            )
                            ORDER BY CategoryBasis ASC, m.bs_Id ASC`, cb_cond, cs_cond, c_cond)).Find(&result)
	return result, err
}

func CarSerialListByMBrand(mbid int,sale bool) ([]model.ViewCarSelectorItemInfo, error) {
	cb_cond := "b.CbSaleState = '在销'"
	cs_cond := "s.CsSaleState = '在销'"
	c_cond := "c.car_SaleState = 95"
	if !sale {
		cb_cond = "(b.CbSaleState = '在销' OR b.CbSaleState = '停销')"
		cs_cond = "(s.CsSaleState = '在销' OR s.CsSaleState = '停销')"
		c_cond = "( c.car_SaleState = 95 OR c.car_SaleState = 96 )"
	}

	result := make([]model.ViewCarSelectorItemInfo, 0)
	err := engine.SQL(fmt.Sprintf(`;WITH w AS
                            (
	                            SELECT s.cs_Id, MIN(c.car_ReferPrice) AS PriceMin, MAX(c.car_ReferPrice) AS PriceMax
	                            FROM Car_MasterBrand_Rel AS r WITH (NOLOCK)
	                            INNER JOIN Car_Brand AS b WITH (NOLOCK) ON r.cb_Id = b.cb_Id
	                            INNER JOIN Car_Serial AS s WITH (NOLOCK) ON b.cb_Id = s.cb_Id
	                            INNER JOIN Car_relation AS c WITH (NOLOCK) ON s.cs_Id = c.Cs_Id
	                            WHERE r.bs_Id = %[1]d AND r.IsState = 0
	                            AND b.IsState = 0 AND b.IsLock = 0
	                            AND %[2]s
	                            AND s.IsState = 0 AND s.IsLock = 0
	                            AND %[3]s
	                            AND c.IsState = 0 AND c.IsLock = 0 AND c.car_ReferPrice IS NOT NULL AND c.Car_YearType IS NOT NULL
	                            AND %[4]s
	                            GROUP BY s.cs_Id
                            )
                            SELECT b.cb_country AS CountryId, b.cb_Name AS CategoryBasis, b.cb_Id AS CategoryBasisId, s.cs_Id AS Id, s.csShowName AS Name, s.allSpell AS Spell, CONVERT(varchar(50), CONVERT(varchar(20), w.PriceMin) + '~' + CONVERT(varchar(20), w.PriceMax) + '万') AS Price
                            FROM w
                            INNER JOIN Car_Serial AS s ON w.cs_Id = s.cs_Id
                            INNER JOIN Car_Brand AS b ON s.cb_Id = b.cb_Id
                            ORDER BY b.cb_country ASC, b.cb_Id ASC, s.CsSaleState DESC, s.csShowName ASC`, mbid, cb_cond, cs_cond, c_cond)).
		Find(&result)

	return result, err
}

func GetMatchingSkinRuleInfo(from string) (model.C2BSkinRuleInfo, error) {
	var entity model.C2BSkinRuleInfo
	has, err := engine.Table("C2B_SkinRule").
		Cols("ID,HighlightText").
		Where(`IsDeleted=0 AND IsEnabled=1
					AND ( (? = '' AND IsUsedByNullFrom=1)
						  OR (? <> '' AND ExChannelCode LIKE '%' + ? + '%')
						  OR (? <> '' AND ExChannelCode = '-1') )
					AND HighlightText IS NOT NULL`, from, from, from, from).
		Get(&entity)
	if has {
		fmt.Println(entity)
	}
	return entity, err
}

func GetMatchSkinRuleCarRelationIdList(level, ruleId int) ([]int, error) {
	var idlist []int

	idFieldName := ""

	switch
	{
	case level == 1:
		idFieldName = "bs_Id"
		break
	case level == 2:
		idFieldName = "cb_Id"
		break
	case level == 3:
		idFieldName = "cs_Id"
		break
	case level == 4:
		idFieldName = "Car_Id"
		break
	}

	err := engine.Table("C2B_SkinRule_HighlightedCar").Alias("h").
		Join("Inner", []string{"ViewLevelCar", "c"},
		`? <= h.RelationType AND
				(
					(h.RelationType=1 AND c.bs_Id=h.RelationID) OR
					(h.RelationType=2 AND c.cb_Id=h.RelationID) OR
					(h.RelationType=3 AND c.cs_Id=h.RelationID) OR
					(h.RelationType=4 AND c.Car_Id=h.RelationID)
				)`, level).
		Where("h.IsDeleted=0 and h.SkinRuleID = ? ", ruleId).
		Distinct("c." + idFieldName).
		Find(&idlist)

	return idlist, err
}
