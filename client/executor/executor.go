package executor

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/SergeiSkv/pow/pkg/pow"
	"github.com/pkg/errors"
)

type Executor interface {
	Execute(ctx context.Context) (string, error)
}

type executor struct {
	conn         net.Conn
	targetPrefix string
}

func NewExecutor(host string, targetPrefix string) (Executor, error) {
	conn, err := establishConnection(host)
	if err != nil {
		return nil, errors.Wrap(err, "establishConnection")
	}
	return &executor{conn: conn, targetPrefix: targetPrefix}, nil
}

func (e *executor) Execute(ctx context.Context) (string, error) {
	challenge, err := e.readChallenge(ctx)
	if err != nil {
		return "", err
	}

	nonce, result := pow.SolvePoW(challenge, e.targetPrefix)
	return e.sendPoWResult(ctx, nonce, result)
}

func establishConnection(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, errors.Wrap(err, "dial to server")
	}
	return conn, nil
}

func (e *executor) performNetworkOperationWithContext(ctx context.Context, op func() error) error {
	done := make(chan error, 1)
	go func() {
		done <- op()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (e *executor) sendPoWResult(ctx context.Context, nonce string, result string) (string, error) {
	msg := fmt.Sprintf("%s:%s", nonce, result)

	if err := e.performNetworkOperationWithContext(ctx, func() error {
		return e.sendString(msg)
	}); err != nil {
		return "", err
	}

	response, err := e.readString()
	if err != nil {
		return "", errors.Wrap(err, "read server response")
	}

	return response, nil
}

func (e *executor) sendString(msg string) error {
	length := len(msg)
	if err := binary.Write(e.conn, binary.BigEndian, uint16(length)); err != nil {
		return errors.Wrap(err, "failed to send message length")
	}
	_, err := e.conn.Write([]byte(msg))
	return errors.Wrap(err, "failed to send message")
}

func (e *executor) readString() (string, error) {
	var length uint16
	if err := binary.Read(e.conn, binary.BigEndian, &length); err != nil {
		return "", errors.Wrap(err, "failed to read message length")
	}

	buffer := make([]byte, length)
	_, err := io.ReadFull(e.conn, buffer)
	if err != nil {
		return "", errors.Wrap(err, "failed to read message")
	}

	return string(buffer), nil
}

func (e *executor) readChallenge(ctx context.Context) (string, error) {
	var challenge string
	err := e.performNetworkOperationWithContext(ctx, func() error {
		var readErr error
		challenge, readErr = e.readString()
		return readErr
	})
	if err != nil {
		return "", errors.Wrap(err, "read challenge")
	}
	return challenge, nil
}
