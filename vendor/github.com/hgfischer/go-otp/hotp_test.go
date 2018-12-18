package otp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHOTP(t *testing.T) {
	secret := `12345678901234567890`
	table := map[uint64]string{
		0: `755224`,
		1: `287082`,
		2: `359152`,
		3: `969429`,
		4: `338314`,
		5: `254676`,
		6: `287922`,
		7: `162583`,
		8: `399871`,
		9: `520489`,
	}

	for cnt, expected := range table {
		hotp := &HOTP{Secret: secret, Counter: cnt, Length: 6}
		result := hotp.Get()
		assert.Equal(t, expected, result)
	}
}

func TestHOTPBase32(t *testing.T) {
	secret := `JBSWY3DPEHPK3PXP`
	table := map[uint64]string{
		0: `282760`,
		1: `996554`,
		2: `602287`,
		3: `143627`,
		4: `960129`,
		5: `768897`,
		6: `883951`,
		7: `449891`,
		8: `964230`,
		9: `924769`,
	}

	for cnt, expected := range table {
		hotp := &HOTP{Secret: secret, Counter: cnt, Length: 6, IsBase32Secret: true}
		result := hotp.Get()
		assert.Equal(t, expected, result)
	}
}

func TestHOTPShouldBeCroppedToMaxLength(t *testing.T) {
	hotp := &HOTP{Length: 20}
	result := hotp.Get()
	assert.Equal(t, MaxLength, len(result))
}

func TestHOTPShouldUseDefaultValues(t *testing.T) {
	hotp := &HOTP{}
	result := hotp.Get()
	assert.Equal(t, uint8(DefaultLength), hotp.Length)
	assert.NotEmpty(t, hotp.Secret)
	assert.Equal(t, hotp.Length, uint8(len(result)))
}
