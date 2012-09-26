package ts3

import (
	"errors"
	"strconv"
	"strings"
)

const (
	resultSeparator   = "|"
	errorPrefix       = "error "
	errorPrefixLength = len(errorPrefix)
)

var (
	ErrNotError      = errors.New("the given line is not an error line")
	ErrBadError      = errors.New("the given line was not a correctly formatted error")
	ErrInvalidResult = errors.New("invalid result")
)

type Result map[string]string

type ErrorID int64

// Results contains the result of a command sent to a TeamSpeak 3 server.
type Results struct {
	// StatusID is the `id` field returned by the "error" response of the TeamSpeak server.
	StatusID ErrorID
	// StatusMessage is the `msg` field returned by the "error" response of the TeamSpeak server.
	StatusMessage string
	// Data is all of the data returned by the TeamSpeak server.
	Data []Result
}

func decodeResult(str string) (*Results, error) {
	var newResults *Results = new(Results)

	rawResults := strings.Split(str, resultSeparator)

	if len(rawResults) == 1 && len(rawResults[0]) <= 0 {
		return newResults, nil
	}

	newResults.Data = make([]Result, len(rawResults))

	for num, rawResult := range rawResults {
		newResults.Data[num] = make(Result)

		params := strings.Split(rawResult, " ")

		for _, param := range params {
			values := strings.SplitN(param, "=", 2)

			if len(values) == 1 {
				newResults.Data[num][UnescapeTS3String(values[0])] = ""
			} else if len(values) == 2 {
				newResults.Data[num][UnescapeTS3String(values[0])] = UnescapeTS3String(values[1])
			}
		}
	}

	return newResults, nil
}

func parseError(str string) (id ErrorID, msg string, err error) {
	if !strings.HasPrefix(str, errorPrefix) {
		err = ErrNotError
		return
	}

	dataString := str[errorPrefixLength:]

	result, err := decodeResult(dataString)

	if err != nil {
		return
	}

	if len(result.Data) <= 0 {
		err = ErrBadError
		return
	}

	ourData := result.Data[0]

	if len(ourData["id"]) <= 0 || len(ourData["msg"]) <= 0 {
		err = ErrBadError
		return
	}

	tempId, err := strconv.ParseInt(ourData["id"], 10, 0)

	if err != nil {
		return
	}

	id = ErrorID(tempId)
	msg = ourData["msg"]

	err = nil

	return
}
