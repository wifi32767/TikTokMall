package model

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	UserId    uint32
	ProductId uint32
	Quantity  int32
}

type Cart []*CartItem

func (c CartItem) TableName() string {
	return "cart"
}

func AddItem(db *gorm.DB, ctx context.Context, cart *CartItem) error {
	oldCart := CartItem{}
	// 追加下单的时候要在原来记录的基础上修改
	err := db.WithContext(ctx).Where(&CartItem{
		UserId:    cart.UserId,
		ProductId: cart.ProductId,
	},
	).First(&oldCart).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// 商品数量是int，可以为负，表示从购物车中删除物品
		// 因此在商品数量为0的时候删除这条记录
		// 个人觉得一个比较方便的做法是前端不显示商品数量为0的记录
		// 因此这里就没有加入额外的代码
		if cart.Quantity <= 0 {
			return kerrors.NewBizStatusError(http.StatusBadRequest, "quantity can't be negative")
		}
		err = db.WithContext(ctx).Create(cart).Error
	} else {
		if oldCart.Quantity+cart.Quantity < 0 {
			return kerrors.NewBizStatusError(http.StatusBadRequest, "quantity can't be negative")
		}
		err = db.WithContext(ctx).Model(&oldCart).Update("quantity", oldCart.Quantity+cart.Quantity).Error
	}
	return err
}

func EmptyCart(db *gorm.DB, ctx context.Context, userid uint32) error {
	return db.WithContext(ctx).Delete(&CartItem{}, "user_id = ?", userid).Error
}

func GetCart(db *gorm.DB, ctx context.Context, userid uint32) (*Cart, error) {
	res := make(Cart, 0)
	err := db.WithContext(ctx).Where("user_id = ?", userid).Find(&res).Error
	return &res, err
}
