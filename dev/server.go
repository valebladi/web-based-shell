// package main
//
// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"syscall"
// 	"unsafe"
//
// 	"github.com/gliderlabs/ssh"
// 	"github.com/kr/pty"
// )
//
// func setWinsize(f *os.File, w, h int) {
// 	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
// 		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
// }
//
// func main() {
// 	ssh.Handle(func(s ssh.Session) {
// 		_, winCh, isPty := s.Pty()
// 		if isPty {
// 			cmd := exec.Command("ls")
// 			f, err := pty.Start(cmd)
// 			if err != nil {
// 				panic(err)
// 			}
// 			go func() {
// 				for win := range winCh {
// 					setWinsize(f, win.Width, win.Height)
// 				}
// 			}()
// 			go func() {
// 				io.Copy(f, s) // stdin
// 			}()
// 			io.Copy(s, f) // stdout
// 		} else {
// 			io.WriteString(s, fmt.Sprintf("Hello %s\n%s\n", s.User(), s.Command()))
// 		}
// 	})
//
// 	log.Println("starting ssh server on port 2222...")
// 	log.Fatal(ssh.ListenAndServe(":2222", nil))
// }

package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	// "bytes"
	// "fmt"
  //  "io"
  //  "time"
	"golang.org/x/crypto/ssh"
)

var (
	hostPrivateKeySigner ssh.Signer
)

func init() {
	keyPath := "./host_key"
	if os.Getenv("HOST_KEY") != "" {
		keyPath = os.Getenv("HOST_KEY")
	}

	hostPrivateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	hostPrivateKeySigner, err = ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		panic(err)
	}
}

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Println(conn.RemoteAddr(), "authenticate with", key.Type())
	return nil, nil
}

// func executeCmd(command, hostname string, port string, config *ssh.ClientConfig) string {
//     conn, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), config)
//     session, _ := conn.NewSession()
//     defer session.Close()
//
//     var stdoutBuf bytes.Buffer
//     session.Stdout = &stdoutBuf
//     session.Run(command)
//
//     return fmt.Sprintf("%s -> %s", hostname, stdoutBuf.String())
// }

func main() {

		log.Println("starting ssh server on port 2222...")

	config := ssh.ServerConfig{
		PublicKeyCallback: keyAuth,
	}
	config.AddHostKey(hostPrivateKeySigner)

	port := "2222"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		// From a standard TCP connection to an encrypted SSH connection
		sshConn, _, _, err := ssh.NewServerConn(conn, &config)
		if err != nil {
			panic(err)
		}

		log.Println("Connection from", sshConn.RemoteAddr())

		// for _, hostname := range hosts {
		// 		go func(hostname string, port string) {
		// 				results <- executeCmd(cmd, hostname, port, &config)
		// 		}(hostname, port)
		// }

		sshConn.Close()
	}
}
