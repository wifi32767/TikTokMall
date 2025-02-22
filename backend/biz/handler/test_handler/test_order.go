package test_handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wifi32767/TikTokMall/backend/rpc"
	"github.com/wifi32767/TikTokMall/backend/utils"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
)

type updateOrderStateReq struct {
	UserId  uint32 `json:"user_id"`
	OrderId string `json:"order_id"`
	State   string `json:"state"`
}

// @Summary		更新订单状态
// @Description	更新订单状态
// @Tags			test_order
// @Accept			json
// @Produce		json
// @Param			req	body	updateOrderStateReq	true	"请求参数"
// @Success		200
// @Failure		400
// @Failure		500
// @Router			/test/order/update [GET]
func UpdateOrderState(c *gin.Context) {
	var req updateOrderStateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	_, err := rpc.OrderClient.UpdateOrderState(c, &order.UpdateOrderStateReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		State:   req.State,
	})
	if err != nil {
		utils.ErrorHandle(c, err)
		return
	}
	c.Status(http.StatusOK)
}
