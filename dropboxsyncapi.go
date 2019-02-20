package main

import (
	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dropboxsync/proto"
)

//AddSyncConfig adds a basic sync config
func (s *Server) AddSyncConfig(ctx context.Context, req *pb.AddSyncConfigRequest) (*pb.AddSyncConfigResponse, error) {
	s.config.SyncConfigs = append(s.config.SyncConfigs, req.ToAdd)
	s.save(ctx)

	return &pb.AddSyncConfigResponse{}, nil
}
