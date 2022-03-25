package chain

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"

	sdk "github.com/irisnet/core-sdk-go"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
)

var sdkClients map[string]*client

type (
	client struct {
		Client sdk.Client
		Gas    uint64
		Denom  string
		Amount int64
	}
	sdkClient struct {
		RcpAddr          string `json:"rcpAddr"`
		GrpcAddr         string `json:"grpcAddr"`
		WsAddr           string `json:"wsAddr"`
		ChainID          string `json:"chainID"`
		ProjectID        string `json:"projectID"`
		ProjectKey       string `json:"projectKey"`
		ChainAccountAddr string `json:"chainAccountAddr"`

		GasCoefficient float64 `json:"gasCoefficient"`
		Gas            uint64  `json:"gas"`
		Denom          string  `json:"denom"`
		Amount         int64   `json:"amount"`
		AccoutGas      int64   `json:"accoutGas"`

		ChainEncryption string `json:"chainEncryption"`
	}
)

func GetSdkClients() map[string]*client {
	return sdkClients
}

func NewSdkClient(confs map[string]interface{}, db *sql.DB) {
	sdkClients = make(map[string]*client, len(confs))
	for k, v := range confs {
		var sdkclient sdkClient
		valueByte, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(valueByte, &sdkclient)
		if err != nil {
			panic(err)
		}
		authToken := NewAuthToken(sdkclient.ProjectID, sdkclient.ProjectKey, sdkclient.ChainAccountAddr)
		// overwrite grpcOpts
		grpcOpts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithPerRPCCredentials(&authToken),
		}

		httpHeader := http.Header{}
		if projectKey := authToken.GetProjectKey(); projectKey != "" {
			httpHeader.Set("x-api-key", authToken.GetProjectKey())
		}

		options := []sdktype.Option{
			sdktype.KeyDAOOption(NewMsqlKeyDao(db, sdkclient.ChainEncryption)),
			sdktype.AlgoOption(sdkclient.ChainEncryption),
			sdktype.TimeoutOption(60),
			sdktype.CachedOption(false),
			sdktype.GRPCOptions(grpcOpts),
			sdktype.HeaderOption(httpHeader),
			sdktype.WSAddrOption(sdkclient.WsAddr),
		}

		cfg, err := sdktype.NewClientConfig(sdkclient.RcpAddr, sdkclient.GrpcAddr, sdkclient.ChainID, options...)
		if err != nil {
			panic(err)
		}
		sct := sdk.NewClient(cfg)
		nftClient := nft.NewClient(sct.BaseClient, sct.AppCodec())
		sct.RegisterModule(nftClient)
		value := &client{
			sct,
			sdkclient.Gas,
			sdkclient.Denom,
			sdkclient.Amount,
		}
		sdkClients[k] = value
	}
}
