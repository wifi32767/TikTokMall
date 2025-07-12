package handler

// basic

type errorReturn struct {
	Error string
}

// user

type userIdReq struct {
	UserId uint32 `json:"user_id" binding:"required"`
}

type registerReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type updateReq struct {
	Username    string `json:"username" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type grantReq struct {
	UserId     uint32 `json:"user_id" binding:"required"`
	Permission uint32 `json:"permission" binding:"required"`
}

// product

type productCreateReq struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Picture     string   `json:"picture" binding:"required"`
	Price       float32  `json:"price" binding:"required"`
	Categories  []string `json:"categories"`
}

type simpleProduct struct {
	ID          uint32   `json:"id" binding:"required"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Picture     string   `json:"picture"`
	Price       float32  `json:"price"`
	Categories  []string `json:"categories"`
}

type idReq struct {
	ID uint32 `json:"id" binding:"required"`
}

type productListReq struct {
	PageSize     int32  `json:"page_size" binding:"required"`
	Page         int32  `json:"page" binding:"required"`
	CategoryName string `json:"category_name"`
}

type searchReq struct {
	Query string `json:"query" binding:"required"`
}

// order

type address struct {
	StreetAddress string `json:"street_address" binding:"required"`
	City          string `json:"city" binding:"required"`
	State         string `json:"state" binding:"required"`
	Country       string `json:"country" binding:"required"`
	ZipCode       int32  `json:"zip_code" binding:"required"`
}

type item struct {
	ProductId uint32  `json:"product_id" binding:"required"`
	Quantity  int32   `json:"quantity" binding:"required"`
	Cost      float32 `json:"cost" binding:"required"`
}

type orderPlaceReq struct {
	UserId       uint32  `json:"user_id" binding:"required"`
	UserCurrency string  `json:"user_currency" binding:"required"`
	Address      address `json:"address"`
	Items        []item  `json:"items" binding:"required"`
}

type orderIdReq struct {
	UserId  uint32 `json:"user_id" binding:"required"`
	OrderId string `json:"order_id" binding:"required"`
}

type orderDetails struct {
	OrderId      string  `json:"order_id"`
	UserId       uint32  `json:"user_id"`
	UserCurrency string  `json:"user_currency"`
	Address      address `json:"address"`
	Email        string  `json:"email"`
	CreatedAt    int32   `json:"created_at"` // 时间戳
}

type orderListRes struct {
	Orders []orderDetails `json:"orders" binding:"required"`
}

// cart

type cartItem struct {
	ProductId uint32 `json:"product_id" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required"`
}

type addItemReq struct {
	UserId uint32   `json:"userid" binding:"required"`
	Item   cartItem `json:"item" binding:"required"`
}

// checkout

type checkoutReq struct {
	UserId     uint32  `json:"user_id" binding:"required"`
	Email      string  `json:"email" binding:"required"`
	Address    address `json:"address"`
	CreditCard *struct {
		CreditCardNumber          string `json:"credit_card_number" binding:"required"`
		CreditCardCvv             int32  `json:"credit_card_cvv" binding:"required"`
		CreditCardExpirationMonth int32  `json:"credit_card_expiration_month" binding:"required"`
		CreditCardExpirationYear  int32  `json:"credit_card_expiration_year" binding:"required"`
	} `json:"credit_card"`
}

type checkoutResp struct {
	OrderId       string `json:"order_id"`
	TransactionId string `json:"transaction_id"`
}
