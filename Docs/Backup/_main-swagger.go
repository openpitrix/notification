package Backup

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/philips/go-bindata-assetfs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gw "notificationService/pkg/nf/pb"
	"io"
	"strings"
	"notificationService/apigateway/pkg/ui/data/swagger"
)
var (
	greeterEndpoint = flag.String("helloworld_endpoint", "localhost:50051", "endpoint of Greeter gRPC Service")
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
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, *greeterEndpoint, opts)
	if err != nil {
		return err
	}

	mux.Handle("/", gwmux)
	serveSwagger(mux)

	log.Print("Greeter gRPC Server gateway start at port 8080...")
	http.ListenAndServe(":8080", mux)
	return nil
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}