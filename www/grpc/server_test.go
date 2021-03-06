package grpc

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/state"
	"github.com/zarbchain/zarb-go/sync"
	"github.com/zarbchain/zarb-go/txpool"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var tMockState *state.MockState
var tMockPool txpool.TxPool
var tMockSync *sync.MockSync
var tListener *bufconn.Listener
var tCtx context.Context

func init() {
	logger.InitLogger(logger.TestConfig())

	const bufSize = 1024 * 1024

	committee, _ := committee.GenerateTestCommittee()
	tListener = bufconn.Listen(bufSize)
	tMockState = state.MockingState(committee)
	tMockPool = txpool.MockingTxPool()
	tMockSync = sync.MockingSync()
	tCtx = context.Background()

	s := grpc.NewServer()
	server := &zarbServer{
		state:  tMockState,
		store:  tMockState.StoreReader(),
		txPool: tMockPool,
		sync:   tMockSync,
		logger: logger.NewLogger("_grpc", nil),
	}
	zarb.RegisterZarbServer(s, server)
	go func() {
		if err := s.Serve(tListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return tListener.Dial()
}

func callServer(t *testing.T) (*grpc.ClientConn, zarb.ZarbClient) {
	conn, err := grpc.DialContext(tCtx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return conn, zarb.NewZarbClient(conn)
}
