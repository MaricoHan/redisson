package kit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	en2 "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	errors2 "gitlab.bianjie.ai/avata/utils/errors"

	"gitlab.bianjie.ai/avata/open-api/internal/pkg/constant"
	"gitlab.bianjie.ai/avata/open-api/internal/pkg/initialize"
)

type (
	// Endpoint Router define a router for http Handler
	Endpoint struct {
		URI     string
		Method  string
		Handler http.Handler
	}

	Handler            = endpoint.Endpoint
	Server             = httptransport.Server
	RequestFunc        = httptransport.RequestFunc
	ServerResponseFunc = httptransport.ServerResponseFunc

	//IController define a interface for all http Controller
	IController interface {
		GetEndpoints() []Endpoint
	}

	Application interface {
		IController
		//Initialize()
		//Stop()
	}

	Controller struct {
		validate *validator.Validate
	}
)

var trans ut.Translator
var log = initialize.Log

func NewController() Controller {
	validate := validator.New()
	en := en2.New()
	trans, _ = ut.New(en, en).GetTranslator("en")
	entranslations.RegisterDefaultTranslations(validate, trans)
	return Controller{validate}
}

// MakeHandler create a http hander for request
func (c Controller) MakeHandler(handler endpoint.Endpoint, request interface{},
	before []httptransport.RequestFunc,
	mid []httptransport.ServerOption,
	after []httptransport.ServerResponseFunc,
) *httptransport.Server {
	return httptransport.NewServer(
		handler,
		c.decodeRequest(request),
		c.encodeResponse,
		c.serverOptions(before, mid, after)...,
	)
}

func (c Controller) GetIntValue(ctx context.Context, key string) (int, error) {
	value := ctx.Value(key)
	if value == nil {
		return 0, errors.Errorf("Not found key: %s in Context", key)
	}

	v, err := strconv.ParseInt(value.(string), 10, 64)
	if err != nil {
		log.Error("Invalid key, must be int type")
		return 0, errors.Errorf("Value: %s is not int type", value)
	}
	return int(v), nil
}

func (c Controller) GetStringValue(ctx context.Context, key string) (string, error) {
	value := ctx.Value(key)
	if value == nil {
		return "", errors.Errorf("Not found key: %s in Context", key)
	}

	v, ok := value.(string)
	if !ok {
		log.Error("Invalid key, must be string type")
		return "", errors.Errorf("Value: %s is not string type", value)
	}
	return v, nil
}

func (c Controller) GetDateValue(ctx context.Context, key string) (*time.Time, error) {
	value := ctx.Value(key)
	if value == nil {
		return nil, errors.Errorf("Not found key: %s in Context", key)
	}

	tim, err := time.Parse(constant.TimeLayout, value.(string))
	if err != nil {
		log.Error("Invalid key, must be string type")
		return nil, errors.Errorf("Value: %s is not string type", value)
	}
	return &tim, nil
}

func (c Controller) GetPagation(ctx context.Context) (int, int) {
	page, err := c.GetIntValue(ctx, "page")
	if err != nil {
		page = 1
	}

	size, err := c.GetIntValue(ctx, "size")
	if err != nil {
		size = 10
	}
	return page, size
}

//深度克隆，可以克隆任意数据类型
func DeepClone(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()                          //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem()            //创建对象
		b, _ := json.Marshal(src)                 //导出json
		json.Unmarshal(b, dst.Addr().Interface()) //json序列化
		return dst.Addr().Interface()             //返回指针
	} else {
		dst := reflect.New(typ).Elem()            //创建对象
		b, _ := json.Marshal(src)                 //导出json
		json.Unmarshal(b, dst.Addr().Interface()) //json序列化
		return dst.Interface()                    //返回值
	}
}

