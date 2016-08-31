package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
)

const (
	WAV_HEAD_SIZE = 44
)

type WavHead struct {
	RiffFlag [4]byte // "RIFF"
	FileSize uint32
	WaveFlag [4]byte // "WAVE"
	FmtFlag  [4]byte // "fmt/0"
	FmtSize  uint32
	// start of PCM WAVE FORMAT
	FormatType    uint16 // 0x0001
	Channels      uint16
	SampleRate    uint32
	BytePerSecond uint32 // samplerate * bitsPerSample * channel / 8
	Align         uint16
	BitsPerSample uint16
	// start of DATA
	DataFlag [4]byte // "data"
	DataSize uint32
}

func (head *WavHead) printf() {
	p := fmt.Println
	f := fmt.Printf
	p("------------- WAV HEAD ------------------")
	f("[%s]\n", string(head.RiffFlag[:len(head.RiffFlag)]))
	p("  FileSize:", head.FileSize)

	f("[%s]\n", string(head.WaveFlag[:len(head.WaveFlag)]))
	p("  FmtFlag:", string(head.FmtFlag[:len(head.FmtFlag)]))
	p("  FmtSize:", head.FmtSize)
	p("  FormatType:", head.FormatType)
	p("  Channels:", head.Channels)
	p("  SampleRate:", head.SampleRate)
	p("  BytePerSecond:", head.BytePerSecond)
	p("  Align:", head.Align)
	p("  BitsPerSample:", head.BitsPerSample)

	f("[%s]\n", string(head.DataFlag[:len(head.DataFlag)]))
	p("  DataSize:", head.DataSize)
}

func parseWavHead(data []byte) (*WavHead, error) {
	var head WavHead
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &head)
	if err != nil {
		log.Println("binary.Read failed:", err)
		return nil, err
	}

	return &head, nil
}

func parseWavHeadFromFile(f *os.File) (*WavHead, error) {
	data := make([]byte, WAV_HEAD_SIZE)
	n, err := f.Read(data)
	if err != nil {
		log.Println("read", err)
		return nil, err
	}

	log.Printf("read %d bytes\n", n)
	head, err := parseWavHead(data)
	if err != nil {
		log.Println("parseWavHead failed:", err)
		return nil, err
	}
	return head, nil
}

func doHead(file string) {
	trace()

	f, err := os.Open(file)
	if err != nil {
		log.Println("open", err)
		return
	}
	defer f.Close()

	head, err := parseWavHeadFromFile(f)
	head.printf()
}

func doStrip(file string) {
	trace()

	f, err := os.Open(file)
	if err != nil {
		log.Println("open", err)
		return
	}
	defer f.Close()

	head, err := parseWavHeadFromFile(f)
	head.printf()

	out, err := os.OpenFile("output.pcm", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Println("open output failed:", err)
		return
	}
	defer out.Close()

	chunkSize := 1024
	buf := make([]byte, chunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("read finished")
				break
			} else {
				log.Println("readBuf error:", err)
				break
			}
		}
		out.Write(buf[:n])
	}
}

func main() {
	log.SetOutput(os.Stdout)

	var cmd string
	var inputFile string

	flag.StringVar(&cmd, "c", "head", "parse the wave head")
	flag.StringVar(&inputFile, "f", "test.wav", "input wave file")

	flag.Parse()

	switch cmd {
	case "head":
		doHead(inputFile)
	case "strip":
		doStrip(inputFile)
	default:
		log.Printf("invalid cmd: %s\n", cmd)
	}
}

func trace() {

	p, f, l, ok := runtime.Caller(1)
	if !ok {
		return
	}
	fmt.Printf("%s (%s:%d)\n", runtime.FuncForPC(p).Name(), f, l)
}
