package ts3

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"sync"
)

type Connection interface {
	SendCommand(command *Command) (*Results, error)
	Close()
}

type interalConnection struct {
	connection  *net.TCPConn
	commandLock *sync.Mutex
	readBuffer  *bufio.Reader
	isClosed    bool
}

var (
	ErrNotTeamSpeak     = errors.New("the provided address is not a TeamSpeak 3 ServerQuery address")
	ErrConnectionClosed = errors.New("this connect has been closed and therefore cannot send data")
)

func Dial(addr *net.TCPAddr) (Connection, error) {
	println("Dialing server...")
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, err
	}

	println("Creating connection structure...")
	newConnection := &interalConnection{
		connection:  conn,
		commandLock: &sync.Mutex{},
		readBuffer:  bufio.NewReader(conn),
	}

	println("Waiting for identification line...")
	line, _, err := newConnection.readBuffer.ReadLine()

	if err != nil {
		newConnection.connection.Close()
		return nil, err
	} else if string(line) != "TS3" {
		newConnection.connection.Close()
		return nil, ErrNotTeamSpeak
	}

	println("Waiting for welcome line...")
	_, _, err = newConnection.readBuffer.ReadLine()

	if err != nil {
		newConnection.connection.Close()
		return nil, err
	}

	println("All done.")
	return newConnection, nil
}

func (conn *interalConnection) SendCommand(command *Command) (*Results, error) {
	println("Waiting for lock...")
	conn.commandLock.Lock()
	defer conn.commandLock.Unlock()

	println("Checking for closed connection...")
	if conn.isClosed {
		return nil, ErrConnectionClosed
	}

	println("Sending data...")
	_, err := conn.connection.Write([]byte(command.Encode() + "\n"))

	if err != nil {
		return nil, err
	}

	var (
		incomingLine    []byte
		incomingResults *Results
	)

	println("Waiting for lines...")
	for ; err == nil; incomingLine, _, err = conn.readBuffer.ReadLine() {
		lineStr := strings.Trim(strings.Trim(string(incomingLine), "\r"), "\n")

		if strings.HasPrefix(lineStr, errorPrefix) {
			if incomingResults == nil {
				incomingResults = new(Results)
			}

			var (
				errorId  ErrorID
				errorMsg string
			)

			errorId, errorMsg, err = parseError(lineStr)

			if err != nil {
				break
			}

			incomingResults.StatusID = errorId
			incomingResults.StatusMessage = errorMsg

			break
		} else {
			incomingResults, err = decodeResult(lineStr)
		}
	}

	if err != nil {
		return nil, err
	}

	return incomingResults, nil
}

func (conn *interalConnection) Close() {
	conn.commandLock.Lock()
	defer conn.commandLock.Unlock()

	conn.connection.Close()
	conn.isClosed = true
}
