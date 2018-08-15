package sarsar

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParseSarFile(t *testing.T) {
	f, err := parseSarFile("../test/sa14.out")
	assert.NoError(t, err)
	assert.Equal(t, 18 /* sections */, len(f.sections))
	cpuUtil := f.sections[SECTION_CPU_UTIL]
	assert.Equal(t, 8060 /* rows */, len(cpuUtil.records))
	record := cpuUtil.records[0]
	assert.Equal(t, 11 /* columns */, len(record.data))
}
