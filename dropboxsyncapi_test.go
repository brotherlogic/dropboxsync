package main

import (
	"context"
	"testing"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/dropboxsync/proto"
)

func InitServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

func TestAddSyncKey(t *testing.T) {
	s := InitServer()
	_, err := s.AddSyncConfig(context.Background(), &pb.AddSyncConfigRequest{ToAdd: &pb.SyncConfig{Key: "newkey"}})
	if err != nil {
		t.Fatalf("Error on update: %v", err)
	}
	if len(s.config.SyncConfigs) != 1 && s.config.SyncConfigs[0].Key != "newkey" {
		t.Errorf("Update did not take: %v", s.config)
	}
}
