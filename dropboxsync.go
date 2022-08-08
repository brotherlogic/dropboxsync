package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/dropboxsync/proto"
	pbg "github.com/brotherlogic/goserver/proto"
)

var (
	lastRun = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "dropboxsync_run",
		Help: "ImageAge",
	}, []string{"err"})
)

const (
	// KEY - where the wants are stored
	KEY = "/github.com/brotherlogic/dropboxsync/config"
)

//Server main server type
type Server struct {
	*goserver.GoServer
	config   *pb.Config
	dropbox  dropboxbridge
	copies   int64
	listTime time.Duration
	copyTime time.Duration
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

	s.Log(fmt.Sprintf("Loaded config %v", s.config))
	return nil
}

// Init builds the server
func Init() *Server {
	s := &Server{
		&goserver.GoServer{},
		&pb.Config{},
		&dbProd{},
		int64(0),
		time.Millisecond,
		time.Millisecond,
	}
	s.dropbox = &dbProd{log: s.Log}
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

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.save(ctx)
	return nil
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
		&pbg.State{Key: "thing", Value: int64(1)},
	}
}

func (s *Server) runAllUpdates(ctx context.Context) (time.Time, error) {
	err := s.load(ctx)
	if err == nil {
		s.Log(fmt.Sprintf("Running update: %v", len(s.config.SyncConfigs)))
		for _, syncConfig := range s.config.SyncConfigs {
			time.Sleep(time.Second * 5)
			s.Log(fmt.Sprintf("Running update for %v", syncConfig))
			s.runUpdate(ctx, syncConfig)
		}
	}

	return time.Now().Add(time.Minute * 10), nil
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var wipe = flag.Bool("wipe", false, "Clear configs")
	var token = flag.String("token", "", "Initial token")
	var origin = flag.String("origin", "", "Origin")
	var dest = flag.String("dest", "", "Destination")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer("dropboxsync")
	server.Register = server
	err := server.RegisterServerV2(false)
	if err != nil {
		return
	}

	if *wipe {
		ctx, cancel := utils.BuildContext("dropboxysync", "dropboxsync")
		defer cancel()

		server.config.SyncConfigs = []*pb.SyncConfig{}
		server.save(ctx)
		return
	}

	if len(*token) > 0 {
		ctx, cancel := utils.BuildContext("dropboxysync", "dropboxsync")
		defer cancel()

		server.config.SyncConfigs = append(server.config.SyncConfigs, &pb.SyncConfig{Key: *token, Origin: *origin, Destination: *dest})
		server.save(ctx)
		return
	}

	go func() {
		for true {
			ctx, cancel := utils.ManualContext("dropboxsync", time.Minute)
			_, err = server.runAllUpdates(ctx)
			lastRun.With(prometheus.Labels{"err": fmt.Sprintf("%v", err)}).Set(float64(time.Now().Unix()))
			if err != nil {
				log.Fatalf("Cannot run the update: %v", ctx)
			}
			cancel()
			time.Sleep(time.Hour * 2)
		}

	}()

	fmt.Printf("%v", server.Serve())
}
