package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/dropboxsync/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/dropboxsync/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config  *pb.Config
	dropbox dropboxbridge
	copies  int64
}

type dropboxbridge interface {
	copyFile(key, origin, dest string) error
	listFiles(key, path string) ([]string, error)
}

func (s *Server) save(ctx context.Context) {
	s.KSclient.Save(ctx, KEY, s.config)
}

func (s *Server) load(ctx context.Context) error {
	config := &pb.Config{}
	data, _, err := s.KSclient.Read(ctx, KEY, config)

	if err != nil {
		return err
	}

	s.config = data.(*pb.Config)
	return nil
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&pb.Config{},
		&dbProd{},
		int64(0),
	}
	return s
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterDropboxSyncServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	if master {
		err := s.load(ctx)
		return err
	}

	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{
		&pbg.State{Key: "num_sync_configs", Value: int64(len(s.config.SyncConfigs))},
		&pbg.State{Key: "copies", Value: s.copies},
	}
}

func (s *Server) runAllUpdates(ctx context.Context) {
	for _, syncConfig := range s.config.SyncConfigs {
		s.runUpdate(ctx, syncConfig)
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var wipe = flag.Bool("wipe", false, "Clear configs")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.GoServer.KSclient = *keystoreclient.GetClient(server.GetIP)
	server.PrepServer()
	server.Register = server
	server.RegisterServer("dropboxsync", false)
	server.RegisterRepeatingTask(server.runAllUpdates, "run_update", time.Minute*5)

	if *wipe {
		server.config.SyncConfigs = []*pb.SyncConfig{}
		//server.save(context.ackground())
		return
	}

	fmt.Printf("%v", server.Serve())
}
