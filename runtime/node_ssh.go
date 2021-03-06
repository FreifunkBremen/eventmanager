package runtime

import (
	"fmt"

	"github.com/bdlm/log"

	"github.com/FreifunkBremen/freifunkmanager/ssh"
)

func (n *Node) SSHUpdate(sshmgmt *ssh.Manager) bool {
	client, err := sshmgmt.ConnectTo(n.GetAddress())
	if err != nil {
		return false
	}
	defer client.Close()

	if n.Hostname != n.HostnameRespondd {
		ssh.Execute(n.Address, client, fmt.Sprintf(`
			uci set system.@system[0].hostname='%s';
			uci commit; echo $(uci get system.@system[0].hostname) > /proc/sys/kernel/hostname;`,
			n.Hostname))
	}
	if n.Owner != n.OwnerRespondd {
		ssh.Execute(n.Address, client, fmt.Sprintf(`
			uci set gluon-node-info.@owner[0].contact='%s';
			uci commit gluon-node-info;`,
			n.Owner))
	}
	if !locationEqual(n.Location, n.LocationRespondd) {
		ssh.Execute(n.Address, client, fmt.Sprintf(`
			uci set gluon-node-info.@location[0].latitude='%f';
			uci set gluon-node-info.@location[0].longitude='%f';
			uci set gluon-node-info.@location[0].share_location=1;
			uci commit gluon-node-info;`,
			n.Location.Latitude, n.Location.Longitude))

	}
	runWifi := false
	defer func() {
		if runWifi {
			ssh.Execute(n.Address, client, "wifi")
			// send warning for running wifi, because it kicks clients from node
			log.Warn("[cmd] wifi ", n.NodeID)
		}
	}()

	// update settings for 2.4GHz (802.11g) radio
	result, err := ssh.Run(n.Address, client, `
		if [ "$(uci get wireless.radio0.hwmode | grep -c g)" -ne 0 ]; then
			echo "radio0";
		elif [ "$(uci get wireless.radio1.hwmode | grep -c g)" -ne 0 ]; then
			echo "radio1";
		fi;`)
	if err != nil {
		return true
	}
	radio := ssh.SSHResultToString(result)
	ch := GetChannel(n.Wireless.Channel24)
	if radio != "" {
		if ch != nil {
			if n.Wireless.TxPower24 != n.WirelessRespondd.TxPower24 {
				ssh.Execute(n.Address, client, fmt.Sprintf(`
					uci set wireless.%s.txpower='%d';
					uci commit wireless;`,
					radio, n.Wireless.TxPower24))
				//runWifi = true
			}
			if n.Wireless.Channel24 != n.WirelessRespondd.Channel24 {
				//ubus call hostapd.%s switch_chan '{"freq":%d}'
				ssh.Execute(n.Address, client, fmt.Sprintf(`
					uci set wireless.%s.channel='%d';
					uci commit wireless;`,
					//strings.Replace(radio, "radio", "client", 1), ch.Frequency,
					radio, n.Wireless.Channel24))
				runWifi = true

			}
		}

		if n.Hostname != n.HostnameRespondd {
			ssh.Execute(n.Address, client, fmt.Sprintf(`
				uci set wireless.priv_%s.ssid='offline-%s';
				uci commit wireless;`,
				radio, n.Hostname))
			// runWifi = true <- not needed for offline-ssid
		}
	}

	// update settings for 5GHz (802.11a) radio
	result, err = ssh.Run(n.Address, client, `
		if [ "$(uci get wireless.radio0.hwmode | grep -c a)" -ne 0 ]; then
			echo "radio0";
		elif [ "$(uci get wireless.radio1.hwmode | grep -c a)" -ne 0 ]; then
			echo "radio1";
		fi;`)
	if err != nil {
		return true
	}
	radio = ssh.SSHResultToString(result)
	ch = GetChannel(n.Wireless.Channel5)
	if radio != "" {
		if ch != nil {
			if n.Wireless.TxPower5 != n.WirelessRespondd.TxPower5 {
				ssh.Execute(n.Address, client, fmt.Sprintf(`
					uci set wireless.%s.txpower='%d';
					uci commit wireless;`,
					radio, n.Wireless.TxPower5))
				//runWifi = true
			}
			if n.Wireless.Channel5 != n.WirelessRespondd.Channel5 {
				//ubus call hostapd.%s switch_chan '{"freq":%d}'
				ssh.Execute(n.Address, client, fmt.Sprintf(`
					uci set wireless.%s.channel='%d';
					uci commit wireless;`,
					//strings.Replace(radio, "radio", "client", 1), ch.Frequency,
					radio, n.Wireless.Channel5))
				runWifi = true
			}
		}
		if n.Hostname != n.HostnameRespondd {
			ssh.Execute(n.Address, client, fmt.Sprintf(`
				uci set wireless.priv_%s.ssid='offline-%s';
				uci commit wireless;`,
				radio, n.Hostname))
			// runWifi = true <- not needed for offline-ssid
		}
	}

	return true
}
