package ssh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

var (
	//Allow cross-domain
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	user     = ""
	password = ""
	host     = ""
	port     = 22
	count    = 0
)

type SSHConnect struct {
	session    *ssh.Session
	stdinPipe  io.WriteCloser
	stdoutPipe io.Reader
	//stdout     *bytes.Buffer
	//stderr     *bytes.Buffer
}

func Home(w http.ResponseWriter, r *http.Request) {
	temp, e := template.ParseFiles("./internal/ssh/template/ssh.html")
	if e != nil {
		fmt.Println(e)
	}

	// Log the raw request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the request body is empty
	if len(bodyBytes) == 0 {
		http.Error(w, "Empty request body", http.StatusBadRequest)
		return
	}

	// Decode the request body into the requestBody struct
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		IP       string `json:"ip"`
	}
	if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
		http.Error(w, "Failed to decode request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Extract username, password, and IP from the request body
	username := requestBody.Username
	pass := requestBody.Password
	ip := requestBody.IP
	count++

	fmt.Println("Received credentials: username:", username, ", password:", password, ", ip:", ip, " count:", count) // Log received data

	// Assign values to global pointers (consider security implications!)
	user = username
	password = pass
	host = ip
	temp.Execute(w, nil)

	// *user = "some"  // Remove hardcoded values
	// *password = "SOME"
	// *host = "35.247.21.236"

}

// create ssh client
func CreateSSHClient(user, password, host string, port int) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		//session      *ssh.Session
		err error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			//Handling the host key
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", host, port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}
	return client, nil
	/*if session, err = client.NewSession(); err != nil {
		return nil, err
	}
	return session, nil*/
}

func RunSSH(client *ssh.Client, command string) (string, error) {
	var err error
	var session *ssh.Session
	if session, err = client.NewSession(); err == nil {
		session.StdinPipe()
		defer session.Close()
		var stdOut bytes.Buffer

		session.Stdout = &stdOut
		err = session.Run(command)
		if err != nil {
			return "", err
		}

		return string(stdOut.Bytes()), nil
	}
	return "", err

}

func NewSSHConnect(client *ssh.Client) (sshConn *SSHConnect, err error) {
	var (
		session *ssh.Session
		//stdout, stderr *bytes.Buffer
	)
	if session, err = client.NewSession(); err != nil {
		return
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("linux", 80, 40, modes); err != nil {
		return
	}

	/*stdout = new(bytes.Buffer)
	stderr = new(bytes.Buffer)*/

	pipe, _ := session.StdinPipe()
	stdoutPipe, _ := session.StdoutPipe()

	/*session.Stdout = stdout
	session.Stderr = stderr*/

	if err = session.Shell(); err != nil {
		return
	}

	return &SSHConnect{
		session:    session,
		stdinPipe:  pipe,
		stdoutPipe: stdoutPipe,
		/*stdout:    stdout,
		stderr:    stderr,*/
	}, nil
}

// Receive messages from websocket
func (s *SSHConnect) Recv(conn *websocket.Conn, quit chan int) {
	defer Quit(quit)
	var (
		bytes []byte
		err   error
	)
	for {
		if bytes, err = WsRecv(conn); err != nil {
			return
		}
		if len(bytes) > 0 {
			if _, e := s.stdinPipe.Write(bytes); e != nil {
				return
			}
		}
	}
}

func (s *SSHConnect) Output(conn *websocket.Conn) (err error) {
	done := make(chan struct{}) // Channel to signal completion

	go func() {
		defer close(done)
		for {
			data := make([]byte, 4096) // Pre-allocate memory for 4096 bytes
			if _, err = io.ReadFull(s.stdoutPipe, data); err != nil {
				if err != io.EOF {
					fmt.Println("Error reading stdout:", err)
				}
				break
			}
			if err = WsSendText(conn, data); err != nil {
				fmt.Println("Error sending data to websocket:", err)
				break
			}
		}
	}()

	select {
	case <-done:
		// Output loop exited, connection might be closed or EOF reached
	case <-time.After(10 * time.Second):
		// Timeout to prevent hanging indefinitely
		fmt.Println("Timeout waiting for SSH output")
	}

	return err
}

// test
func (s *SSHConnect) recvv(command string) {
	if _, err := s.stdinPipe.Write([]byte(command)); err != nil {
		fmt.Println(err)
	}
}

// test
func (s *SSHConnect) output() {
	tick := time.NewTicker(120 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			i := make([]byte, 1024)
			if read, err := s.stdoutPipe.Read(i); err == nil {
				i2 := string(i[:read])
				//Get head
				//split := strings.Split( i2,"\n")
				//fmt.Println(split[len(split)-1])
				fmt.Println(i2)
			}
		}
	}
}

func Quit(quit chan int) {
	quit <- 1
}

// func main() {
// 	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("./internal/ssh/static/css/"))))
// 	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("./internal/ssh/static/js/"))))

// 	http.HandleFunc("/ssh", home)
// 	http.HandleFunc("/ws/v1", wsHandle)
// 	http.ListenAndServe(":8080", nil)
// }

func WsHandle(w http.ResponseWriter, r *http.Request) {
	var (
		conn    *websocket.Conn
		client  *ssh.Client
		sshConn *SSHConnect
		err     error
	)
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}
	defer conn.Close()

	//Create ssh client
	if client, err = CreateSSHClient(user, password, host, port); err != nil {
		WsSendText(conn, []byte(err.Error()))
		return
	}
	defer client.Close()

	//connect to ssh
	if sshConn, err = NewSSHConnect(client); err != nil {
		WsSendText(conn, []byte(err.Error()))
		return
	}

	quit := make(chan int)
	go sshConn.Output(conn)
	go sshConn.Recv(conn, quit)
	<-quit
}

type WSMessage struct {
	Command string
}

func WsRecv(conn *websocket.Conn) ([]byte, error) {
	var (
		err  error
		data []byte
	)
	_, data, err = conn.ReadMessage()
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	return data, nil

}

func WsSendText(conn *websocket.Conn, b []byte) error {
	if err := conn.WriteMessage(1, b); err != nil {
		return err
	}
	return nil
}
