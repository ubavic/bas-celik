package card

import "fmt"

// Constructs an APDU (Application Protocol Data Unit) command according
// to the specifications from the ISO 7816-4 (5. Organization for interchange).
func buildAPDU(cla, ins, p1, p2 byte, data []byte, ne uint) []byte {
	length := len(data)

	if length > 0xFFFF {
		panic(fmt.Errorf("APDU command length too large"))
	}

	apdu := make([]byte, 4)
	apdu[0] = cla
	apdu[1] = ins
	apdu[2] = p1
	apdu[3] = p2

	if length == 0 {
		if ne != 0 {
			if ne <= 256 {
				l := byte(0x00)
				if ne != 256 {
					l = byte(ne)
				}
				apdu = append(apdu, l)
			} else {
				var l1, l2 byte
				if ne == 65536 {
					l1 = 0
					l2 = 0
				} else {
					l1 = byte(ne >> 8)
					l2 = byte(ne)
				}
				apdu = append(apdu, []byte{l1, l2}...)
			}
		}
	} else {
		if ne == 0 {
			if length <= 255 {
				apdu = append(apdu, byte(length))
				apdu = append(apdu, data...)
			} else {
				l := []byte{0x0, byte(length >> 8), byte(length)}
				apdu = append(apdu, l...)
				apdu = append(apdu, data...)
			}
		} else {
			if length <= 255 && ne <= 256 {
				apdu = append(apdu, byte(length))
				apdu = append(apdu, data...)
				if ne != 256 {
					apdu = append(apdu, byte(ne))
				} else {
					apdu = append(apdu, 0x00)
				}
			} else {
				l := []byte{0x00, byte(length >> 8), byte(length)}
				apdu = append(apdu, l...)
				apdu = append(apdu, data...)
				if ne != 65536 {
					neB := []byte{byte(ne >> 8), byte(ne)}
					apdu = append(apdu, neB...)
				}
			}
		}
	}

	return apdu
}