// decodeRequest decode request(http.request -> model.request)
func (c Controller) decodeRequest(req interface{}) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		if req == nil {
			return nil, err
		}
		p := reflect.ValueOf(req).Elem()
		p.Set(reflect.Zero(p.Type()))

		tmpReq := DeepClone(req)
		if r.Method != constant.Delete {
			if err := json.NewDecoder(r.Body).Decode(&tmpReq); err != nil {
				log.Error("Execute decode request failed,", "error", err.Error())
				return nil, errors2.New(errors2.ClientParams, constant.ErrClientParams)
			}
		} else if r.Method == constant.Delete && fmt.Sprintf("%s", r.Body) != "{}" {
			if err := json.NewDecoder(r.Body).Decode(&tmpReq); err != nil {
				log.Error("Execute decode request failed", "error", err.Error())
				return nil, errors2.New(errors2.ClientParams, constant.ErrClientParams)
			}
		}

		switch p.Type().Kind() {
		case reflect.Struct:
			//validate request
			if err := c.validate.Struct(tmpReq); err != nil {
				log.Error("Execute decode request failed,", "validate struct", err.Error(), "req:", req)
				return nil, errors2.New(errors2.ClientParams, Translate(err))
			}
		case reflect.Array:
			if err := c.validate.Var(tmpReq, ""); err != nil {
				log.Error("Execute decode request failed,", "validate struct", err.Error(), "req:", req)
				return nil, errors2.New(errors2.ClientParams, Translate(err))
			}
		}
		return tmpReq, nil
	}
}

// encodeResponse encode the  response(model.response -> http.response)
func (c Controller) encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	log.Debug("Execute encode response", "method", "encodeResponse")
	response := constant.Response{
		Data: resp,
	}
	operationIdKey := ctx.Value("X-App-Operation-Key")
	if operationIdKey != nil {
		operationIdKey = operationIdKey.([]string)[0]
		if operationIdKey != "" {
			// 清除operation缓存
			if err := initialize.RedisClient.Delete(operationIdKey.(string)); err != nil {
				log.Infof("del operation id key：%s,err:%s", operationIdKey, err)
			}
		}
	}
	//uri := ctx.Value(httptransport.ContextKeyRequestURI)
	//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "200"}...).Add(1)

	return httptransport.EncodeJSONResponse(ctx, w, response)
}

