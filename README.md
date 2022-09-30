# Generate audio

`gst-launch-1.0 audiotestsrc num-buffers=5000 ! audio/x-raw, rate=8000, channels=1 ! mulawenc ! filesink location=ulaw.raw`

# Running

`go run .`
