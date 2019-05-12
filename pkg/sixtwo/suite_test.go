package sixtwo

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type sixtwoSuite struct {
	suite.Suite

	physDisk   *data.Segment
	logSector  *data.Segment
	physSector *data.Segment
	imageType  int

	baseDir       string
	physDiskPath  string
	logSectorPath string
}

func (s *sixtwoSuite) SetupSuite() {
	dir, err := os.Getwd()
	assert.NoError(s.T(), err)

	s.baseDir = dir + "/../../data"
	s.physDiskPath = s.baseDir + "/physical.disk"
	s.logSectorPath = s.baseDir + "/logical.sector"

	s.imageType = DOS33
	s.physDisk = data.NewSegment(NibSize)
	s.logSector = data.NewSegment(LogSectorLen)
	s.physSector = data.NewSegment(PhysSectorLen)

	assert.NoError(s.T(), loadFile(s.physDisk, s.baseDir+"/physical.disk"))
	assert.NoError(s.T(), loadFile(s.logSector, s.baseDir+"/logical.sector"))
	assert.NoError(s.T(), loadFile(s.physSector, s.baseDir+"/physical.sector"))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(sixtwoSuite))
}

func loadFile(seg *data.Segment, path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for i, b := range bytes {
		seg.Mem[i] = data.Byte(b)
	}

	return nil
}

func fileMatches(seg *data.Segment, path string) bool {
	fseg := data.NewSegment(len(seg.Mem))
	err := loadFile(fseg, path)
	if err != nil {
		log.Printf("Couldn't complete fileMatches: file doesn't exist: %v\n", path)
		return false
	}

	if len(seg.Mem) != len(fseg.Mem) {
		log.Printf("Seg mem size (%v) mismatches fseg mem size (%v)\n", len(seg.Mem), len(fseg.Mem))
		return false
	}

	for i, b := range seg.Mem {
		if b != fseg.Mem[i] {
			log.Printf("seg byte %x (%x) mismatches fseg byte (%x)\n", i, b, fseg.Mem[i])
			return false
		}
	}

	return true
}
