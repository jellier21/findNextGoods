package main

import (
	"GoFindNextHomeObj/Server-web/Server-Side/Goods"
	Self_made "GoFindNextHomeObj/Server-web/Server-Side/Self-made"
	"GoFindNextHomeObj/Server-web/Server-Side/User"
	User_Location "GoFindNextHomeObj/Server-web/Server-Side/User-Location"
	"GoFindNextHomeObj/Server-web/Server-Side/UserNeed"
	image_op "GoFindNextHomeObj/Server-web/image-op"
	"fmt"
	"net/http"
)

func main() {
	User.UseTheUserServer()
	User_Location.UseTheSiteServer()
	Self_made.UseSelfMadeApiServer()
	UserNeed.UseNeedApiServer()
	Goods.UseGoodsApiServer()
	image_op.UseTheImageFileApiServer()
	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
