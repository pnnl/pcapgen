package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"git.cyberfire.ninja/devs/pcapgen/pkg/pcapwriter"
)

func junk(n int) []byte {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		log.Fatal(err)
	}
	return buf
}

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
	fmt.Fprintf(out, "Usage: %s FILE [FILE...]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Runs a netarch 25000 session, transferring all listed files, multiplexed.")
}

func converse(cli, srv io.Writer, filenames []string) {
	files := make([]io.ReadCloser, len(flag.Args()))

	// Pretend there's some sort of undecodable initialization handshake
	cli.Write(junk(0x40))
	srv.Write(junk(0x20))
	cli.Write(junk(12))

	// Server: request files
	for i, name := range filenames {
		f, err := os.Open(name)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := f.Stat()
		if err != nil {
			log.Fatal(err)
		}

		srv.Write(XferBeginPacket(uint8(i), uint32(fi.Size()), name))
		cli.Write(AckPacket())
		files[i] = f
	}

	filesLeft := len(files)
	for filesLeft > 0 {
		filesLeft = 0
		for i, f := range files {
			buf, err := XferPacket(uint8(i), f)
			if err == io.EOF {
				continue
			} else if err != nil {
				log.Fatal(err)
			}
			cli.Write(buf)
			srv.Write(AckPacket())
			filesLeft += 1
		}
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if len(flag.Args()) < 1 {
		flag.Usage()
		return
	}

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

	converse(cli, srv, flag.Args())

	cli.Close()
	srv.Close()
	wg.Wait()
}
