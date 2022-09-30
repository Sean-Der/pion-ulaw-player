//go:build !js
// +build !js

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/pkg/media"
)

func doSignaling(w http.ResponseWriter, r *http.Request) {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed: %s\n", connectionState.String())
	})

	audioTrack, err := webrtc.NewTrackLocalStaticSample(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypePCMU}, "audio", "pion")
	if err != nil {
		panic(err)
	} else if _, err := peerConnection.AddTrack(audioTrack); err != nil {
		panic(err)
	}

	go sendAudioFile(audioTrack)

	var offer webrtc.SessionDescription
	if err = json.NewDecoder(r.Body).Decode(&offer); err != nil {
		panic(err)
	}

	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	} else if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}
	<-gatherComplete

	response, err := json.Marshal(*peerConnection.LocalDescription())
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(response); err != nil {
		panic(err)
	}
}

func sendAudioFile(track *webrtc.TrackLocalStaticSample) {
	buff := make([]byte, 1024)

	file, err := os.Open("ulaw.raw")
	if err != nil {
		panic(err)
	}

	for ; true; <-time.NewTicker(time.Millisecond * 20).C {
		if _, err := file.Read(buff); err != nil {
			panic(err)
		}

		if err = track.WriteSample(media.Sample{Data: buff, Duration: time.Millisecond * 20}); err != nil {
			panic(err)
		}
	}

}

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/doSignaling", doSignaling)

	fmt.Println("Open http://localhost:8080 to access this demo")
	panic(http.ListenAndServe(":8080", nil))
}
