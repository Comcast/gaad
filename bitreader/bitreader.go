package bitreader

import (
	"fmt"
	"math"
)

type BitReader struct {
	bytes            []byte
	currentByte      byte
	currentByteIndex uint
	posInCurrentByte uint
	length           uint
}

var netmasks = []byte{
	0x01, // 0000 0001
	0x02, // 0000 0010
	0x04, // 0000 0100
	0x08, // 0000 1000
	0x10, // 0001 0000
	0x20, // 0010 0000
	0x40, // 0100 0000
	0x80, // 1000 0000
}

func NewBitReader(input []byte) *BitReader {
	return &BitReader{
		bytes:            input,
		currentByte:      input[0],
		currentByteIndex: 0,
		posInCurrentByte: 7,
		length:           uint(len(input)),
	}
}

// Return the total number of bits left in the stream
func (p *BitReader) BitsLeft() uint {
	return ((p.BytesLeft() - 1) * 8) + (uint(p.posInCurrentByte) + 1)
}

// Return the number of bytes left (even if partially read)
func (p *BitReader) BytesLeft() uint {
	bytesLeft := p.length - p.currentByteIndex
	if bytesLeft < 0 {
		return 0
	}
	return uint(bytesLeft)
}

func (p *BitReader) ReadBitAsBool() (bool, error) {
	val, err := p.ReadBit()
	if err != nil {
		return false, err
	}
	if val == 0 {
		return false, nil
	}
	return true, nil
}

// Return the number of bits as an unsigned integer
func (p *BitReader) ReadBitsAsUInt(n uint) (uint, error) {
	var result uint
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		result += uint(bit) << uint(n-i-1)
	}
	return result, nil
}

// Return the number of bits as an unsigned integer
func (p *BitReader) ReadBitsAsUInt8(n uint) (uint8, error) {
	var result uint8
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		result += uint8(bit) << uint(n-i-1)
	}
	return result, nil
}

// Return the number of bits as an unsigned integer
func (p *BitReader) ReadBitsAsUInt32(n uint) (uint32, error) {
	var result uint32
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		result += uint32(bit) << uint(n-i-1)
	}
	return result, nil
}

// Return the number of bits as an unsigned integer
func (p *BitReader) ReadBitsAsUInt16(n uint) (uint16, error) {
	var result uint16
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		result += uint16(bit) << uint(n-i-1)
	}
	return result, nil
}

// Return the number of bits as a signed integer
func (p *BitReader) ReadBitsAsInt(n uint) (int, error) {
	var result int
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		result += int(bit) << uint(n-i-1)
	}
	return result, nil
}

// Return n number of bits into a byte array
func (p *BitReader) ReadBitsToByteArray(n uint) ([]byte, error) {
	result := make([]byte, int(math.Ceil(float64(n)/8)))

	temp := make([]byte, n)
	for i := uint(0); i < n; i++ {
		bit, err := p.ReadBit()
		if err != nil {
			return nil, err
		}
		temp[i] = bit
	}

	bitmask := 0
	for i := n; i > 0; i-- {
		index := len(result) - (bitmask / 8) - 1
		result[index] |= temp[i-1] << uint(bitmask%8)
		bitmask++
	}

	return result, nil
}

// Return n number of bits
func (p *BitReader) ReadBits(n uint) (byte, error) {
	if n > 8 {
		fmt.Printf("ReadBits(n) can only handle upto 8 bits, use ReadBitsToByteArray(n)")
	}

	var r byte

	for i := n; i > 0; i-- {
		r <<= 1
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		r |= bit
	}
	return r, nil
}

// Return the next bit from the buffer
func (p *BitReader) ReadBit() (byte, error) {
	if p.BitsLeft() == 0 {
		return 0, fmt.Errorf("Not enough bits left to read")
	}
	r := (p.currentByte & netmasks[p.posInCurrentByte]) >> p.posInCurrentByte
	p.SkipBits(1)
	return r, nil
}

// Return n number of bytes read from the buffer
func (p *BitReader) ReadBytes(n uint) ([]byte, error) {
	arr := make([]byte, n)

	if p.BytesLeft() < n {
		return nil, fmt.Errorf("Not enough bytes left to read")
	} else {
		copy(arr, p.bytes[p.currentByteIndex:p.currentByteIndex+n])
		p.currentByteIndex += n
		p.posInCurrentByte = 7
	}

	return arr, nil
}

