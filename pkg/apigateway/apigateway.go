package apigateway

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"openpitrix.io/logger"
	"openpitrix.io/notification/pkg/config"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/philips/go-bindata-assetfs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gw "openpitrix.io/notification/pkg/pb"
	"io"
	"strings"
	"openpitrix.io/notification/pkg/apigateway/pkg/ui/data/swagger"

)


var (
	notificationEndpoint *string
)


func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(w, strings.NewReader(gw.Swagger))
	})

	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterNotificationHandlerFromEndpoint(ctx, gwmux, *notificationEndpoint, opts)
	if err != nil {
		return err
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

//	log.Print("Notification gRPC Server gateway start at port 8080...")
	logger.Infof(nil,"Gateway Service Started:%+v",*notificationEndpoint)


	http.ListenAndServe(":8080", mux)
	return nil
}

func Serve() {

	config.GetInstance().LoadConf()
	host:=config.GetInstance().App.Host
	port:=config.GetInstance().App.Port
	address:=host+port

	notificationEndpoint = flag.String("notification_endpoint", address, "endpoint of Notification gRPC Service")

	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}