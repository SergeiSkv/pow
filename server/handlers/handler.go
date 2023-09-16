package handlers

import (
	"context"
	"encoding/binary"
	"io"
	"math/rand"
	"net"
	"strings"

	"github.com/SergeiSkv/pow/pkg/pow"
	"github.com/SergeiSkv/pow/server/assets"
	"github.com/pkg/errors"
)

type Handler interface {
	Handle(ctx context.Context, conn net.Conn) error
}

type handler struct {
	quotes          []assets.Quote
	targetPrefix    string
	challengeLength int
}

func NewHandler(targetPrefix string, length int, quotes []assets.Quote) Handler {
	return &handler{
		quotes:          quotes,
		targetPrefix:    targetPrefix,
		challengeLength: length,
	}
}

func (h *handler) Handle(ctx context.Context, conn net.Conn) error {
	challenge := h.generateAndSendChallenge(ctx, conn)
	if challenge == "" {
		return errors.New("failed to send challenge")
	}

	response := h.receiveResponse(ctx, conn)
	if response == "" {
		return errors.New("failed to read response")
	}

	return h.handleResponse(ctx, conn, challenge, response)
}

func (h *handler) generateAndSendChallenge(ctx context.Context, conn net.Conn) string {
	challenge := pow.GenerateChallenge(h.challengeLength)
	if h.sendStringWithContext(ctx, conn, challenge) == nil {
		return challenge
	}
	return ""
}

func (h *handler) receiveResponse(ctx context.Context, conn net.Conn) string {
	response, _ := h.readStringWithContext(ctx, conn)
	return response
}

func (h *handler) handleResponse(ctx context.Context, conn net.Conn, challenge, response string) error {
	if h.isValid(challenge, response) {
		quote := h.randomQuote()
		return h.sendStringWithContext(ctx, conn, quote.String())
	}
	return h.sendStringWithContext(ctx, conn, "Invalid PoW")
}

// Existing code for sendStringWithContext, readStringWithContext, etc.

func (h *handler) sendStringWithContext(ctx context.Context, conn net.Conn, msg string) error {
	done := make(chan error, 1)
	go func() {
		err := h.sendString(conn, msg)
		done <- err
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}

func (h *handler) readStringWithContext(ctx context.Context, conn net.Conn) (string, error) {
	done := make(chan struct {
		data string
		err  error
	}, 1)

	go func() {
		msg, err := h.readString(conn)
		done <- struct {
			data string
			err  error
		}{data: msg, err: err}
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-done:
		return result.data, result.err
	}
}

func (h *handler) sendString(conn net.Conn, msg string) error {
	length := len(msg)
	if err := binary.Write(conn, binary.BigEndian, uint16(length)); err != nil {
		return errors.Wrap(err, "failed to send message length")
	}
	_, err := conn.Write([]byte(msg))
	return err
}

func (h *handler) readString(conn net.Conn) (string, error) {
	var length uint16
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return "", errors.Wrap(err, "failed to read message length")
	}

	buffer := make([]byte, length)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func (h *handler) isValid(challenge, response string) bool {
	parts := strings.Split(response, ":")
	if len(parts) != 2 {
		return false
	}
	nonce, result := parts[0], parts[1]
	return pow.ValidatePoW(challenge, nonce, result, h.targetPrefix)
}

func (h *handler) randomQuote() assets.Quote {
	return h.quotes[rand.Intn(len(h.quotes))]
}
