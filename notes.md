notes.md

# pi

## find pi
`arp -na | grep -i b8:27:eb`

## Connect
`ssh pi@192.168.0.117` pass: RV<=fun?
TODO: ip not static
<!-- TODO: add ssh key to pi so I don't need password -->
TODO: make makefile



# build
`cd blink-1`
`GOARM=6 GOARCH=arm GOOS=linux go build`

# move
scp blink-1 pi@192.168.0.117:/home/pi/
ssh -t pi@192.168.0.117 "./blink-1"




# docker
https://github.com/rpi-ws281x/rpi-ws281x-go
in ws2811

`docker build --tag rpi-ws281x-go-builder .`

`docker run --rm -ti -v "$(pwd)":/go/src/test rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/test/test -v test"`
or
`docker run --rm -ti -v "$(pwd)":/go/src/swiss rpi-ws281x-go-builder /usr/bin/qemu-arm-static /bin/sh -c "go build -o src/swiss/swiss -v swiss"`


`scp swis pi@192.168.0.117:/home/pi/`







# vs code setup go on pi
## setup vscode on pi
https://medium.com/@pythonpow/remote-development-on-a-raspberry-pi-with-ssh-and-vscode-a23388e24bc7

## setup go on pi
https://www.jeremymorgan.com/tutorials/raspberry-pi/install-go-raspberry-pi/






use this one!!!! for now
https://josemanuelpita.blogspot.com/2017/02/build-monitor-20-using-ws2812-led-strip.html





