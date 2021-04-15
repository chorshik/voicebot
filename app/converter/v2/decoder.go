package v2

import (
	"github.com/hajimehoshi/go-mp3"
	"gopkg.in/hraban/opus.v2"
	"io/ioutil"
	"log"
	"net/http"
)

// Path ...
const (
	workDir    = "voice_messages/"
	channels   = 1 // mono; 2 for stereo
	bufferSize = 1000
)

// DecodeMP3 ...
func DecodeMP3(resp *http.Response) *mp3.Decoder {
	f := resp.Body
	decode, err := mp3.NewDecoder(f)
	if err != nil {
		log.Print(err)
	}

	return decode
}

// EncodeToOPUS ...
func EncodeToOPUS(decode *mp3.Decoder, resp *http.Response) []byte {
	sampleRate := decode.SampleRate()
	sampleRate = 48000

	byteSlice, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	intSlice := make([]int16, len(byteSlice))
	for i, b := range byteSlice {
		intSlice[i] = int16(b)
	}

	enc, err := opus.NewEncoder(sampleRate, channels, opus.AppVoIP)
	if err != nil {
		log.Print(err)
	}

	var pcm []int16 = intSlice

	//frameSize := len(pcm)
	//frameSizeMs := float32(frameSize) / channels * 1000 / float32(sampleRate)

	data := make([]byte, bufferSize)

	n, err := enc.Encode(pcm, data)
	if err != nil {
		log.Print(err)
	}
	data = data[:n]

	return data

}
