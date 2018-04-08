package main

import (
	"github.com/jacobsa/go-serial/serial"
	"log"
	"io"
	"fmt"
	"time"
)

var (
	portName              = "/dev/serial0"
	baudRate         uint = 19200
	byteTime              = 11.0 / baudRate
	defaultHeatTime  byte = 120
	defaultSleepTime      = time.Second / 2
	ESC              byte = 27
	column                = 0
	resumeTime       time.Time
	maxColumn        = 32
	charHeight       = 24
	prevByte         byte
	lineSpacing      = 8
	dotPrintTime     = 0.03
	dotFeedTime      = 0.0021
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

	sleep(defaultSleepTime)
	write([]byte("Dzien dobry"))

	write([]byte(" Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec condimentum non diam quis luctus. Nam ultricies dapibus massa, in ultrices enim. Praesent tempus est eu ex mollis, id luctus nisi vehicula. Ut sed ultrices neque, ac luctus nulla. Suspendisse at venenatis dui. Nunc tempor congue mauris, sit amet elementum nibh dictum a. Cras condimentum velit at vulputate condimentum. Vivamus nec orci ipsum. Morbi ut ex lorem. Sed nisl nulla, posuere non eleifend eu, ornare nec tortor. Curabitur ultrices blandit mi a tempus. Duis et ante sed libero egestas vulputate auctor a purus. Nunc neque enim, sollicitudin id efficitur sit amet, tincidunt in metus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Nunc molestie luctus ligula, a porttitor diam scelerisque at."))

}
func sleep(duration time.Duration) {
	time.Sleep(duration)
}

func write(bytes []byte) {
	for _, oneByte := range bytes {
		if oneByte != 0x13 {
			timeoutWait()
			n, err := port.Write(bytes)
			if err != nil {
				log.Fatalf("port.Write: %v", err)
			}
			fmt.Println("Wrote", n, "bytes.")
			d := float64(byteTime)
			if oneByte == '\n' || column == maxColumn {
				if prevByte == '\n' {
					d += float64(charHeight) +
						float64(lineSpacing) *
							dotFeedTime
					prevByte = oneByte
				} else {
					d += (float64(charHeight) *
						float64(dotPrintTime)) +
						(float64(lineSpacing) *
							float64(dotFeedTime))
					column = 0
					prevByte = '\n'
				}
			} else {
				column += 1
				prevByte = oneByte
			}
			timeoutSet(time.Duration(d))
		}
	}
}

/*
  # Because there's no flow control between the printer and computer,
	# special care must be taken to avoid overrunning the printer's
	# buffer.  Serial output is throttled based on serial speed as well
	# as an estimate of the device's print and feed rates (relatively
	# slow, being bound to moving parts and physical reality).  After
	# an operation is issued to the printer (e.g. bitmap print), a
	# timeout is set before which any other printer operations will be
	# suspended.  This is generally more efficient than using a delay
	# in that it allows the calling code to continue with other duties
	# (e.g. receiving or decoding an image) while the printer
	# physically completes the task.

	# Sets estimated completion time for a just-issued task.
	def timeoutSet(self, x):
		self.resumeTime = time.time() + x

	# Waits (if necessary) for the prior task to complete.
	def timeoutWait(self):
		if self.writeToStdout is False:
			while (time.time() - self.resumeTime) < 0: pass
 */

func timeoutSet(x time.Duration) {
	resumeTime = time.Now().Add(x)
}
func timeoutWait() {
	for time.Now().Sub(resumeTime) < 0 {
		print("waiting")
	}
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
