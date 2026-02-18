package main

import (
    "context"
    "flag"
    "fmt"
    "net/http"
    "strings"
    "time"

    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    gw "github.com/lindaprotocol/grpc-gateway/pkg/api/protocol"
)

var (
    grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
    httpPort           = flag.String("http-port", ":18890", "HTTP port")
)

func corsMiddleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        h.ServeHTTP(w, r)
    })
}

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

    // Register Wallet service
    if err := gw.RegisterWalletHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
        return err
    }

    // Register WalletSolidity service
    if err := gw.RegisterWalletSolidityHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
        return err
    }

    // Register Scan service (from your api.proto)
    if err := gw.RegisterScanServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
        return err
    }

    glog.Infof("HTTP server listening on %s", *httpPort)
    return http.ListenAndServe(*httpPort, corsMiddleware(mux))
}

func main() {
    flag.Parse()
    defer glog.Flush()

    if err := run(); err != nil {
        glog.Fatal(err)
    }
}
