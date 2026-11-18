package wal

import (
	"bufio"
	"databaseInGo/backend"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"os"
)

const headerSize = 8

var partialRecordError error = errors.New("record is partially torn")
var checksumMismatch error = errors.New("checksums do not match")

type WAL struct {
	File *os.File
}

type Record struct {
	checksum   uint32
	dataLength uint32
	data       string
}

type Decoder struct {
	reader *bufio.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		bufio.NewReader(r),
	}
}

func (d *Decoder) Next() (*Record, error) {
	var header [headerSize]byte
	if _, err := io.ReadFull(d.reader, header[:]); err != nil {
		if errors.Is(err, io.ErrUnexpectedEOF) {
			return nil, partialRecordError
		}
		return nil, err
	}

	cs := binary.BigEndian.Uint32(header[:4])
	dataLen := binary.BigEndian.Uint32(header[4:8])

	dataBuf := make([]byte, dataLen)
	_, err := io.ReadFull(d.reader, dataBuf)
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, partialRecordError
	}

	if err != nil {
		return nil, err
	}

	gotCs := crc32.Update(0, crc32.IEEETable, header[4:8])
	gotCs = crc32.Update(gotCs, crc32.IEEETable, dataBuf)
	if cs != gotCs {
		return nil, checksumMismatch
	}

	record := string(dataBuf)
	return &Record{
		cs,
		dataLen,
		record,
	}, nil
}
func NewWal(path string) (*WAL, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &WAL{
		f,
	}, nil
}

func EncodeStatement(ss string) []byte {
	encoding := make([]byte, 8+len(ss))
	binary.BigEndian.PutUint32(encoding[4:8], uint32(len(ss)))
	copy(encoding[8:], []byte(ss))
	crcChecksum := crc32.Checksum(encoding[4:], crc32.IEEETable)
	binary.BigEndian.PutUint32(encoding[:4], crcChecksum)
	return encoding
}

func (w *WAL) Append(stmt string) error {
	if _, err := w.File.Write(EncodeStatement(stmt)); err != nil {
		return err
	}
	if err := w.File.Sync(); err != nil {
		return err
	}
	return nil
}

func (w *WAL) Replay(mb *backend.Mb) error {
	_, err := w.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	dec := NewDecoder(w.File)
	for {
		rec, err := dec.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		_, err = mb.Execute(rec.data)
		if err != nil {
			return err
		}
	}
	return nil
}
