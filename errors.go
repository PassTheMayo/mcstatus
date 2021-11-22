package mcstatus

import "errors"

var (
	// ErrUnexpectedResponse means the server sent an unexpected response to the client
	ErrUnexpectedResponse = errors.New("received an unexpected response from the server")
	// ErrEmptyBuffer is a generic error for any read methods where the buffer array doesn't contain enough data to read the whole type
	ErrEmptyBuffer = errors.New("packet does not contain enough data to read this type")
	// ErrInvalidBoolean means the server sent a value expected as a boolean but the value was neither 0 or 1
	ErrInvalidBoolean = errors.New("cannot ReadBoolean() as value is neither 0 or 1")
	// ErrVarIntTooBig means the server sent a varint which was beyond the protocol size of a varint
	ErrVarIntTooBig = errors.New("size of VarInt exceeds maximum data size")
	// ErrNotConnected means the client attempted to send data but there was no connection to the server
	ErrNotConnected = errors.New("client attempted to send data but connection is non-existent")
	// ErrAlreadyLoggedIn means the RCON client was already logged in after a second login attempt was made
	ErrAlreadyLoggedIn = errors.New("RCON client is already logged in after a second login attempt was made")
	// ErrInvalidPassword means the password used in the RCON loggin was incorrect
	ErrInvalidPassword = errors.New("incorrect RCON password")
	// ErrNotLoggedIn means the client attempted to execute a command before a login was successful
	ErrNotLoggedIn = errors.New("RCON client attempted to send message before successful login")
	// ErrUnknownVersion means the server returned a Votifier version that is unsupported
	ErrUnknownVersion = errors.New("unsupported server Votifier version")
)
