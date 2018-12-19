package main

import (
	"time"

	"golang.org/x/net/context"

	pb "github.com/brotherlogic/dropboxsync/proto"
	pbt "github.com/brotherlogic/tracer/proto"
)

//UpdateConfig updates the config
func (s *Server) UpdateConfig(ctx context.Context, req *pb.UpdateConfigRequest) (*pb.UpdateConfigResponse, error) {
	ctx = s.LogTrace(ctx, "UpdateConfig", time.Now(), pbt.Milestone_START_FUNCTION)
	s.config.CoreKey = req.NewCoreKey
	s.save(ctx)

	s.LogTrace(ctx, "UpdateConfig", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.UpdateConfigResponse{}, nil
}

//AddSyncConfig adds a basic sync config
func (s *Server) AddSyncConfig(ctx context.Context, req *pb.AddSyncConfigRequest) (*pb.AddSyncConfigResponse, error) {
	ctx = s.LogTrace(ctx, "AddSyncConfig", time.Now(), pbt.Milestone_START_FUNCTION)
	s.config.SyncConfigs = append(s.config.SyncConfigs, req.ToAdd)
	s.save(ctx)

	s.LogTrace(ctx, "AddSyncConfig", time.Now(), pbt.Milestone_END_FUNCTION)
	return &pb.AddSyncConfigResponse{}, nil
}
