module gitlab.bianjie.ai/avata/open-api

go 1.17

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/go-kit/kit v0.12.0
	github.com/go-playground/locales v0.14.0
	github.com/go-playground/universal-translator v0.18.0
	github.com/go-playground/validator/v10 v10.10.1
	github.com/go-redis/redis/v8 v8.11.4
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.10.1
	github.com/volatiletech/sqlboiler v3.7.1+incompatible
	google.golang.org/grpc v1.45.0
	gorm.io/gorm v1.23.4
)

require (
	github.com/deckarep/golang-set v1.8.0
	gitlab.bianjie.ai/avata/chains/api v0.7.2-0.20220830014600-0360685e7537
	gitlab.bianjie.ai/avata/utils/commons v0.0.0-00010101000000-000000000000
	gitlab.bianjie.ai/avata/utils/errors v0.0.0-00010101000000-000000000000
	gorm.io/driver/mysql v1.3.4
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/ericlagergren/decimal v0.0.0-20181231230500-73749d4874d5 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lib/pq v1.2.1-0.20191011153232-f91d3411e481 // indirect
	github.com/magiconair/properties v1.8.5 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/pelletier/go-toml v1.9.4 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.30.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/spf13/afero v1.6.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/null v8.0.0+incompatible // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20211210111614-af8b64212486 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	gitlab.bianjie.ai/avata/chains/api => gitlab.bianjie.ai/avata/chains/api.git v0.11.1-0.20230316073345-53024131ea27

	gitlab.bianjie.ai/avata/utils/commons => gitlab.bianjie.ai/avata/utils/commons.git v0.2.0

	gitlab.bianjie.ai/avata/utils/errors => gitlab.bianjie.ai/avata/utils/errors.git v0.1.10-0.20230317063241-b437928125ba
)
