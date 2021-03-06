package main

import (
	"log"
	"os"
	"strconv"

	"github.com/brotherlogic/goserver/utils"
	"google.golang.org/grpc"

	pb "github.com/brotherlogic/dropboxsync/proto"

	//Needed to pull in gzip encoding init
	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
	host, port, err := utils.Resolve("dropboxsync", "dropboxysync-cli")
	if err != nil {
		log.Fatalf("Unable to reach server: %v", err)
	}
	conn, err := grpc.Dial(host+":"+strconv.Itoa(int(port)), grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Fatalf("Unable to dial: %v", err)
	}

	client := pb.NewDropboxSyncServiceClient(conn)
	ctx, cancel := utils.BuildContext("dropboxsync-cli", "dropboxsync")
	defer cancel()

	switch os.Args[1] {
	case "config":
		_, err := client.AddSyncConfig(ctx, &pb.AddSyncConfigRequest{ToAdd: &pb.SyncConfig{Key: os.Args[2], Origin: os.Args[3], Destination: os.Args[4]}})
		if err != nil {
			log.Fatalf("Error on GET: %v", err)
		}
	}
}