func (c Controller) serverOptions(before []httptransport.RequestFunc, mid []httptransport.ServerOption, after []httptransport.ServerResponseFunc) []httptransport.ServerOption {
	//copy params from Form,PostForm to Context
	copyParams := func(ctx context.Context, request *http.Request) context.Context {
		log.Debug("Merge request params to Context,", "method,", "serverBefore")
		if err := request.ParseForm(); err != nil {
			log.Error("Parse form failed", "error", err.Error())
			return ctx
		}
		improveValue := func(vs []string) interface{} {
			if len(vs) == 1 {
				return vs[0]
			}
			return vs
		}
		for k, v := range request.Form {
			ctx = context.WithValue(ctx, k, improveValue(v))
		}

		for k, v := range request.PostForm {
			ctx = context.WithValue(ctx, k, improveValue(v))
		}

		for k, v := range mux.Vars(request) {
			ctx = context.WithValue(ctx, k, v)
		}

		for k, v := range request.Header {
			ctx = context.WithValue(ctx, k, v)
		}
		return ctx
	}

	//format error
	errorEncoderOption := func(ctx context.Context, err error, w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		var response constant.Response
		method := ctx.Value(httptransport.ContextKeyRequestMethod)
		if method == "POST" {
			operationIdKey := ctx.Value("X-App-Operation-Key")
			if operationIdKey != nil {
				operationIdKey = operationIdKey.([]string)[0]
				if operationIdKey != "" {
					// 清除operation缓存
					if err := initialize.RedisClient.Delete(operationIdKey.(string)); err != nil {
						log.Infof("del operation id key：%s,err:%s", operationIdKey, err)
					}
				}
			}
		}
		//uri := ctx.Value(httptransport.ContextKeyRequestURI)
		urlPath := ctx.Value(httptransport.ContextKeyRequestPath)
		url := strings.SplitN(urlPath.(string)[1:], "/", 3)
		codeSpace := strings.ToUpper(url[1])

		respErr := errors2.Convert(err)
		errMesg, ok := errors2.StrToCode[respErr.Code()]

		log.Debugf("code: %s, string: %s, details: %s , err: %s, message:%s, ErrMesg:%s \n", respErr.Code(), respErr.String(), respErr.Details(), respErr.Err(), respErr.Message(), errMesg)

		if strings.Contains(respErr.String(), "produced zero addresses") {
			response = constant.Response{
				ErrorResp: &constant.ErrorResp{
					CodeSpace: codeSpace,
					Code:      constant.InternalFailed,
					Message:   "produced zero addresses",
				},
			}
		}

		if errMesg != "" && respErr.Message() != "" {
			switch respErr.Code() {
			case errors2.ClientParams, errors2.StatusFailed, errors2.ChainFailed, errors2.DuplicateRequest, errors2.OrderFailed, errors2.StateGatewayFailed:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "400"}...).Add(1)
				w.WriteHeader(http.StatusBadRequest) //400
			case errors2.Authentication:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "403"}...).Add(1)
				w.WriteHeader(http.StatusForbidden) //401
			case errors2.NotFound:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "404"}...).Add(1)
				w.WriteHeader(http.StatusNotFound) //404
			case errors2.NotImplemented:
				w.WriteHeader(http.StatusNotImplemented) //501
			default:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "500"}...).Add(1)
				w.WriteHeader(http.StatusInternalServerError) //500
			}
			response = constant.Response{
				ErrorResp: &constant.ErrorResp{
					CodeSpace: codeSpace,
					Code:      errMesg,
					Message:   respErr.Message(),
				},
			}
		}
		if (respErr.Code().String() == "Unknown" || respErr.Code().String() == "Unavailable" || respErr.Code().String() == "DeadlineExceeded") && respErr.Message() != "" {
			//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "500"}...).Add(1)
			w.WriteHeader(http.StatusInternalServerError) //500
			response = constant.Response{
				ErrorResp: &constant.ErrorResp{
					CodeSpace: codeSpace,
					Code:      constant.InternalFailed,
					Message:   constant.ErrInternalFailed,
				},
			}
		}

		appErr, ok := err.(constant.IError)
		if ok {
			switch appErr.Code() {
			case constant.ClientParamsError, constant.FrequentRequestsNotSupports, constant.NftStatusAbnormal,
				constant.NftClassStatusAbnormal, constant.MaximumLimitExceeded, constant.ErrOutOfGas, constant.ModuleFailed, constant.AccountFailed:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "400"}...).Add(1)
				w.WriteHeader(http.StatusBadRequest) //400
			case constant.AuthenticationFailed, constant.StructureSignTransactionFailed:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "403"}...).Add(1)
				w.WriteHeader(http.StatusForbidden) //403
			case constant.NotFound:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "404"}...).Add(1)
				w.WriteHeader(http.StatusNotFound) //404
			default:
				//metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "500"}...).Add(1)
				w.WriteHeader(http.StatusInternalServerError) //500
				appErr = constant.ErrInternal
			}
			response = constant.Response{ErrorResp: &constant.ErrorResp{
				CodeSpace: codeSpace,
				Code:      appErr.Code(),
				Message:   appErr.Error(),
			}}
		}

		bz, _ := json.Marshal(response)
		_, _ = w.Write(bz)
	}

	var options []httptransport.ServerOption
	before = append(
		[]httptransport.RequestFunc{httptransport.PopulateRequestContext, copyParams},
		before...,
	)
	options = append(options, httptransport.ServerBefore(before...))
	options = append(options, append(mid, httptransport.ServerErrorEncoder(errorEncoderOption))...)
	options = append(options, httptransport.ServerAfter(after...))
	return options
}

// Translate 错误返回
func Translate(err error) (errMsg string) {
	errs := err.(validator.ValidationErrors)
	for _, err := range errs {
		errMsg = strings.ToLower(err.Translate(trans))
	}
	return
}
