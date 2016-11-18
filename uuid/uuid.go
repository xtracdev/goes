package uuid

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

//GenerateUuidV4 will generate a uuid that is a valid v4 uuid (https://tools.ietf.org/html/rfc4122)
/*The version 4 UUID is meant for generating UUIDs from truly-random or
pseudo-random numbers.

	The algorithm is as follows:

o  Set the two most significant bits (bits 6 and 7) of the
clock_seq_hi_and_reserved to zero and one, respectively.

	o  Set the four most significant bits (bits 12 through 15) of the
time_hi_and_version field to the 4-bit version number from
Section 4.1.3.

	o  Set all the other bits to randomly (or pseudo-randomly) chosen
values.
*/
func GenerateUuidV4() (string, error) {
	//16 random bytes
	bytes, err := GenerateRandomBytes(16)
	if err != nil {
		return "", err
	}
	//set v4 byte
	bytes[6] = (bytes[6] & 0xf) | 0x4<<4
	//set version rfc4122
	bytes[8] = (bytes[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:]), nil
}
