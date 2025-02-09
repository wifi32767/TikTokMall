.PHONY: gen-auth
gen-auth:
	@cd auth && kitex -service authservice -module "github.com/wifi32767/TikTokMall/app/auth" -I ../../idl -use "github.com/wifi32767/TikTokMall/rpc/kitex_gen" ../../idl/auth.proto

.PHONY: gen-user
gen-user:
	@cd user && kitex -service userservice -module "github.com/wifi32767/TikTokMall/app/user" -I ../../idl -use "github.com/wifi32767/TikTokMall/rpc/kitex_gen" ../../idl/user.proto

.PHONY: gen-product
gen-product:
	@cd product && kitex -service productservice -module "github.com/wifi32767/TikTokMall/app/product" -I ../../idl -use "github.com/wifi32767/TikTokMall/rpc/kitex_gen" ../../idl/product.proto

.PHONY: gen-cart
gen-cart:
	@cd cart && kitex -service cartservice -module "github.com/wifi32767/TikTokMall/app/cart" -I ../../idl -use "github.com/wifi32767/TikTokMall/rpc/kitex_gen" ../../idl/cart.proto