// Read Unsigned Exp-Golomb
func (p *BitReader) ReadUE() (uint32, error) {
	zeros := 0
	val := uint32(0)

	// Count leading zeros
	for {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		if bit == 0 {
			zeros++
		} else {
			break
		}
	}

	if zeros == 0 {
		return 0, nil
	} else {
		val = 1
	}

	// Shift bits
	for {
		bit, err := p.ReadBit()
		if err != nil {
			return 0, err
		}
		val <<= 1
		val |= uint32(bit)
		zeros--
		if zeros == 0 {
			return val - 1, nil // Subtract onebecause we stole a bit for 0.
		}
	}

	// Should not get here
	return 0, ErrExpGolombDecode
}

// Read Signed Exp-Golomb
func (p *BitReader) ReadSE() (int32, error) {
	u, err := p.ReadUE()
	if err != nil {
		return 0, err
	}
	s := int32(u)
	sign := ((s & 0x1) << 1) - 1
	return ((s >> 1) + (s & 0x1)) * sign, nil
}

// Return n number of bits from the buffer, do not adv. the cursor
func (p *BitReader) PeekBits(n uint) (uint, error) {
	var r uint

	currentByte := p.currentByte
	posInCurrentByte := p.posInCurrentByte
	currentByteIndex := p.currentByteIndex
	length := p.length
	for i := n; i > 0; i-- {
		r <<= 1
		bit := (currentByte & netmasks[posInCurrentByte]) >> posInCurrentByte
		r |= uint(bit)

		if posInCurrentByte > 0 {
			posInCurrentByte--
		} else {
			currentByteIndex++
			if currentByteIndex <= length-1 {
				currentByte = p.bytes[currentByteIndex]
			} else {
				currentByte = byte(0)
			}
			posInCurrentByte = 7
		}
	}

	return r, nil
}

// Return the next bit from the buffer, do not adv. the cursor
func (p *BitReader) PeekBit() (byte, error) {
	if p.BitsLeft() == 0 {
		return 0, fmt.Errorf("Not enough bits left to read")
	}
	r := (p.currentByte & netmasks[p.posInCurrentByte]) >> p.posInCurrentByte
	return r, nil
}

// Skip n number of bits in the buffer
func (p *BitReader) SkipBits(n uint) error {
	if p.BitsLeft() < n {
		return fmt.Errorf("Not enough bits left to skip")
	} else {
		for i := n; i > 0; i-- {
			if p.posInCurrentByte > 0 {
				p.posInCurrentByte--
			} else {
				p.currentByteIndex++
				if p.currentByteIndex <= p.length-1 {
					p.currentByte = p.bytes[p.currentByteIndex]
				} else {
					p.currentByte = byte(0)
				}
				p.posInCurrentByte = 7
			}
		}
	}

	return nil
}

// Skip n number of bytes in the buffer
func (p *BitReader) SkipBytes(n uint) error {
	if p.BytesLeft() < n {
		return fmt.Errorf("Not enough bytes left to skip")
	} else {
		p.currentByteIndex += n
		p.posInCurrentByte = 7
		if p.currentByteIndex < p.length {

			p.currentByte = p.bytes[p.currentByteIndex]
		} else {
			p.currentByte = byte(0)
		}
	}
	return nil
}

// Return if there's a bit left in the stream
func (p *BitReader) HasBitLeft() bool {
	return p.BitsLeft() > 0
}

// Return if there is a byte left in the stream
func (p *BitReader) HasByteLeft() bool {
	return p.HasBytesLeft(0)
}

// Return if there are n bytes left in the stream
func (p *BitReader) HasBytesLeft(n uint) bool {
	return p.BytesLeft() > n
}

// Reset the stream reader back to the start of the buffer
func (p *BitReader) Reset() {
	p.currentByteIndex = 0
	p.currentByte = p.bytes[0]
	p.posInCurrentByte = 7
}

// Pefrom byte alignment (skip any remaining bits of current byte)
func (p *BitReader) ByteAlign() {
	if p.posInCurrentByte != 7 {
		p.SkipBits(p.posInCurrentByte + 1)
	}
}

func (p *BitReader) ByteOffset() uint64 {
	return uint64(p.currentByteIndex)
}
