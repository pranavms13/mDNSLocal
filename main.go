package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/TwiN/go-color"
	"github.com/pion/mdns/v2"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func GetFileExtension() string {
	switch operatingSystem := runtime.GOOS; operatingSystem {
	case "windows":
		return ".exe"
	case "linux":
		return ""
	case "darwin":
		return ""
	default:
		fmt.Println("OS Unsupported")
		os.Exit(0)
		return ""
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func printUsage() {
	ext := GetFileExtension()
	println(color.Ize(color.Green, "mDNSLocal - A simple mDNS server"))
	println("Usage:")
	print("\tmDNSLocal", ext, color.Colorize(color.Blue, " <hostname> <ip>\n"))
	print("\tmDNSLocal", ext, color.Colorize(color.Blue, " -h, --help, help\n"))
	println("Examples:")
	print("\tmDNSLocal", ext)
	print("\tmDNSLocal", ext, "", color.Colorize(color.Blue, " my-laptop.local\n"))
	print("\tmDNSLocal", ext, "", color.Colorize(color.Blue, " my-laptop.local 192.168.0.1\n"))
}

func getHostname() (string, net.IP) {
	args := os.Args[1:]
	hostName := ""
	if len(args) > 0 {
		if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
			printUsage()
			os.Exit(0)
		}
		hostName = args[0]
		if !strings.HasSuffix(hostName, ".local") {
			fmt.Println("The hostname must end with .local")
			printUsage()
			os.Exit(1)
		}
	}
	ip := net.ParseIP("")
	if len(args) > 1 {
		ip = net.ParseIP(args[1])
		if ip == nil {
			fmt.Println("The IP address is invalid")
			printUsage()
			os.Exit(1)
		}
	}

	return hostName, ip
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		println(color.Ize(color.Green, "mDNSLocal - A simple mDNS server"))
		<-c
		println(color.Ize(color.RedBackground, "    "), "mDNS server stopped")
		os.Exit(0)
	}()

	hostnameOpts, ipOpts := getHostname()

	if hostnameOpts == "" {
		hostname, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		hostnameOpts = hostname + ".local"
		println(color.Colorize(color.Yellow, "No hostname provided\nUsing the system hostname :"), color.Colorize(color.Blue, hostname))
	}

	if ipOpts == nil {
		ipOpts = GetOutboundIP()
	}

	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		panic(err)
	}

	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		panic(err)
	}

	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		panic(err)
	}

	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		panic(err)
	}

	_, err = mdns.Server(ipv4.NewPacketConn(l4), ipv6.NewPacketConn(l6), &mdns.Config{
		LocalNames:   []string{hostnameOpts},
		LocalAddress: net.ParseIP(ipOpts.String()),
	})

	if err != nil {
		panic(err)
	}

	println(color.Ize(color.GreenBackground, "    "), "mDNS server started")
	println(color.Ize(color.BlueBackground, "    "), "Hostname :", color.Colorize(color.Green, hostnameOpts), "==> IP :", color.Colorize(color.Green, ipOpts.String()))
	select {}
}
