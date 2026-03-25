package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword_Valid(t *testing.T) {
	err := ValidatePassword("P@ssw0rd")
	assert.NoError(t, err)
}

func TestValidatePassword_TooShort(t *testing.T) {
	err := ValidatePassword("P@ss1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 8 characters")
}

func TestValidatePassword_NoUppercase(t *testing.T) {
	err := ValidatePassword("p@ssw0rd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Maj")
}

func TestValidatePassword_NoLowercase(t *testing.T) {
	err := ValidatePassword("P@SSW0RD")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "lower case")
}

func TestValidatePassword_NoNumber(t *testing.T) {
	err := ValidatePassword("P@ssword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Number")
}

func TestValidatePassword_NoSpecial(t *testing.T) {
	err := ValidatePassword("Passw0rd")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "special")
}

func TestValidatePassword_AllSpecialChars(t *testing.T) {
	specials := []string{"!", "@", "#", "$", "%", "^", "&", "*", "."}
	for _, s := range specials {
		err := ValidatePassword("Passw0r" + s)
		assert.NoError(t, err, "should accept special char: %s", s)
	}
}

func TestGenerateTempPassword_Length(t *testing.T) {
	pw, err := GenerateTempPassword(12)
	assert.NoError(t, err)
	assert.Len(t, pw, 12)
}

func TestGenerateTempPassword_MeetsValidation(t *testing.T) {
	// Run multiple times to account for randomness
	for i := 0; i < 20; i++ {
		pw, err := GenerateTempPassword(12)
		assert.NoError(t, err)
		assert.NoError(t, ValidatePassword(pw), "generated password should pass validation: %s", pw)
	}
}

func TestGenerateTempPassword_MinLengthFloor(t *testing.T) {
	pw, err := GenerateTempPassword(2)
	assert.NoError(t, err)
	assert.Len(t, pw, 8, "length below 4 should be raised to 8")
}

func TestGenerateTempPassword_Randomness(t *testing.T) {
	pw1, _ := GenerateTempPassword(12)
	pw2, _ := GenerateTempPassword(12)
	assert.NotEqual(t, pw1, pw2, "two generated passwords should differ")
}
