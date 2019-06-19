module openpitrix.io/notification

require (
	github.com/bitly/go-simplejson v0.5.0
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/coreos/go-systemd v0.0.0-20190212144455-93d5ec2c7f76 // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20190204142019-df6d76eb9289 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/fatih/camelcase v1.0.0
	github.com/fatih/structs v1.1.0
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/gin-contrib/sse v0.0.0-20190125020943-a7658810eb74 // indirect
	github.com/gin-gonic/gin v1.3.0
	github.com/go-redis/redis v6.15.2+incompatible
	github.com/go-sql-driver/mysql v1.4.1 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef // indirect
	github.com/golang/protobuf v1.3.1
	github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c // indirect
	github.com/google/go-cmp v0.2.0 // indirect
	github.com/google/gops v0.0.0-20180903072510-f341a40f99ec
	github.com/gorilla/websocket v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.7.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v0.0.0-20181116074157-8ec929ed50c3 // indirect
	github.com/jonboulle/clockwork v0.1.0 // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/koding/multiconfig v0.0.0-20171124222453-69c27309b2d7
	github.com/kr/pretty v0.1.0 // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v0.9.2 // indirect
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90 // indirect
	github.com/sirupsen/logrus v1.3.0 // indirect
	github.com/soheilhy/cmux v0.1.4 // indirect
	github.com/sony/sonyflake v0.0.0-20181109022403-6d5bd6181009
	github.com/speps/go-hashids v2.0.0+incompatible
	github.com/stretchr/testify v1.2.2
	github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5 // indirect
	github.com/ugorji/go/codec v0.0.0-20190128213124-ee1426cffec0 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	go.etcd.io/bbolt v1.3.2 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1 // indirect
	golang.org/x/net v0.0.0-20190520210107-018c4d40a106
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	golang.org/x/tools v0.0.0-20190311212946-11955173bddd
	google.golang.org/genproto v0.0.0-20190128161407-8ac453e89fca
	google.golang.org/grpc v1.20.1
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/yaml.v2 v2.2.2 // indirect
	openpitrix.io/libqueue v0.3.1
	openpitrix.io/logger v0.1.0
)

replace openpitrix.io/libqueue v0.3.1 => github.com/openpitrix/libqueue v0.3.1

replace gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df => github.com/go-mail/mail v2.3.1+incompatible
