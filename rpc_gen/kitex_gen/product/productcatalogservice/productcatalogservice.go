// Code generated by Kitex v0.9.1. DO NOT EDIT.

package productcatalogservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	product "gomall/rpc_gen/kitex_gen/product"
	proto "google.golang.org/protobuf/proto"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"ListProducts": kitex.NewMethodInfo(
		listProductsHandler,
		newListProductsArgs,
		newListProductsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"GetProduct": kitex.NewMethodInfo(
		getProductHandler,
		newGetProductArgs,
		newGetProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"SearchProducts": kitex.NewMethodInfo(
		searchProductsHandler,
		newSearchProductsArgs,
		newSearchProductsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"UpdateProduct": kitex.NewMethodInfo(
		updateProductHandler,
		newUpdateProductArgs,
		newUpdateProductResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
	"ListProductIds": kitex.NewMethodInfo(
		listProductIdsHandler,
		newListProductIdsArgs,
		newListProductIdsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingUnary),
	),
}

var (
	productCatalogServiceServiceInfo                = NewServiceInfo()
	productCatalogServiceServiceInfoForClient       = NewServiceInfoForClient()
	productCatalogServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return productCatalogServiceServiceInfo
}

// for client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return productCatalogServiceServiceInfoForStreamClient
}

// for stream client
func serviceInfoForClient() *kitex.ServiceInfo {
	return productCatalogServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "ProductCatalogService"
	handlerType := (*product.ProductCatalogService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "product",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.9.1",
		Extra:           extra,
	}
	return svcInfo
}

func listProductsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.ListProductsReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.ProductCatalogService).ListProducts(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *ListProductsArgs:
		success, err := handler.(product.ProductCatalogService).ListProducts(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ListProductsResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newListProductsArgs() interface{} {
	return &ListProductsArgs{}
}

func newListProductsResult() interface{} {
	return &ListProductsResult{}
}

type ListProductsArgs struct {
	Req *product.ListProductsReq
}

func (p *ListProductsArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.ListProductsReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ListProductsArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ListProductsArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ListProductsArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *ListProductsArgs) Unmarshal(in []byte) error {
	msg := new(product.ListProductsReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ListProductsArgs_Req_DEFAULT *product.ListProductsReq

func (p *ListProductsArgs) GetReq() *product.ListProductsReq {
	if !p.IsSetReq() {
		return ListProductsArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ListProductsArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *ListProductsArgs) GetFirstArgument() interface{} {
	return p.Req
}

type ListProductsResult struct {
	Success *product.ListProductsResp
}

var ListProductsResult_Success_DEFAULT *product.ListProductsResp

func (p *ListProductsResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.ListProductsResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ListProductsResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ListProductsResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ListProductsResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *ListProductsResult) Unmarshal(in []byte) error {
	msg := new(product.ListProductsResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ListProductsResult) GetSuccess() *product.ListProductsResp {
	if !p.IsSetSuccess() {
		return ListProductsResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ListProductsResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.ListProductsResp)
}

func (p *ListProductsResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ListProductsResult) GetResult() interface{} {
	return p.Success
}

func getProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.GetProductReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.ProductCatalogService).GetProduct(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *GetProductArgs:
		success, err := handler.(product.ProductCatalogService).GetProduct(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*GetProductResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newGetProductArgs() interface{} {
	return &GetProductArgs{}
}

func newGetProductResult() interface{} {
	return &GetProductResult{}
}

type GetProductArgs struct {
	Req *product.GetProductReq
}

func (p *GetProductArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.GetProductReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *GetProductArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *GetProductArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *GetProductArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *GetProductArgs) Unmarshal(in []byte) error {
	msg := new(product.GetProductReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var GetProductArgs_Req_DEFAULT *product.GetProductReq

func (p *GetProductArgs) GetReq() *product.GetProductReq {
	if !p.IsSetReq() {
		return GetProductArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *GetProductArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *GetProductArgs) GetFirstArgument() interface{} {
	return p.Req
}

type GetProductResult struct {
	Success *product.GetProductResp
}

var GetProductResult_Success_DEFAULT *product.GetProductResp

func (p *GetProductResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.GetProductResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *GetProductResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *GetProductResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *GetProductResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *GetProductResult) Unmarshal(in []byte) error {
	msg := new(product.GetProductResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *GetProductResult) GetSuccess() *product.GetProductResp {
	if !p.IsSetSuccess() {
		return GetProductResult_Success_DEFAULT
	}
	return p.Success
}

func (p *GetProductResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.GetProductResp)
}

func (p *GetProductResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *GetProductResult) GetResult() interface{} {
	return p.Success
}

func searchProductsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.SearchProductsReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.ProductCatalogService).SearchProducts(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *SearchProductsArgs:
		success, err := handler.(product.ProductCatalogService).SearchProducts(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*SearchProductsResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newSearchProductsArgs() interface{} {
	return &SearchProductsArgs{}
}

func newSearchProductsResult() interface{} {
	return &SearchProductsResult{}
}

type SearchProductsArgs struct {
	Req *product.SearchProductsReq
}

func (p *SearchProductsArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.SearchProductsReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *SearchProductsArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *SearchProductsArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *SearchProductsArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *SearchProductsArgs) Unmarshal(in []byte) error {
	msg := new(product.SearchProductsReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var SearchProductsArgs_Req_DEFAULT *product.SearchProductsReq

func (p *SearchProductsArgs) GetReq() *product.SearchProductsReq {
	if !p.IsSetReq() {
		return SearchProductsArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *SearchProductsArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *SearchProductsArgs) GetFirstArgument() interface{} {
	return p.Req
}

type SearchProductsResult struct {
	Success *product.SearchProductsResp
}

var SearchProductsResult_Success_DEFAULT *product.SearchProductsResp

func (p *SearchProductsResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.SearchProductsResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *SearchProductsResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *SearchProductsResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *SearchProductsResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *SearchProductsResult) Unmarshal(in []byte) error {
	msg := new(product.SearchProductsResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *SearchProductsResult) GetSuccess() *product.SearchProductsResp {
	if !p.IsSetSuccess() {
		return SearchProductsResult_Success_DEFAULT
	}
	return p.Success
}

func (p *SearchProductsResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.SearchProductsResp)
}

func (p *SearchProductsResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *SearchProductsResult) GetResult() interface{} {
	return p.Success
}

func updateProductHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.UpdateProductReq)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.ProductCatalogService).UpdateProduct(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *UpdateProductArgs:
		success, err := handler.(product.ProductCatalogService).UpdateProduct(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*UpdateProductResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newUpdateProductArgs() interface{} {
	return &UpdateProductArgs{}
}

func newUpdateProductResult() interface{} {
	return &UpdateProductResult{}
}

type UpdateProductArgs struct {
	Req *product.UpdateProductReq
}

func (p *UpdateProductArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.UpdateProductReq)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *UpdateProductArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *UpdateProductArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *UpdateProductArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *UpdateProductArgs) Unmarshal(in []byte) error {
	msg := new(product.UpdateProductReq)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var UpdateProductArgs_Req_DEFAULT *product.UpdateProductReq

func (p *UpdateProductArgs) GetReq() *product.UpdateProductReq {
	if !p.IsSetReq() {
		return UpdateProductArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *UpdateProductArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *UpdateProductArgs) GetFirstArgument() interface{} {
	return p.Req
}

type UpdateProductResult struct {
	Success *product.UpdateProductResp
}

var UpdateProductResult_Success_DEFAULT *product.UpdateProductResp

func (p *UpdateProductResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.UpdateProductResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *UpdateProductResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *UpdateProductResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *UpdateProductResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *UpdateProductResult) Unmarshal(in []byte) error {
	msg := new(product.UpdateProductResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *UpdateProductResult) GetSuccess() *product.UpdateProductResp {
	if !p.IsSetSuccess() {
		return UpdateProductResult_Success_DEFAULT
	}
	return p.Success
}

func (p *UpdateProductResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.UpdateProductResp)
}

func (p *UpdateProductResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *UpdateProductResult) GetResult() interface{} {
	return p.Success
}

func listProductIdsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(product.Empty)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(product.ProductCatalogService).ListProductIds(ctx, req)
		if err != nil {
			return err
		}
		return st.SendMsg(resp)
	case *ListProductIdsArgs:
		success, err := handler.(product.ProductCatalogService).ListProductIds(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*ListProductIdsResult)
		realResult.Success = success
		return nil
	default:
		return errInvalidMessageType
	}
}
func newListProductIdsArgs() interface{} {
	return &ListProductIdsArgs{}
}

func newListProductIdsResult() interface{} {
	return &ListProductIdsResult{}
}

type ListProductIdsArgs struct {
	Req *product.Empty
}

func (p *ListProductIdsArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(product.Empty)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *ListProductIdsArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *ListProductIdsArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *ListProductIdsArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, nil
	}
	return proto.Marshal(p.Req)
}

func (p *ListProductIdsArgs) Unmarshal(in []byte) error {
	msg := new(product.Empty)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var ListProductIdsArgs_Req_DEFAULT *product.Empty

func (p *ListProductIdsArgs) GetReq() *product.Empty {
	if !p.IsSetReq() {
		return ListProductIdsArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *ListProductIdsArgs) IsSetReq() bool {
	return p.Req != nil
}

func (p *ListProductIdsArgs) GetFirstArgument() interface{} {
	return p.Req
}

type ListProductIdsResult struct {
	Success *product.ListProductIdsResp
}

var ListProductIdsResult_Success_DEFAULT *product.ListProductIdsResp

func (p *ListProductIdsResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(product.ListProductIdsResp)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *ListProductIdsResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *ListProductIdsResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *ListProductIdsResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, nil
	}
	return proto.Marshal(p.Success)
}

func (p *ListProductIdsResult) Unmarshal(in []byte) error {
	msg := new(product.ListProductIdsResp)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *ListProductIdsResult) GetSuccess() *product.ListProductIdsResp {
	if !p.IsSetSuccess() {
		return ListProductIdsResult_Success_DEFAULT
	}
	return p.Success
}

func (p *ListProductIdsResult) SetSuccess(x interface{}) {
	p.Success = x.(*product.ListProductIdsResp)
}

func (p *ListProductIdsResult) IsSetSuccess() bool {
	return p.Success != nil
}

func (p *ListProductIdsResult) GetResult() interface{} {
	return p.Success
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) ListProducts(ctx context.Context, Req *product.ListProductsReq) (r *product.ListProductsResp, err error) {
	var _args ListProductsArgs
	_args.Req = Req
	var _result ListProductsResult
	if err = p.c.Call(ctx, "ListProducts", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetProduct(ctx context.Context, Req *product.GetProductReq) (r *product.GetProductResp, err error) {
	var _args GetProductArgs
	_args.Req = Req
	var _result GetProductResult
	if err = p.c.Call(ctx, "GetProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) SearchProducts(ctx context.Context, Req *product.SearchProductsReq) (r *product.SearchProductsResp, err error) {
	var _args SearchProductsArgs
	_args.Req = Req
	var _result SearchProductsResult
	if err = p.c.Call(ctx, "SearchProducts", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) UpdateProduct(ctx context.Context, Req *product.UpdateProductReq) (r *product.UpdateProductResp, err error) {
	var _args UpdateProductArgs
	_args.Req = Req
	var _result UpdateProductResult
	if err = p.c.Call(ctx, "UpdateProduct", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListProductIds(ctx context.Context, Req *product.Empty) (r *product.ListProductIdsResp, err error) {
	var _args ListProductIdsArgs
	_args.Req = Req
	var _result ListProductIdsResult
	if err = p.c.Call(ctx, "ListProductIds", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
