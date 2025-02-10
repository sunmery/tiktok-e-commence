package biz

type PlaceOrderReq struct {
	Name         string  // 订单名称
	UserCurrency string  // 货币类型
	Address      Address // 地址信息
	Items        []Item  // 商品列表
	Email        string  // 用户邮箱
	Owner        string  // 订单所有者
}

type Address struct {
	StreetAddress string // 街道地址
	City          string // 城市
	State         string // 州/省
	Country       string // 国家
	ZipCode       string // 邮政编码
}

type Item struct {
	Id          int32   // 商品ID
	Name        string  // 商品名称
	Description string  // 商品描述
	Price       float32 // 商品单价
	Quantity    int32   // 商品数量
}

type PlaceOrderResp struct {
	Order OrderResult // 订单结果
}

type OrderResult struct {
	OrderId int32 // 订单ID
}

type ListOrderReq struct {
	UserId string // 用户ID，用于查询该用户的所有订单
}

type ListOrderResp struct {
	Orders []OrderSummary // 订单列表
}

type OrderSummary struct {
	OrderId      string  // 订单ID
	OrderName    string  // 订单名称
	OrderStatus  string  // 订单状态: pending, paid, shipped, etc.
	CreatedAt    int32   // 创建时间，使用字符串或时间戳，具体依据需求
	TotalAmount  float32 // 订单总金额
	Address      Address // 确保这里有 Address 字段
	UserId       uint32  // 订单所有者的用户ID
	UserCurrency string  // 货币类型
	Email        string  // 用户邮箱
	State        string  // 订单状态
	OrderItems   []OrderItem
}

type OrderItem struct {
	Id          int32   // 商品ID
	Name        string  // 商品名称
	Description string  // 商品描述
	Price       float32 // 商品单价
	Quantity    int32   // 商品数量
}
