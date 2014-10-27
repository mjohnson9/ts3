package ts3

import (
	"bufio"
	"errors"
	"net"
	"strings"
)

// A Connection is a connection between you and a TeamSpeak server.
type Connection interface {
	// SendCommand sends the given command to this connection.
	// It will return an error if the connection is closed or it encounters any other problems.
	SendCommand(command *Command) (*Results, error)
	// Close closes this connection. Whether or not Close closes the underlying socket is up to the implementation.
	Close()
}

type interalConnection struct {
	connection  *net.TCPConn
	commandChan chan internalCommandRequest
}

type internalCommandRequest struct {
	Command      Command
	responseChan chan internalCommandResponse
}

type internalCommandResponse struct {
	Error   error
	Results *Results
}

var (
	ErrNotTeamSpeak = errors.New("the provided address is not a TeamSpeak 3 ServerQuery address")
)

// This dials the TeamSpeak 3 ServerQuery interface at the given address.
func Dial(addr *net.TCPAddr) (Connection, error) {
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return nil, err
	}

	newConnection := &interalConnection{
		connection:  conn,
		commandChan: make(chan internalCommandRequest),
	}

	bufferedConnection := bufio.NewReader(newConnection.connection)

	// Read the initial line and remove the \n from it.
	line, err := bufferedConnection.ReadString('\n')
	line = strings.TrimSuffix(line, "\n")

	if err != nil {
		newConnection.Close()
		return nil, err
	} else if line != "TS3" {
		newConnection.Close()
		return nil, ErrNotTeamSpeak
	}

	// Read and discard the welcome line
	_, err = bufferedConnection.ReadString('\n')

	if err != nil {
		newConnection.Close()
		return nil, err
	}

	go newConnection.process()

	return newConnection, nil
}

func (conn *interalConnection) process() {
	bufferedConnection := bufio.NewReader(conn.connection)

	for toProcess := range conn.commandChan {
		_, err := conn.connection.Write([]byte(toProcess.Command.Encode() + "\n"))

		if err != nil {
			toProcess.responseChan <- internalCommandResponse{
				Error: err,
			}
			close(toProcess.responseChan)
			continue
		}

		incomingResults := new(Results)

		for {
			var lineStr string
			lineStr, err = bufferedConnection.ReadString('\n')
			if err != nil {
				break
			}
			lineStr = strings.TrimPrefix(strings.TrimSuffix(lineStr, "\n"), "\r")

			if len(lineStr) > 0 {
				if strings.HasPrefix(lineStr, errorPrefix) {
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
			} else {
				break
			}
		}

		if err != nil {
			toProcess.responseChan <- internalCommandResponse{
				Error: err,
			}
			close(toProcess.responseChan)
			continue
		}

		toProcess.responseChan <- internalCommandResponse{
			Results: incomingResults,
		}
		close(toProcess.responseChan)
	}
}

func (conn *interalConnection) SendCommand(command *Command) (*Results, error) {
	responseChan := make(chan internalCommandResponse)
	conn.commandChan <- internalCommandRequest{
		Command:      *command,
		responseChan: responseChan,
	}

	response := <-responseChan
	return response.Results, response.Error
}

func (conn *interalConnection) Close() {
	close(conn.commandChan)
	conn.connection.Close()
}
