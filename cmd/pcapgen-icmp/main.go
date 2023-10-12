package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"git.cyberfire.ninja/devs/pcapgen/pkg/pcapwriter"
)

func sink(r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	buf := make([]byte, 4096)
	for {
		if _, err := r.Read(buf); (err == io.EOF) || (err == io.ErrClosedPipe) {
			return
		} else if err != nil {
			log.Println("Read error:", err)
		}
	}
}

func usage() {
	out := flag.CommandLine.Output()
	fmt.Fprintf(out, "Usage: cat script.txt | %s > out.pcap\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Reads from stdin, writes a PCAP file to stdout.")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "# Example script")
	fmt.Fprintln(out, "##################")
	fmt.Fprintln(out, "# Client query (4 bytes)")
	fmt.Fprintln(out, "C: 3e 29 008a")
	fmt.Fprintln(out, "# Delay 12 seconds")
	fmt.Fprintln(out, "Z: 12s")
	fmt.Fprintln(out, "# Server response (2 bytes + 1 byte)")
	fmt.Fprintln(out, "S: 7f c3")
	fmt.Fprintln(out, "S: 04")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	begin := time.Date(2010, 2, 22, 22, 57, 23, 71877000, time.UTC)
	pcap, err := pcapwriter.NewWriter(os.Stdout, begin, 20*time.Millisecond)
	if err != nil {
		log.Fatal(err)
	}
	pcap.WriteStandardHeader()

	cli, srv := pcapwriter.NewICMPv4Taps(pcap, 11, 55)

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go sink(srv, wg)
	go sink(cli, wg)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			// Skip blank lines
			continue
		} else if line[0] == '#' {
			// Skip comment lines
			continue
		}

		directive, data, found := strings.Cut(line, ":")
		if !found {
			log.Fatal("line")
		}
		data = strings.ReplaceAll(data, " ", "")
		data = strings.ReplaceAll(data, "\t", "")
		switch directive {
		case "C":
			if buf, err := hex.DecodeString(data); err != nil {
				log.Fatal(err)
			} else {
				cli.Write(buf)
			}
		case "S":
			if buf, err := hex.DecodeString(data); err != nil {
				log.Fatal(err)
			} else {
				srv.Write(buf)
			}
		case "Z":
			if d, err := time.ParseDuration(data); err != nil {
				log.Fatal(err)
			} else {
				pcap.Sleep(d)
			}
		default:
			log.Fatal("Unknown directive:", directive)
		}
	}

	cli.Close()
	srv.Close()

	wg.Wait()
}
