package websocket

import (
	log "github.com/sirupsen/logrus"

	wsLib "dev.sum7.eu/genofire/golang-lib/websocket"

	"github.com/FreifunkBremen/freifunkmanager/runtime"
)

var wifi24Channels []uint32
var wifi5Channels []uint32

func (ws *WebsocketServer) connectHandler(logger *log.Entry, msg *wsLib.Message) error {
	//msg.From.Write(&wsLib.Message{Subject: MessageTypeStats, Body: ws.nodes.Statistics})
	var nodes []*runtime.Node
	var count int

	ws.db.Find(&nodes).Count(&count)

	ws.nodes.Lock()
	i := 0
	for _, node := range ws.nodes.List {
		n := runtime.NewNode(node, "")
		if n == nil {
			continue
		}
		n.Lastseen = node.Lastseen
		msg.From.Write(&wsLib.Message{Subject: MessageTypeCurrentNode, Body: n})
		i++
	}
	ws.nodes.Unlock()
	for _, node := range nodes {
		msg.From.Write(&wsLib.Message{Subject: MessageTypeSystemNode, Body: node})
	}
	msg.From.Write(&wsLib.Message{Subject: MessageTypeChannelsWifi24, Body: wifi24Channels})
	msg.From.Write(&wsLib.Message{Subject: MessageTypeChannelsWifi5, Body: wifi5Channels})
	logger.Debugf("done - fetch %d nodes and send %d", count, i)
	return nil
}

func init() {
	for ch, channel := range runtime.ChannelList {
		if runtime.ChannelEU && !channel.AllowedInEU {
			continue
		}
		if channel.Frequenz > runtime.FREQ_THREASHOLD {
			wifi5Channels = append(wifi5Channels, ch)
		} else {
			wifi24Channels = append(wifi24Channels, ch)
		}
	}
}
