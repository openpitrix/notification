module openpitrix.io/notification

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/gin-gonic/gin v1.7.0
	github.com/golang/protobuf v1.3.3
	github.com/google/gops v0.0.0-20180903072510-f341a40f99ec
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/grpc-gateway v1.9.5
	github.com/jinzhu/gorm v1.9.11
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/sony/sonyflake v1.0.0
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.4.0
	golang.org/x/tools v0.0.0-20190312170243-e65039ee4138
	google.golang.org/genproto v0.0.0-20190404172233-64821d5d2107
	google.golang.org/grpc v1.20.1
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	openpitrix.io/libqueue v0.4.1
	openpitrix.io/logger v0.1.0
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20190212144455-93d5ec2c7f76 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.13.0 // indirect
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/go-redis/redis v6.15.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gogo/protobuf v1.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/google/uuid v1.0.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/sirupsen/logrus v1.3.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/ugorji/go/codec v1.1.7 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1 // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/net v0.0.0-20190520210107-018c4d40a106 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	golang.org/x/text v0.3.2 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df => github.com/go-mail/mail v2.3.1+incompatible

replace openpitrix.io/libqueue v0.4.1 => github.com/openpitrix/libqueue v0.4.1

replace openpitrix.io/logger v0.1.0 => github.com/openpitrix/logger v0.1.0
