module gitlab.bianjie.ai/irita-paas/open-api

go 1.17

require (
	github.com/bsm/redislock v0.7.2
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/validator/v10 v10.10.0
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gorilla/mux v1.8.0
	github.com/irisnet/core-sdk-go v0.0.0-20220106085924-448b745f3429
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	github.com/volatiletech/sqlboiler/v4 v4.8.3
	gitlab.bianjie.ai/irita-paas/orms/orm-nft v0.0.0-00010101000000-000000000000
)

require (
	github.com/ChainSafe/go-schnorrkel v0.0.0-20200405005733-88cbf1b4c40d // indirect
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cosmos/go-bip39 v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gtank/merlin v0.1.1 // indirect
	github.com/gtank/ristretto255 v0.1.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/copier v0.3.4 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.2.1-0.20191011153232-f91d3411e481 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/mimoo/StrobeGo v0.0.0-20181016162300-f8f6d4d2b643 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/petermattis/goid v0.0.0-20180202154549-b0b1615b78e5 // indirect
	github.com/sasha-s/go-deadlock v0.2.0 // indirect
	github.com/spf13/afero v1.8.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15 // indirect
	github.com/tendermint/go-amino v0.16.0 // indirect
	github.com/tendermint/tendermint v0.34.11 // indirect
	github.com/tjfoc/gmsm v1.4.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/null/v8 v8.1.2 // indirect
	github.com/volatiletech/randomize v0.0.1 // indirect
	github.com/volatiletech/sqlboiler v3.7.1+incompatible // indirect
	github.com/volatiletech/strmangle v0.0.1 // indirect
	golang.org/x/crypto v0.0.0-20211108221036-ceb1ce70b4fa // indirect
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/grpc v1.42.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

replace gitlab.bianjie.ai/irita-paas/orms/orm-nft => gitlab.bianjie.ai/irita-paas/orms/orm-nft.git v0.0.0-20220119031143-9b260fb80637

replace github.com/tendermint/tendermint => github.com/bianjieai/tendermint v0.34.1-irita-210113
