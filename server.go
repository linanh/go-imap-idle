package idle

import (
	"bufio"
	"errors"
	"strings"

	"github.com/linanh/go-imap"
	"github.com/linanh/go-imap/server"
)

type Handler struct {
	Command
}

func (h *Handler) Handle(conn server.Conn) error {
	cont := &imap.ContinuationReq{Info: "idling"}
	if err := conn.WriteResp(cont); err != nil {
		return err
	}
	conn.SetIdling(true)
	defer conn.SetIdling(false)

	// Wait for DONE
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return err
	}

	if strings.ToUpper(scanner.Text()) != doneLine {
		return errors.New("Expected DONE")
	}
	return nil
}

type extension struct{}

func (ext *extension) Capabilities(c server.Conn) []string {
	return []string{Capability}
}

func (ext *extension) Command(name string) server.HandlerFactory {
	if name != commandName {
		return nil
	}

	return func() server.Handler {
		return &Handler{}
	}
}

func NewExtension() server.Extension {
	return &extension{}
}
