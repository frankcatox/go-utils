package process

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	mfdCloexec     = 0x0001
	memfdCreateX64 = 319
	fork           = 57
)

// Execute run
func Execute(command ...string) (int, error) {
	signal.Ignore(syscall.SIGHUP)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Env = os.Environ()
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

// RunFromMemory shellcode
func RunFromMemory(runName string, runDir string, buffer []byte, env []string, args ...string) (int, error) {
	fdName := ""
	fd, _, _ := syscall.Syscall(uintptr(memfdCreateX64), uintptr(unsafe.Pointer(&fdName)), uintptr(mfdCloexec), 0)
	_, _ = syscall.Write(int(fd), buffer)
	fdPath := fmt.Sprintf("/proc/self/fd/%d", fd)

	child, _, _ := syscall.Syscall(uintptr(fork), 0, 0, 0)
	switch child {
	case 0:
		break
	case 1:
		// Fork failed!
		return 0, errors.New("fork failed")
	default:
		// Parent exiting...
		return int(child), nil
	}

	_ = syscall.Umask(0)
	_, _ = syscall.Setsid()
	_ = syscall.Chdir(runDir)

	file, _ := os.OpenFile("/dev/null", os.O_RDWR, 0)
	syscall.Dup2(int(file.Fd()), int(os.Stdin.Fd()))
	syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
	syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
	file.Close()

	progWithArgs := append([]string{runName}, args...)
	pid := os.Getpid()
	err := syscall.Exec(fdPath, progWithArgs, env)
	return pid, err
}

// Daemon process
func Daemon(nochdir, noclose int) int {
	var ret, ret2 uintptr
	var err syscall.Errno
	darwin := runtime.GOOS == "darwin"

	// already a daemon
	if syscall.Getppid() == 1 {
		return 0
	}

	// fork off the parent process
	ret, ret2, err = syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		return -1
	}

	// failure
	if ret2 < 0 {
		os.Exit(-1)
	}

	// handle exception for darwin
	if darwin && ret2 == 1 {
		ret = 0
	}

	// if we got a good PID, then we call exit the parent process.
	if ret > 0 {
		os.Exit(0)
	}

	signal.Ignore(syscall.SIGHUP)

	/* Change the file mode mask */
	syscall.Umask(0)

	// create a new SID for the child process
	sRet, sErrno := syscall.Setsid()
	if sErrno != nil {
		fmt.Println("Error: syscall.Setsid errno: ", sErrno)
	}
	if sRet < 0 {
		return -1
	}
	if nochdir == 0 {
		os.Chdir("/tmp")
	}

	if noclose == 0 {
		f, e := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if e == nil {
			fd := f.Fd()
			syscall.Dup2(int(fd), int(os.Stdin.Fd()))
			syscall.Dup2(int(fd), int(os.Stdout.Fd()))
			syscall.Dup2(int(fd), int(os.Stderr.Fd()))
		}
	}

	return 0
}
