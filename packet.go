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

func (p *Packet) ReadUnsignedByte() (uint8, error) {
	val, err := p.ReadByte()

	return uint8(val), err
}

func (p *Packet) WriteUnsignedByte(val uint8) error {
	return p.WriteByte(byte(val))
}

func (p *Packet) ReadShort() (int16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return int16(binary.BigEndian.Uint16(data)), err
}

func (p *Packet) WriteShort(val int16) error {
	data := make([]byte, 2)

	binary.BigEndian.PutUint16(data, uint16(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadShortLE() (int16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return int16(binary.LittleEndian.Uint16(data)), err
}

func (p *Packet) WriteShortLE(val int16) error {
	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, uint16(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedShort() (uint16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint16(data), err
}

func (p *Packet) WriteUnsignedShort(val uint16) error {
	data := make([]byte, 2)

	binary.BigEndian.PutUint16(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedShortLE() (uint16, error) {
	data := make([]byte, 2)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint16(data), err
}

func (p *Packet) WriteUnsignedShortLE(val uint16) error {
	data := make([]byte, 2)

	binary.LittleEndian.PutUint16(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadInt() (int32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return int32(binary.BigEndian.Uint32(data)), err
}

func (p *Packet) WriteInt(val int32) error {
	data := make([]byte, 4)

	binary.BigEndian.PutUint32(data, uint32(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadIntLE() (int32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return int32(binary.LittleEndian.Uint32(data)), err
}

func (p *Packet) WriteIntLE(val int32) error {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, uint32(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedInt() (uint32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint32(data), err
}

func (p *Packet) WriteUnsignedInt(val uint32) error {
	data := make([]byte, 4)

	binary.BigEndian.PutUint32(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedIntLE() (uint32, error) {
	data := make([]byte, 4)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint32(data), err
}

func (p *Packet) WriteUnsignedIntLE(val uint32) error {
	data := make([]byte, 4)

	binary.LittleEndian.PutUint32(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadLong() (int64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return int64(binary.BigEndian.Uint64(data)), err
}

func (p *Packet) WriteLong(val int64) error {
	data := make([]byte, 8)

	binary.BigEndian.PutUint64(data, uint64(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadLongLE() (int64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return int64(binary.LittleEndian.Uint64(data)), err
}

func (p *Packet) WriteLongLE(val int64) error {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, uint64(val))

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedLong() (uint64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return binary.BigEndian.Uint64(data), err
}

func (p *Packet) WriteUnsignedLong(val uint64) error {
	data := make([]byte, 8)

	binary.BigEndian.PutUint64(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadUnsignedLongLE() (uint64, error) {
	data := make([]byte, 8)

	_, err := p.b.Read(data)

	return binary.LittleEndian.Uint64(data), err
}

func (p *Packet) WriteUnsignedLongLE(val uint64) error {
	data := make([]byte, 8)

	binary.LittleEndian.PutUint64(data, val)

	_, err := p.b.Write(data)

	return err
}

func (p *Packet) ReadFloat() (float32, error) {
	val, err := p.ReadUnsignedInt()

	return math.Float32frombits(val), err
}

func (p *Packet) WriteFloat(val float32) error {
	res := math.Float32bits(val)

	return p.WriteUnsignedInt(res)
}

func (p *Packet) ReadDouble() (float64, error) {
	val, err := p.ReadUnsignedLong()

	return math.Float64frombits(val), err
}

func (p *Packet) WriteDouble(val float64) error {
	res := math.Float64bits(val)

	return p.WriteUnsignedLong(res)
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
