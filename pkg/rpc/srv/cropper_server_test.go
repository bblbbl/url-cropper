package srv

import (
	"context"
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"path/filepath"
	"regexp"
	"testing"
	"urls/internal/service"
	"urls/pkg/etc"
	cropper "urls/pkg/rpc/proto"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := cropper.NewUrlCropperClient(conn)
	resp, err := client.CropUrl(ctx, &cropper.CropRequest{Url: "test_url"})
	if err != nil {
		t.Fatalf("rpc crop url failed: %v", err)
	}

	cnf := etc.GetConfig()
	reg := regexp.MustCompile(fmt.Sprintf("%s:\\/\\/%s\\/go\\/(.)+", cnf.Http.Schema, cnf.App.Host))

	assert.MatchRegex(t, resp.Url, reg)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	path, err := filepath.Abs("../../../.env.test")
	if err != nil {
		log.Fatal("failed to get root path")
	}

	err = godotenv.Load(path)
	if err != nil {
		log.Fatal("failed to load .env")
	}

	etc.InitLogger()
	etc.InitConfig()

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	cropper.RegisterUrlCropperServer(s, NewCropperServer(service.NewWriteExecutor().Start()))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}
