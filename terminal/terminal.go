package terminal

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	ioctlReadTermios  = syscall.TCGETS
	ioctlWriteTermios = syscall.TCSETS
)

// ClearScreen clears the terminal screen.
//
// No parameters.
// No return types.
func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// SetStdinNonBlocking sets the standard input to non-blocking mode.
//
// It does this by getting the file descriptor of the standard input and checking the runtime operating system.
// On Windows, non-blocking mode is not supported, so an error is returned.
// On other operating systems, the file descriptor is set to non-blocking mode using the syscall.SetNonblock() function.
// The function returns an error if the operation fails.
func SetStdinNonBlocking() error {
	fd := int(os.Stdin.Fd())
	if runtime.GOOS == "windows" {
		return fmt.Errorf("non-blocking mode is not supported on Windows")
	}
	return syscall.SetNonblock(fd, true)
}

// SetRawMode sets the terminal to raw mode.
//
// It takes an integer file descriptor as a parameter and returns a pointer to a syscall.Termios struct and an error.
func SetRawMode(fd int) (*syscall.Termios, error) {
	var oldState syscall.Termios // Get the current terminal state

	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	raw := oldState
	raw.Iflag &^= syscall.ICRNL | syscall.INLCR | syscall.IXON | syscall.IXOFF
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IEXTEN
	raw.Cflag &^= syscall.CSIZE | syscall.PARENB
	raw.Cflag |= syscall.CS8
	_, _, errno = syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&raw)), 0, 0, 0)
	if errno != 0 {
		return nil, errno
	}
	return &oldState, nil
}

// ReadInput reads input from standard input and sends each character to the given channel.
//
// It takes a channel 'ch' as the first parameter which is used to send the input characters.
// The second parameter 'kill' is a channel used to stop the function.
func ReadInput(ch chan<- rune, kill <-chan bool) {
	var char [1]byte
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-kill:
			return
		case <-tick.C:
			n, err := os.Stdin.Read(char[:])
			if err != nil {
				continue
			}
			if n > 0 {
				ch <- rune(char[0])
			}
		}
	}
}

// Restore restores the terminal to the specified state.
//
// The function takes an integer file descriptor `fd` and a pointer to a `syscall.Termios` struct `oldState`
// as parameters. It returns an `error` value.
func Restore(fd int, oldState *syscall.Termios) error {
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL,
		uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(oldState)), 0, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
