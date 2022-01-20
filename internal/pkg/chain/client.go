package chain

import (
	"database/sql"
	"net/http"

	sdk "github.com/irisnet/core-sdk-go"
	sdktype "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/irismod-sdk-go/nft"
	configs "gitlab.bianjie.ai/irita-paas/open-api/config"
	"google.golang.org/grpc"
)

var sdkClient sdk.Client

func GetSdkClient() sdk.Client {
	return sdkClient
}

func NewSdkClient(conf configs.Chain, db *sql.DB) {
	authToken := NewAuthToken(conf.ProjectID, conf.ProjectKey, conf.ChainAccountAddr)
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
		sdktype.KeyDAOOption(NewMsqlKeyDao(db)),
		sdktype.TimeoutOption(60),
		sdktype.CachedOption(true),
		sdktype.GRPCOptions(grpcOpts),
		sdktype.HeaderOption(httpHeader),
		sdktype.WSAddrOption(conf.WsAddr),
	}

	cfg, err := sdktype.NewClientConfig(conf.RcpAddr, conf.GrpcAddr, conf.ChainID, options...)
	if err != nil {
		panic(err)
	}

	sdkClient = sdk.NewClient(cfg)
	nftClient := nft.NewClient(sdkClient.BaseClient, sdkClient.AppCodec())
	sdkClient.RegisterModule(nftClient)
}
