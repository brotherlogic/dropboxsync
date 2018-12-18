package main

import "golang.org/x/net/context"
import pb "github.com/brotherlogic/dropboxsync/proto"

//UpdateConfig updates the config
func (s *Server) UpdateConfig(ctx context.Context, req *pb.UpdateConfigRequest) (*pb.UpdateConfigResponse, error) {
	s.config.CoreKey = req.NewCoreKey
	s.save(ctx)

	return &pb.UpdateConfigResponse{}, nil
}

//AddSyncConfig adds a basic sync config
func (s *Server) AddSyncConfig(ctx context.Context, req *pb.AddSyncConfigRequest) (*pb.AddSyncConfigResponse, error) {
	s.config.SyncConfigs = append(s.config.SyncConfigs, req.ToAdd)
	s.save(ctx)

	return &pb.AddSyncConfigResponse{}, nil
}
