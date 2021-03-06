package network_api

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/sync/message"
)

type NetworkAPI interface {
	Start() error
	Stop()
	PublishMessage(msg *message.Message) error
	JoinDownloadTopic() error
	LeaveDownloadTopic()
	SelfID() peer.ID
}
