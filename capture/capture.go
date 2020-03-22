package capture

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type Tracker struct {
	message string
}

func NewTracker() *Tracker {
	tracker := Tracker{}
	tracker.message = ""
	return &tracker
}

func (t *Tracker) Track(device, src, dst string, promiscuos bool, savePath string) {

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)

	go t.track(device, src, dst, promiscuos)
	<-killSignal

	if t.message != "" {
		fmt.Println("Saving the message")

	}

	if t.save(savePath) == false {
		fmt.Println("Failed to save data")
	}

	fmt.Println("data saved")
}

func (t *Tracker) save(dstPath string) bool {

	if t.message == "" {
		fmt.Println("No packet message to save")
		return false
	}

	f, err := os.Create(dstPath)

	if err != nil {
		fmt.Println("Error creating a file to save the message")
		fmt.Println(err)
		return false
	}

	n, err := f.Write([]byte(t.message))

	if err != nil {
		fmt.Println("Error writting the message to the dst")
		fmt.Println(err)
		return false
	}

	fmt.Printf("wrote %d bytes.\n", n)
	return true

}

func (t *Tracker) track(device, src, dst string, promiscuos bool) {

	var snaphotLen int32

	snaphotLen = 1024
	timeout := 30 * time.Second
	handle, err := pcap.OpenLive(device, snaphotLen, promiscuos, timeout)

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	for packet := range packetSource.Packets() {

		ipLayer := packet.Layer(layers.LayerTypeIPv4)

		if ipLayer == nil {
			continue
		}

		ip, _ := ipLayer.(*layers.IPv4)

		str := ""

		if src != "" {
			srcIP := fmt.Sprintf("%s", ip.SrcIP)

			if srcIP == src {
				str = fmt.Sprintf("From %s to %s\n", ip.SrcIP, ip.DstIP)
				fmt.Println(str)
			}
		}

		if dst != "" {

			dstIP := fmt.Sprintf("%s", ip.DstIP)

			if dstIP == dst {
				str = fmt.Sprintf("From %s to %s\n", ip.SrcIP, ip.DstIP)
				fmt.Println(str)
			}
		}

		if str == "" {
			continue
		}

		t.message += str
	}

}
