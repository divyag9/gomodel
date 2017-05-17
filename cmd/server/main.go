package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/divyag9/gomodel/pkg/interface"
	imageIdsCache "github.com/divyag9/gomodel/pkg/listByImageIds/cache"
	imageIdsDatabase "github.com/divyag9/gomodel/pkg/listByImageIds/database"
	orderNumberCache "github.com/divyag9/gomodel/pkg/listByOrderNumber/cache"
	orderNumberDatabase "github.com/divyag9/gomodel/pkg/listByOrderNumber/database"
	pb "github.com/divyag9/gomodel/pkg/pb"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	ora "gopkg.in/rana/ora.v4"
)

// Server contains information required by contentservice server
type Server struct {
	Db              *ora.Ses
	MemClient       *memcache.Client
	SecondsToExpiry int32
}

//ListByImageIds retrieves imagedetails for given imageids
func (s *Server) ListByImageIds(ctx context.Context, in *pb.ImageIdsRequest) (*pb.ListResponse, error) {
	var imageIdsGetter contentserviceinterface.ImageIdsGetter
	//Create database client
	imageIdsGetter = imageIdsDatabase.New(s.Db)

	//Retrieve cache value from context
	var cacheEnabled bool
	headers, ok := metadata.FromIncomingContext(ctx)
	if ok {
		//Get the cache header value
		cacheValue := headers["cache"][0]
		cacheBoolValue, err := strconv.ParseBool(cacheValue)
		if err != nil {
			return nil, err
		}
		cacheEnabled = cacheBoolValue
	}

	// If caching is enabled create cache client
	if cacheEnabled {
		imageIdsGetter = imageIdsCache.New(s.MemClient, s.SecondsToExpiry, imageIdsGetter)
	}

	//Retrieve imagedetails for an ordernumber
	imageDetails, err := imageIdsGetter.GetImageDetailsByImageIds(in.ImageIds)
	if err != nil {
		return nil, err
	}

	//Create response
	listResponse := &pb.ListResponse{}
	listResponse.ImageDetails = imageDetails

	return listResponse, nil
}

//ListByOrderNumber retrieves imagedetails for given ordernumber
func (s *Server) ListByOrderNumber(ctx context.Context, in *pb.OrderNumberRequest) (*pb.ListResponse, error) {
	var orderNumberGetter contentserviceinterface.OrderNumberGetter
	//Create database client
	orderNumberGetter = orderNumberDatabase.New(s.Db)

	//Retrieve cache value from context
	var cacheEnabled bool
	headers, ok := metadata.FromIncomingContext(ctx)
	if ok {
		//Get the cache header value
		cacheValue := headers["cache"][0]
		cacheBoolValue, err := strconv.ParseBool(cacheValue)
		if err != nil {
			return nil, err
		}
		cacheEnabled = cacheBoolValue
	}

	// If caching is enabled create cache client
	if cacheEnabled {
		orderNumberGetter = orderNumberCache.New(s.MemClient, s.SecondsToExpiry, orderNumberGetter)
	}

	//Retrieve imagedetails for an ordernumber
	imageDetails, err := orderNumberGetter.GetImageDetailsByOrderNumber(in.OrderNumber)
	if err != nil {
		return nil, err
	}

	//Create response
	listResponse := &pb.ListResponse{}
	listResponse.ImageDetails = imageDetails

	return listResponse, nil
}

func main() {
	//Command line parameters
	port := flag.Int("port", 10000, "The server port")
	flag.Parse()

	//Create grpc server
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	//Create database session
	dsn := os.Getenv("GO_OCI8_LIB_CONNECT_STRING")
	env, srv, session, err := ora.NewEnvSrvSes(dsn)
	if err != nil {
		fmt.Println(err)
	}
	defer env.Close()
	defer srv.Close()
	defer session.Close()

	//Get new memcache client
	memClient := getMemcacheClient()

	//Get cache expiry time from environment variable
	expiryTime, _ := strconv.Atoi(os.Getenv("EXPIRY_TIME"))

	//Create server struct
	server := &Server{Db: session,
		MemClient:       memClient,
		SecondsToExpiry: int32(expiryTime),
	}

	pb.RegisterContentServiceServer(grpcServer, server)
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Println("Failed to serve: ", err)
	}
}

//getMemcacheClient returns a new memcache client
func getMemcacheClient() *memcache.Client {
	servers := os.Getenv("MEMCACHE_SERVERS")
	memcacheServers := strings.Split(servers, ",")
	mc := memcache.New(memcacheServers...)

	return mc
}
