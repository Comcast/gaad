package bitreader

import (
	"encoding/hex"
	"reflect"
	"testing"
)

var zeroToTenHex = "A64298E2048A16"

func TestReadUE(t *testing.T) {
	byteArray, _ := hex.DecodeString(zeroToTenHex)
	parser := NewBitReader(byteArray)
	for i := uint32(0); i <= 10; i++ {
		got, err := parser.ReadUE()
		if err != nil {
			t.Errorf("TestReadUE - err from index=%d: %v", i, err)
		} else {
			if i != got {
				t.Errorf("TestReadUE - wanted %d, got %d", i, got)
			}
		}
	}
}

func TestReadSE(t *testing.T) {
	byteArray, _ := hex.DecodeString(zeroToTenHex)
	parser := NewBitReader(byteArray)

	wants := []int32{0, 1, -1, 2, -2, 3, -3, 4, -4, 5, -5}

	for i := 0; i <= 10; i++ {
		got, err := parser.ReadSE()
		if err != nil {
			t.Errorf("TestReadSE - err from index=%d: %v", i, err)
		} else {
			want := wants[i]
			if want != got {
				t.Errorf("TestReadSE - index=%d, wanted %d, got %d", i, wants, got)
			}
		}
	}
}

func TestReadBit1(t *testing.T) {
	byteArray := []byte{0xaa} // 1010 1010
	reader := NewBitReader(byteArray)

	bit, err := reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(1) {
		t.Errorf("bit (%d) must equal 1", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(0) {
		t.Errorf("bit (%d) must equal 0", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(1) {
		t.Errorf("bit (%d) must equal 1", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(0) {
		t.Errorf("bit (%d) must equal 0", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(1) {
		t.Errorf("bit (%d) must equal 1", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(0) {
		t.Errorf("bit (%d) must equal 0", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(1) {
		t.Errorf("bit (%d) must equal 1", bit)
	}

	bit, err = reader.ReadBit()
	if err != nil {
		t.Errorf("err (%s) must be nil", err.Error())
	}
	if bit != byte(0) {
		t.Errorf("bit (%d) must equal 0", bit)
	}
}

func TestReadBitCrossoverBytes(t *testing.T) {

	// Input: 0101 0101 0101 0101
	byteArray := []byte{0x55, 0x55}
	reader := NewBitReader(byteArray)

	// Test 1: Skip 7 bytes (0101 010)
	reader.SkipBits(7)

	// Test 2: Read 5 bytes (10101 --> 00010101)
	r1, _ := reader.ReadBits(5)
	if r1 != byte(0x15) {
		t.Errorf("r1 (%d) must be equal to 0x15", r1)
	}

	// Verify: BytesLeft/BitsLeft
	if reader.BytesLeft() != uint(1) {
		t.Errorf("BytesLeft() must return 1")
	}
	if reader.BitsLeft() != uint(4) {
		t.Errorf("BitsLeft() must return 4")
	}
}

func TestReadBitsToByteArray(t *testing.T) {

	// Input: 01010101 01010101
	byteArray := []byte{0x55, 0x55}
	reader := NewBitReader(byteArray)

	reader.Reset()
	buf, _ := reader.ReadBitsToByteArray(3)
	if reflect.DeepEqual([]byte{0x02}, buf) == false {
		t.Errorf("buf %#v must equal {0x02}", buf)
	}

	reader.Reset()
	buf, _ = reader.ReadBitsToByteArray(4)
	if reflect.DeepEqual([]byte{0x05}, buf) == false {
		t.Errorf("buf %#v must equal {0x05}", buf)
	}

	reader.Reset()
	buf, _ = reader.ReadBitsToByteArray(8)
	if reflect.DeepEqual([]byte{0x55}, buf) == false {
		t.Errorf("buf %#v must equal {0x55}", buf)
	}

	reader.Reset()
	buf, _ = reader.ReadBitsToByteArray(12)
	if reflect.DeepEqual([]byte{0x05, 0x55}, buf) == false {
		t.Errorf("buf %#v must equal {0x05, 0x55}", buf)
	}

	reader.Reset()
	buf, _ = reader.ReadBitsToByteArray(16)
	if reflect.DeepEqual([]byte{0x55, 0x55}, buf) == false {
		t.Errorf("buf %#v must equal {0x55, 0x55}", buf)
	}
}

func TestReadBytes(t *testing.T) {

	// Input: 0x00 0x01 0x02 0x03
	byteArray := []byte{0, 1, 2, 3}
	reader := NewBitReader(byteArray)

	// Test 1: Read 2 bytes (0x00 0x01)
	r1, _ := reader.ReadBytes(2)
	if len(r1) != 2 {
		t.Errorf("Length of r1 (%d) must be 2", len(r1))
	}
	if reflect.DeepEqual([]byte{0, 1}, r1) == false {
		t.Errorf("r1 %#v must equal {0x0, 0x1}", r1)
	}

	// Test 2: Read a bit from 0x02 (subsequent ReadBytes still return entire byte)
	r2, _ := reader.ReadBits(1)
	if r2 != byte(0) {
		t.Errorf("r2 (%d) must be euqal to 0", r2)
	}

	// Test 3: Read remaining 2 bytes (0x02 0x03) despite of bit being read from test 2
	r3, _ := reader.ReadBytes(2)
	if len(r3) != 2 {
		t.Errorf("Length of r3 (%d) must be 2", len(r3))
	}
	if reflect.DeepEqual([]byte{0x2, 0x3}, r3) == false {
		t.Errorf("r3 %#v must equal {0x2, 0x3}", r3)
	}
}

func TestSkipBitsBytes(t *testing.T) {

	// Input: 0101 0101 0101 0101 (0x55 0x55)
	byteArray := []byte{0x55, 0x55}
	reader := NewBitReader(byteArray)

	// Test 1: Skip 3 bits (010), verify next bit is 1
	reader.SkipBits(3) // skip 010
	r1, _ := reader.ReadBits(1)
	if r1 != byte(1) {
		t.Errorf("r1 (%d) must equal 1", r1)
	}

	// Test 2: Skip 1 byte, verify posInCurrentByte = 7
	reader.SkipBytes(1)
	if reader.posInCurrentByte != uint(7) {
		t.Errorf("Position in byte (%d) must equal 7", reader.posInCurrentByte)
	}

	// Test 3: Read 8 bytes (should be 0x55)
	r3, _ := reader.ReadBits(8)
	if r3 != byte(0x55) {
		t.Errorf("r3 (%d) must equal 0x55 (%d)", r3, 0x55)
	}
}

func TestSkipBytesToEnd(t *testing.T) {
	byteArray := []byte{0x55, 0x55}
	reader := NewBitReader(byteArray)
	reader.SkipBytes(reader.BytesLeft())
	if reader.HasByteLeft() {
		t.Errorf("HasByteLeft() must return false after skipping all remaining bytes.")
	}
}

func TestBitsBytesLeft(t *testing.T) {

	// Input: 0101 0101 0101 0101 (0x55 0x55)
	byteArray := []byte{0x55, 0x55}
	reader := NewBitReader(byteArray)

	// Test 1: Skip 2; bits left = 14
	reader.SkipBits(2)
	if reader.BitsLeft() != uint(14) {
		t.Errorf("BitsLeft() (%d) must return 14", reader.BitsLeft())
	}
	if reader.HasBitLeft() != true {
		t.Errorf("HasBitLeft() must return true")
	}
	if reader.HasByteLeft() != true {
		t.Errorf("HasByteLeft() must return true")
	}

	// Test 2: Read 4; bits left = 10
	reader.ReadBits(4)
	if reader.BitsLeft() != uint(10) {
		t.Errorf("BitsLeft() (%d) must return 10", reader.BitsLeft())
	}
	if reader.HasBitLeft() != true {
		t.Errorf("HasBitLeft() must return true")
	}
	if reader.HasByteLeft() != true {
		t.Errorf("HasByteLeft() must return true")
	}

	// Test 3: Read 1 byte; bits left = 8
	reader.ReadBytes(1)
	if reader.BitsLeft() != uint(8) {
		t.Errorf("BitsLeft() (%d) must return 8", reader.BitsLeft())
	}
	if reader.HasBitLeft() != true {
		t.Errorf("HasBitLeft() must return true")
	}
	if reader.HasByteLeft() != true {
		t.Errorf("HasByteLeft() must return true")
	}

	// Test 4: Read all remaining 8 bits left
	reader.ReadBits(8)
	if reader.HasBitLeft() != false {
		t.Errorf("HasBitLeft() must return true")
	}
	if reader.HasByteLeft() != false {
		t.Errorf("HasByteLeft() must return true")
	}
}

func TestReadBitsAsUInt(t *testing.T) {

	// Input: 0101 0101 0101 0101 0101 0101 (0x55 0x55, 0x55)
	byteArray := []byte{0x55, 0x55, 0x55}
	reader := NewBitReader(byteArray)
	val, _ := reader.ReadBitsAsUInt(5)
	if val != uint(10) {
		t.Errorf("val (%d) must equal 10", val)
	}
	val, _ = reader.ReadBitsAsUInt(3)
	if val != uint(5) {
		t.Errorf("val (%d) must equal 5", val)
	}
	reader.SkipBits(1)
	val, _ = reader.ReadBitsAsUInt(6)
	if val != uint(42) {
		t.Errorf("val (%d) must equal 42", val)
	}
}

func TestReadBitAsBool(t *testing.T) {

	reader := NewBitReader([]byte{0x55})
	val, _ := reader.ReadBitAsBool()
	if val != false {
		t.Errorf("val must be false")
	}
	val, _ = reader.ReadBitAsBool()
	if val != true {
		t.Errorf("val must be true")
	}
	val, _ = reader.ReadBitAsBool()
	if val != false {
		t.Errorf("val must be false")
	}
	val, _ = reader.ReadBitAsBool()
	if val != true {
		t.Errorf("val must be true")
	}
	val, _ = reader.ReadBitAsBool()
	if val != false {
		t.Errorf("val must be false")
	}
	val, _ = reader.ReadBitAsBool()
	if val != true {
		t.Errorf("val must be true")
	}
	val, _ = reader.ReadBitAsBool()
	if val != false {
		t.Errorf("val must be false")
	}
	val, _ = reader.ReadBitAsBool()
	if val != true {
		t.Errorf("val must be true")
	}
}

func TestByteAlignment(t *testing.T) {
	// 00000000 11111111
	reader := NewBitReader([]byte{0x00, 0xff})

	// Test 1: ByteAlign on bit7 should have no change
	reader.ByteAlign()
	if reader.posInCurrentByte != uint(7) {
		t.Errorf("Position in current byte (%d) must be 7", reader.posInCurrentByte)
	}
	if reader.currentByteIndex != uint(0) {
		t.Errorf("Current Byte Index (%d) must be 0", reader.currentByteIndex)
	}

	// Test 2: ByteAlign after skipping 1 bits should align to next byte
	reader.Reset()
	reader.SkipBits(1)
	reader.ByteAlign()
	if reader.posInCurrentByte != uint(7) {
		t.Errorf("Position in current byte (%d) must be 7", reader.posInCurrentByte)
	}
	if reader.currentByteIndex != uint(1) {
		t.Errorf("Current Byte Index (%d) must be 1", reader.currentByteIndex)
	}
}

func TestErrors(t *testing.T) {

	// 00000000 11111111
	reader := NewBitReader([]byte{0x00, 0xff})

	// Test 1: Read more bytes
	_, err := reader.ReadBytes(3)
	if err == nil {
		t.Errorf("err cannot be nil")
	}

	// Test 2: Read more bits
	reader.SkipBits(15)
	_, err = reader.ReadBits(3)
	if err == nil {
		t.Errorf("err cannot be nil")
	}

	// Test 3: Read more bits
	reader.ReadBits(8)
	reader.ReadBits(8)
	_, err = reader.ReadBits(3)
	if err == nil {
		t.Errorf("err cannot be nil")
	}
}
