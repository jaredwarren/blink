

build:
	docker run --rm -ti -v /Users/jaredwarren/go/src/github.com/jaredwarren/blink/ws2811:/go/src/swiss rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/swiss/swiss -v swiss"

push:
	scp swiss pi@192.168.0.111:/home/pi/

pb:
	docker run --rm -ti -v /Users/jaredwarren/go/src/github.com/jaredwarren/blink/ws2811:/go/src/swiss rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/swiss/swiss -v swiss"
	scp swiss pi@192.168.0.111:/home/pi/