package capture

import (
	"fmt"
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Track(device, dst string, promiscuos bool) {

	var snaphotLen int32

	snaphotLen = 1024
	timeout := 30 * time.Second
	handle, err := pcap.OpenLive(device, snaphotLen, promiscuos, timeout)

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {

		ipLayer := packet.Layer(layers.LayerTypeIPv4)

		if ipLayer == nil {
			continue
		}

		ip, _ := ipLayer.(*layers.IPv4)
		dstIP := fmt.Sprintf("%s", ip.DstIP)

		if dstIP == dst {
			fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
			fmt.Println("Protocol : ", ip.Protocol)
		}
	}
}
