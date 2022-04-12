package kit

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/redis"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	en2 "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/log"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/metric"
	"gitlab.bianjie.ai/irita-paas/open-api/internal/pkg/types"
)

// TimeLayout time format
const (
	TimeLayout = "2006-01-02"
	Delete     = "DELETE"
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
		Initialize()
		Stop()
	}

	Controller struct {
		validate *validator.Validate
	}
)

var trans ut.Translator

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

	tim, err := time.Parse(TimeLayout, value.(string))
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

// decodeRequest decode request(http.request -> model.request)
func (c Controller) decodeRequest(req interface{}) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (request interface{}, err error) {
		log.Debug("Execute decode request", "method", "decodeRequest")
		if req == nil {
			return nil, err
		}
		p := reflect.ValueOf(req).Elem()
		p.Set(reflect.Zero(p.Type()))
		if r.Method != Delete {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Error("Execute decode request failed", "error", err.Error())
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrClientParams)
			}
		} else if r.Method == Delete && fmt.Sprintf("%s", r.Body) != "{}" {
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				log.Error("Execute decode request failed", "error", err.Error())
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, types.ErrClientParams)
			}
		}

		switch p.Type().Kind() {
		case reflect.Struct:
			//validate request
			if err := c.validate.Struct(req); err != nil {
				log.Error("Execute decode request failed", "validate struct", err.Error(), "req:", req)
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, Translate(err))
			}
		case reflect.Array:
			if err := c.validate.Var(req, ""); err != nil {
				log.Error("Execute decode request failed", "validate struct", err.Error(), "req:", req)
				return nil, types.NewAppError(types.RootCodeSpace, types.ClientParamsError, Translate(err))
			}
		}
		return req, nil
	}
}

// encodeResponse encode the  response(model.response -> http.response)
func (c Controller) encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	log.Debug("Execute encode response", "method", "encodeResponse")
	response := Response{
		Data: resp,
	}
	method := ctx.Value(httptransport.ContextKeyRequestMethod)
	uri := ctx.Value(httptransport.ContextKeyRequestURI)
	appID := ctx.Value("X-App-Id")
	operationID := w.Header().Get("X-Operation-ID")
	key := fmt.Sprintf("%s:%s", appID, operationID)
	if operationID != "" {
		if err := redis.Set(key, "1", time.Second*60); err != nil {
			log.Error("encode redis set err: ", err)
		}
	}
	metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "200"}...).Add(1)
	return httptransport.EncodeJSONResponse(ctx, w, response)
}

func (c Controller) serverOptions(
	before []httptransport.RequestFunc,
	mid []httptransport.ServerOption,
	after []httptransport.ServerResponseFunc,
) []httptransport.ServerOption {
	//copy params from Form,PostForm to Context
	copyParams := func(ctx context.Context, request *http.Request) context.Context {
		log.Debug("Merge request params to Context", "method", "serverBefore")
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
		var response Response
		method := ctx.Value(httptransport.ContextKeyRequestMethod)
		uri := ctx.Value(httptransport.ContextKeyRequestURI)
		urlPath := ctx.Value(httptransport.ContextKeyRequestPath)
		url := strings.SplitN(urlPath.(string)[1:], "/", 3)
		codeSpace := strings.ToUpper(url[1])
		appErr, ok := err.(types.IError)
		if !ok {
			metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "500"}...).Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			response = Response{
				ErrorResp: &ErrorResp{
					CodeSpace: codeSpace,
					Code:      types.ErrInternal.Code(),
					Message:   types.ErrInternal.Error(),
				},
			}
		} else {
			switch appErr.Code() {
			case types.ClientParamsError, types.FrequentRequestsNotSupports, types.NftStatusAbnormal,
				types.NftClassStatusAbnormal, types.MaximumLimitExceeded, types.ErrOutOfGas, types.ModuleFailed, types.AccountFAiled:
				metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "400"}...).Add(1)
				w.WriteHeader(http.StatusBadRequest) //400
			case types.AuthenticationFailed, types.StructureSignTransactionFailed:
				metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "403"}...).Add(1)
				w.WriteHeader(http.StatusForbidden) //403
			case types.NotFound:
				metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "404"}...).Add(1)
				w.WriteHeader(http.StatusNotFound) //404
			default:
				metric.NewPrometheus().ApiHttpRequestCount.With([]string{"method", method.(string), "uri", uri.(string), "code", "500"}...).Add(1)
				w.WriteHeader(http.StatusInternalServerError) //500
				appErr = types.ErrInternal
			}
			response = Response{ErrorResp: &ErrorResp{
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
