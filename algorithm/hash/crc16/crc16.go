package crc16

/*

	Name  : CRC-16/IBM
	Poly  : 0x8005    x^16 + x^15 + x^2 + 1
	Polyrv: 0xA001
	Init  : 0xFFFF
	Revert: true
	XorOut: 0x0000
	Check : 0x4B37 ("123456789")
	MaxLen: 4095 bytes

*/

const (
	POLY_IBM uint16 = 0xA001
)

type Table [256]uint16

var TableIBM = MakeTable(POLY_IBM)

func MakeTable(polynom uint16) *Table {
	table := new(Table)
	for i := 0; i < 256; i++ {
		crc := uint16(i)
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ polynom
			} else {
				crc >>= 1
			}
		}
		table[i] = crc
	}
	return table
}

func Update(crc uint16, table *Table, data []byte) uint16 {
	for _, b := range data {
		crc = table[byte(crc)^b] ^ (crc >> 8)
	}
	return crc
}

func ChecksumIBM(data []byte) uint16 {

	return Update(0xFFFF, TableIBM, data)
}
