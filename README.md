# Generate audio

`gst-launch-1.0 audiotestsrc num-buffers=5000 ! mulawenc ! filesink location=ulaw.raw`

# Running

`go run .`
