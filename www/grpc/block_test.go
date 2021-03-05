package grpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	zarb "github.com/zarbchain/zarb-go/www/grpc/proto"
)

func TestGetBlock(t *testing.T) {
	conn, client := callServer(t)

	t.Run("Should return nil for non existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1, Verbosity: 0})
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	b1, trxs := block.GenerateTestBlock(nil, nil)
	tMockState.AddBlock(1, b1, trxs)

	t.Run("Should return an existing block ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1, Verbosity: 0})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		h, err := crypto.HashFromString(res.Hash)
		assert.NoError(t, err)
		assert.Equal(t, h, b1.Hash())
		assert.Empty(t, res.Json)

	})

	t.Run("Should return json object with verbosity 1 ", func(t *testing.T) {
		res, err := client.GetBlock(tCtx, &zarb.BlockRequest{Height: 1, Verbosity: 1})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		h, err := crypto.HashFromString(res.Hash)
		assert.NoError(t, err)
		assert.Equal(t, h, b1.Hash())
		assert.NotEmpty(t, res.Json)
	})

	conn.Close()
}
