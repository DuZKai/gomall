module gomall/app/frontend

go 1.24.4

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	github.com/cloudwego/hertz v0.10.0
	github.com/cloudwego/kitex v0.14.1
	github.com/hertz-contrib/cors v0.1.0
	github.com/hertz-contrib/gzip v0.0.3
	github.com/hertz-contrib/logger/accesslog v0.0.0-20241107070745-e4ce8c54dd97
	github.com/hertz-contrib/logger/logrus v1.0.1
	github.com/hertz-contrib/pprof v0.1.2
	github.com/hertz-contrib/sessions v1.0.3
	github.com/joho/godotenv v1.5.1
	github.com/kitex-contrib/registry-consul v0.2.0
	github.com/kr/pretty v0.3.1
	github.com/redis/go-redis/v9 v9.11.0
	go.uber.org/zap v1.27.0
	google.golang.org/protobuf v1.36.6
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/validator.v2 v2.0.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/bytedance/gopkg v0.1.2 // indirect
	github.com/bytedance/sonic v1.13.3 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cloudwego/configmanager v0.2.3 // indirect
	github.com/cloudwego/dynamicgo v0.6.4 // indirect
	github.com/cloudwego/fastpb v0.0.5 // indirect
	github.com/cloudwego/frugal v0.2.5 // indirect
	github.com/cloudwego/gopkg v0.1.5 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/cloudwego/localsession v0.1.2 // indirect
	github.com/cloudwego/netpoll v0.7.1 // indirect
	github.com/cloudwego/runtimex v0.1.1 // indirect
	github.com/cloudwego/thriftgo v0.4.2 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/google/pprof v0.0.0-20240727154555-813a5fbdbec8 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.2.2 // indirect
	github.com/hashicorp/consul/api v1.26.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/jhump/protoreflect v1.8.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nyaruka/phonenumbers v1.3.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tidwall/gjson v1.17.3 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20240112132812-db7319d0e0e3 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.20.0 // indirect
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
