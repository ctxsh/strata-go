package apex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildProm(t *testing.T) {
	expected := `# HELP apex_test a test metric
# TYPE apex_test counter
apex_test{region="us-east-1"} 1E+00
`
	actual := BuildProm("apex_test", "a test metric", "counter", map[string]string{"region": "us-east-1"}, 1.0)
	assert.Equal(t, expected, actual)
}
