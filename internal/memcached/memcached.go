package memcached

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/maxrasky/crema/internal/model"
)

const (
	defaultTimeout = 300 * time.Millisecond
)

var (
	crlf           = []byte("\r\n")
	resultStored   = []byte("STORED\r\n")
	resultDeleted  = []byte("DELETED\r\n")
	resultNotFound = []byte("NOT_FOUND\r\n")
	resultEnd      = []byte("END\r\n")
)

type Client struct {
	nc net.Conn
	rw *bufio.ReadWriter
}

func New(address string) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	nc, err := dial(addr)
	if err != nil {
		return nil, err
	}

	return &Client{
		nc: nc,
		rw: bufio.NewReadWriter(bufio.NewReader(nc), bufio.NewWriter(nc)),
	}, nil
}

func (c *Client) Set(item *model.Item) error {
	_, err := fmt.Fprintf(c.rw, "%s %s %d %d %d\r\n", "set", item.Key, item.Flags, item.Expiration, len(item.Value))
	if err != nil {
		return err
	}

	if _, err = c.rw.Write(item.Value); err != nil {
		return err
	}
	if _, err := c.rw.Write(crlf); err != nil {
		return err
	}
	if err := c.rw.Flush(); err != nil {
		return err
	}
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return err
	}

	switch {
	case bytes.Equal(line, resultStored):
		return nil
	default:
		return fmt.Errorf("unexpected response 'set': %s", string(line))
	}
}

func (c *Client) Delete(key string) error {
	_, err := fmt.Fprintf(c.rw, "delete %s\r\n", key)
	if err != nil {
		return err
	}
	if err := c.rw.Flush(); err != nil {
		return err
	}
	line, err := c.rw.ReadSlice('\n')
	if err != nil {
		return err
	}

	switch {
	case bytes.Equal(line, resultDeleted), bytes.Equal(line, resultNotFound):
		return nil
	default:
		return fmt.Errorf("unexpected response: %s", string(line))
	}

}

func (c *Client) Get(key string) (*model.Item, error) {
	if _, err := fmt.Fprintf(c.rw, "get %s\r\n", key); err != nil {
		return nil, err
	}

	if err := c.rw.Flush(); err != nil {
		return nil, err
	}

	item, err := parseResponse(c.rw.Reader)
	switch {
	case err != nil:
		return nil, err
	case item == nil:
		return nil, model.ErrNotFound
	default:
		return item, nil
	}
}

func parseResponse(r *bufio.Reader) (*model.Item, error) {
	var res *model.Item

	for {
		line, err := r.ReadSlice('\n')
		if err != nil {
			return res, err
		}
		if bytes.Equal(line, resultEnd) {
			return res, nil
		}

		item := new(model.Item)
		size, err := scanResponseLine(line, item)
		if err != nil {
			return res, err
		}
		item.Value = make([]byte, size+2)
		_, err = io.ReadFull(r, item.Value)
		if err != nil {
			item.Value = nil
			return res, err
		}

		item.Value = item.Value[:size]
		res = item
	}
}

func scanResponseLine(line []byte, it *model.Item) (size int, err error) {
	pattern := "VALUE %s %d %d\r\n"
	_, err = fmt.Sscanf(string(line), pattern, &it.Key, &it.Flags, &size)
	if err != nil {
		return -1, err
	}

	return size, nil
}

func dial(addr net.Addr) (net.Conn, error) {
	return net.DialTimeout(addr.Network(), addr.String(), defaultTimeout)
}
