package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	home "gomall/app/frontend/hertz_gen/frontend/home"
)

type HomeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewHomeService(Context context.Context, RequestContext *app.RequestContext) *HomeService {
	return &HomeService{RequestContext: RequestContext, Context: Context}
}

func (h *HomeService) Run(req *home.Empty) (map[string]any, error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	resp := make(map[string]any)
	items := []map[string]any{
		{"Name": "T-shirt1", "Price": 100, "Picture": "/static/img/t-shirt.jpg"},
		{"Name": "T-shirt2", "Price": 200, "Picture": "/static/img/t-shirt.jpg"},
		{"Name": "T-shirt3", "Price": 300, "Picture": "/static/img/t-shirt.jpg"},
	}
	resp["Title"] = "Hot Sales"
	resp["Items"] = items
	return resp, nil
}
