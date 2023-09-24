package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_String(t *testing.T) {
	tests := []struct {
		name    string
		builder Builder
		want    string
	}{
		{
			builder: Builder{
				After:  "2023/07/31",
				Before: "2023/08/31",
			},
			want: "after:2023/07/31 before:2023/08/31",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.builder.String())
		})
	}
}
