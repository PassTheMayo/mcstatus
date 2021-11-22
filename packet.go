package mcstatus

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"
)

// Packet contains helper functions for reading and writing buffer packets to the TCP stream
type Packet struct {
	b io.ReadWriter
}

// NewPacket creates a new packet with no data
func NewPacket() *Packet {
	return &Packet{
		b: &bytes.Buffer{},
	}
}

// NewPacketFromReader creates a new packet from the reader
func NewPacketFromReader(r io.Reader) *Packet {
	return &Packet{
		b: &bytes.Buffer{},
	}
}

func (p *Packet) ReadBoolean() (bool, error) {
	val, err := p.ReadByte()

	if err != nil {
		return false, err
	}

	return val == 1, nil
}

func (p *Packet) WriteBoolean(val bool) error {
	if val {
		return p.WriteByte(1)
	}

	return p.WriteByte(0)
}

func (p *Packet) ReadByte() (byte, error) {
	data := make([]byte, 1)

	_, err := p.b.Read(data)

	return data[0], err
}

func (p *Packet) WriteByte(val byte) error {
	_, err := p.b.Write([]byte{val})

	return err
}

func (p *Packet) ReadBytes(length int) ([]byte, error) {
	data := make([]byte, length)

	_, err := p.b.Read(data)

	return data, err
}

func (p *Packet) WriteBytes(val []byte) error {
	_, err := p.b.Write(val)

	return err
}

func (p *Packet) ReadUInt8() (uint8, error) {
	val, err := p.ReadByte()

	return uint8(val), err
}

func (p *Packet) WriteUInt8(val uint8) error {
	return p.WriteByte(byte(val))
}

func (p *Packet) ReadInt16BE() (int16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return int16(binary.BigEndian.Uint16(data)), err
}

func (p *Packet) WriteInt16BE(val int16) error {
	data := make([]byte, 2)

	binary.BigEndian.PutUint16(data, uint16(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt16LE() (int16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return int16(binary.LittleEndian.Uint16(data)), err
}

func (p *Packet) WriteInt16LE(val int16) error {
	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, uint16(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt16BE() (uint16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint16(data), err
}

func (p *Packet) WriteUInt16BE(val uint16) error {
	data := make([]byte, 2)

	binary.BigEndian.PutUint16(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt16LE() (uint16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint16(data), err
}

func (p *Packet) WriteUInt16LE(val uint16) error {
	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt32BE() (int32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return int32(binary.BigEndian.Uint32(data)), err
}

func (p *Packet) WriteInt32BE(val int32) error {
	data := make([]byte, 4)

	binary.BigEndian.PutUint32(data, uint32(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt32LE() (int32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return int32(binary.LittleEndian.Uint32(data)), err
}

func (p *Packet) WriteInt32LE(val int32) error {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, uint32(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt32BE() (uint32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint32(data), err
}

func (p *Packet) WriteUInt32BE(val uint32) error {
	data := make([]byte, 4)

	binary.BigEndian.PutUint32(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt32LE() (uint32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint32(data), err
}

func (p *Packet) WriteUInt32LE(val uint32) error {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt64BE() (int64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return int64(binary.BigEndian.Uint64(data)), err
}

func (p *Packet) WriteInt64BE(val int64) error {
	data := make([]byte, 8)

	binary.BigEndian.PutUint64(data, uint64(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt64LE() (int64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return int64(binary.LittleEndian.Uint64(data)), err
}

func (p *Packet) WriteInt64LE(val int64) error {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt64BE() (uint64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint64(data), err
}

func (p *Packet) WriteUint64BE(val uint64) error {
	data := make([]byte, 8)

	binary.BigEndian.PutUint64(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUInt64LE() (uint64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint64(data), err
}

func (p *Packet) WriteUInt64LE(val uint64) error {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadFloat32() (float32, error) {
	val, err := p.ReadUInt32BE()

	return math.Float32frombits(val), err
}

func (p *Packet) WriteFloat32(val float32) error {
	res := math.Float32bits(val)

	return p.WriteUInt32BE(res)
}

func (p *Packet) ReadFloat64() (float64, error) {
	val, err := p.ReadUInt64BE()

	return math.Float64frombits(val), err
}

func (p *Packet) WriteFloat64(val float64) error {
	res := math.Float64bits(val)

	return p.WriteUint64BE(res)
}

func (p *Packet) ReadString() (string, error) {
	return readString(p.b)
}

func (p *Packet) WriteString(val string) error {
	return writeString(val, p.b)
}

func (p *Packet) ReadVarInt() (int32, error) {
	val, _, err := readVarInt(p.b)

	return val, err
}

func (p *Packet) WriteVarInt(val int32) error {
	_, err := writeVarInt(val, p.b)

	return err
}

func (p *Packet) ReadVarLong() (int64, error) {
	val, _, err := readVarLong(p.b)

	return val, err
}

func (p *Packet) WriteVarLong(val int64) error {
	_, err := writeVarLong(val, p.b)

	return err
}

func (p *Packet) WriteLength() error {
	data, err := ioutil.ReadAll(p.b)

	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}

	if _, err := writeVarInt(int32(len(data)), buf); err != nil {
		return err
	}

	if _, err = buf.Write(data); err != nil {
		return err
	}

	p.b = buf

	return nil
}

func (p *Packet) WriteTo(w io.Writer) (int64, error) {
	data, err := ioutil.ReadAll(p.b)

	if err != nil {
		return 0, err
	}

	n, err := w.Write(data)

	return int64(n), err
}

func (p *Packet) ReadAllBytes() ([]byte, error) {
	return ioutil.ReadAll(p.b)
}
