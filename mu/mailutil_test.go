package mu

import (
	"regexp"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBase64(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBase64(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBase64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeBase64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_regexp(t *testing.T) {
	assert.True(t, unicode.IsSpace(' '))
	assert.False(t, regexp.MustCompile(`\s`).MatchString(" "))
}

func TestRemoveHTMLTag(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RemoveHTMLTag(tt.body))
		})
	}
}

func TestNormalizeSpaces(t *testing.T) {
	assert.Equal(t, "a b", NormalizeSpaces("a  b"))
}
