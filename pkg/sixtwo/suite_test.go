package sixtwo

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/pevans/erc/pkg/data"
	"github.com/stretchr/testify/suite"
)

type sixtwoSuite struct {
	suite.Suite

	physDisk   *data.Segment
	physTrack  *data.Segment
	physSector *data.Segment
	logDisk    *data.Segment
	logTrack   *data.Segment
	logSector  *data.Segment
	imageType  int

	baseDir string
}

func (s *sixtwoSuite) SetupSuite() {
	dir, err := os.Getwd()
	s.NoError(err)

	s.baseDir = dir + "/../../data"

	s.imageType = DOS33
	s.physDisk = data.NewSegment(NibSize)
	s.physTrack = data.NewSegment(PhysTrackLen)
	s.physSector = data.NewSegment(PhysSectorLen)
	s.logDisk = data.NewSegment(DosSize)
	s.logTrack = data.NewSegment(LogTrackLen)
	s.logSector = data.NewSegment(LogSectorLen)

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

func loadFile(seg *data.Segment, path string) error {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for i, b := range bytes {
		seg.Mem[i] = uint8(b)
	}

	return nil
}
