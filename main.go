package main

import "packet_capture/capture"

func main() {

	tracker := capture.NewTracker()

	tracker.Track("enp2s0", "192.168.1.2", "192.168.1.2", false, "./record.txt")

}
