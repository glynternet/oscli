package osc

import (
	"fmt"

	"github.com/sander/go-osc/osc"
)

// Print provides a PacketHandler that will print a recieved Packet.
// Any binary output can be decoded and printed as a string if decodeBlobs is set to true
func Print(decodeBlobs bool) PacketHandler {
	if decodeBlobs {
		return printHandler(decodedBlobsPrint)
	}
	return printHandler(rawPrint)
}

func printHandler(printFn func(*osc.Message)) PacketHandler {
	return func(packet osc.Packet) {
		switch p := packet.(type) {
		case *osc.Message:
			fmt.Printf("-- OSC Message: ")
			printFn(p)

		case *osc.Bundle:
			fmt.Println("-- OSC Bundle:")
			for i, message := range p.Messages {
				fmt.Printf("  -- OSC Message #%d: ", i+1)
				printFn(message)
			}

		default:
			fmt.Printf("Unknown packet type: %T!\n", p)
		}
	}
}

func rawPrint(msg *osc.Message) {
	fmt.Println(msg)
}

func decodedBlobsPrint(msg *osc.Message) {
	fmt.Println(msg)
	for i, a := range msg.Arguments {
		if bs, ok := a.([]byte); ok {
			fmt.Printf("element[%d]: %s\n", i, string(bs))
		}
	}
}
