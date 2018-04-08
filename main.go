package main

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"fmt"
	"time"
	"io"
)

var (
	portName              = "/dev/serial0"
	baudRate         uint = 19200
	byteTime              = 11.0 / baudRate
	defaultHeatTime  byte = 120
	defaultSleepTime      = time.Second / 2
	ESC              byte = 27
)
var port io.ReadWriteCloser

func main() {

	// Set up options.
	options := serial.OpenOptions{
		PortName:        portName,
		BaudRate:        baudRate,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	// Open the port.
	var err error
	port, err = serial.Open(options)
	if err != nil {
		log.Fatalf("serial.Open: %v", err)
	}

	// Make sure to close it later.
	defer port.Close()

	sleep(defaultSleepTime)
	wake()
	reset()

	write([]byte{
		ESC,
		55,
		11,
		defaultHeatTime,
		40})

	printDensity := byte(10)
	printBreakTime := byte(2)

	write([]byte{18,
		35,
		(printBreakTime << 5) | printDensity})
	//dotPrintTime := 0.03
	//dotFeedTime := 0.0021

	write([]byte("Dzien dobry"))

}
func sleep(duration time.Duration) {
	time.Sleep(duration)
}
func write(bytes []byte) {
	n, err := port.Write(bytes)
	if err != nil {
		log.Fatalf("port.Write: %v", err)
	}

	fmt.Println("Wrote", n, "bytes.")
}

/*
    def wake(self):
        self.timeoutSet(0)
        self.writeBytes(255)
        if self.firmwareVersion >= 264:
            time.sleep(0.05)  # 50 ms
            self.writeBytes(27, 118, 0)  # Sleep off (important!)
        else:
            for i in range(10):
                self.writeBytes(27)
                self.timeoutSet(0.1)
 */

func wake() {
	write([]byte{255})
	sleep(defaultSleepTime)
	write([]byte{27, 118, 0})
}

/*
 def reset(self):
        self.writeBytes(27, 64)  # Esc @ = init command
        self.prevByte = '\n'  # Treat as if prior line is blank
        self.column = 0
        self.maxColumn = 32
        self.charHeight = 24
        self.lineSpacing = 6
        self.barcodeHeight = 50
        if self.firmwareVersion >= 264:
            # Configure tab stops on recent printers
            self.writeBytes(27, 68)  # Set tab stops
            self.writeBytes(4, 8, 12, 16)  # every 4 columns,
            self.writeBytes(20, 24, 28, 0)  # 0 is end-of-list.
 */
func reset() {
	write([]byte{ESC, '@'})
	write([]byte{ESC, 68})
	write([]byte{4, 8, 12, 16})
	write([]byte{20, 24, 28, 0})
}
