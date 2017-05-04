package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
<<<<<<< HEAD
	pb "github.com/divyag9/gomodel/pkg/pb/github.com/divyag9/proto"
=======
	"github.com/divyag9/gomodel/pkg/cache"
	"github.com/divyag9/gomodel/pkg/database"
	"github.com/divyag9/gomodel/pkg/interface"
	pb "github.com/divyag9/gomodel/pkg/pb"
>>>>>>> Adding more composition
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	ora "gopkg.in/rana/ora.v4"
)

// Server contains information required by server
type Server struct {
	Db              *ora.Ses
	MemClient       *memcache.Client
	SecondsToExpiry int32
}

//ListByImageIds retrieves imagedetails for given imageids
func (s *Server) ListByImageIds(ctx context.Context, in *pb.ImageIdsRequest) (*pb.ListResponse, error) {
	//TO be implemented
	return nil, nil
}

//ListByOrderNumber retrieves imagedetails for given ordernumber
func (s *Server) ListByOrderNumber(ctx context.Context, in *pb.OrderNumberRequest) (*pb.ListResponse, error) {
	var orderNumberGetter contentserviceinterface.OrderNumberGetter
	//Create database client
	orderNumberGetter = database.New(s.Db)

	//Retrieve cache_enabled value from the context
	cacheEnabled := ctx.Value("cache_enabled").(bool)
	// If caching is enabled create cache client
	if cacheEnabled {
		orderNumberGetter = cache.New(s.MemClient, s.SecondsToExpiry, orderNumberGetter)
	}

	//Retrieve imagedetails for an ordernumber
	imageDetails, err := orderNumberGetter.GetImageDetailsByOrderNumber(in.OrderNumber)
	if err != nil {
		return nil, err
	}

	//Create the response
	listResponse := &pb.ListResponse{}
	listResponse.ImageDetails = imageDetails

	return listResponse, nil
}

func main() {
	tls := flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile := flag.String("cert_file", "testdata/server1.pem", "The TLS cert file")
	keyFile := flag.String("key_file", "testdata/server1.key", "The TLS key file")
	port := flag.Int("port", 10000, "The server port")

	flag.Parse()
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			grpclog.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)

	dsn := os.Getenv("GO_OCI8_CONNECT_STRING")
	env, srv, ses, err := ora.NewEnvSrvSes(dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer env.Close()
	defer srv.Close()
	defer ses.Close()

	mc := getMemcacheClient()
	expiryTime, _ := strconv.Atoi(os.Getenv("EXPIRY_TIME"))

	server := &Server{}
	server.Db = ses
	server.MemClient = mc
	server.SecondsToExpiry = int32(expiryTime)
	pb.RegisterContentServiceServer(grpcServer, server)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}

func getMemcacheClient() *memcache.Client {
	servers := os.Getenv("MEMCACHE_SERVERS")
	memcacheServers := strings.Split(servers, ",")
	mc := memcache.New(memcacheServers...)

	return mc
}
