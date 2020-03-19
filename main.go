package main

import "packet_capture/capture"

func main() {
	capture.Track("enp2s0", "192.168.1.7", "192.168.1.7", false)

}
