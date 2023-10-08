package sixtwo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pevans/erc/memory"
	"github.com/stretchr/testify/suite"
)

type sixtwoSuite struct {
	suite.Suite

	physDisk   *memory.Segment
	physTrack  *memory.Segment
	physSector *memory.Segment
	logDisk    *memory.Segment
	logTrack   *memory.Segment
	logSector  *memory.Segment
	imageType  int

	baseDir string
}

func (s *sixtwoSuite) SetupSuite() {
	dir, err := os.Getwd()
	s.NoError(err)

	s.baseDir = dir + "/../../data"

	s.imageType = DOS33
	s.physDisk = memory.NewSegment(NibSize)
	s.physTrack = memory.NewSegment(PhysTrackLen)
	s.physSector = memory.NewSegment(PhysSectorLen)
	s.logDisk = memory.NewSegment(DosSize)
	s.logTrack = memory.NewSegment(LogTrackLen)
	s.logSector = memory.NewSegment(LogSectorLen)

	s.NoError(loadFile(s.physDisk, s.baseDir+"/physical.disk"))
	s.NoError(loadFile(s.physTrack, s.baseDir+"/physical.track"))
	s.NoError(loadFile(s.physSector, s.baseDir+"/physical.sector"))
	s.NoError(loadFile(s.logDisk, s.baseDir+"/logical.disk"))
	s.NoError(loadFile(s.logTrack, s.baseDir+"/logical.track"))
	s.NoError(loadFile(s.logSector, s.baseDir+"/logical.sector"))
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(sixtwoSuite))
}

func loadFile(seg *memory.Segment, path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for i, b := range bytes {
		seg.Mem[i] = uint8(b)
	}

	return nil
}
