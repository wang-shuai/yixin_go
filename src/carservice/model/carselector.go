package model

type ViewCarSelectorInfo struct {
	/// 分类名称
	CategoryName string

	/// 选车控件信息项集合信息
	CategoryCollection []ViewCarSelectorItemInfo

	/// 分类ID
	CategoryId int
}

/// 选车控件信息项集合信息
type ViewCarSelectorItemInfo struct {
	/// 分类依据
	CategoryBasis string

	/// 国别ID（90为中国）
	CountryId int

	Id int

	/// 名称
	Name string

	/// 拼写
	Spell string

	/// 图片地址
	ImgUrl string

	/// 价格信息
	Price string

	/// 车价文本（单位：万元）
	CarPriceText string

	/// 车价前文本（例如：厂家指导价、易车商城价格、地区经销商平均报价、某某经销商报价等）
	TextBeforeCarPrice string

	/// 车款全称

	ShowCarName string

	/// 分类依据ID
	CategoryBasisId int

	/// 标签文本
	TagText string
}
