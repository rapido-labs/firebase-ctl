package firebase

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"path/filepath"
	"testing"
)

type FsTestSuite struct {
	suite.Suite
	fs afero.Fs
}

func (c *FsTestSuite) SetupTest() {
	c.fs = afero.NewMemMapFs()

	c.fs.Mkdir(correctJsonsDir, 0644)
	c.fs.MkdirAll(filepath.Join(wrongJsonDir, "test"), 0644)
	f, _ := c.fs.OpenFile(filepath.Join(correctJsonsDir, "1.json"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	defer f.Close()
	f.WriteString("{}")
	g, _ := c.fs.OpenFile(filepath.Join(wrongJsonDir, "1.json"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	defer g.Close()
	g.WriteString("{")
}

const correctJsonsDir = "correct-json"
const wrongJsonDir = "wrong-json"
const testDir = ""

func (c *FsTestSuite) TestReadJsonFromFile() {
	customFs := customFs{fs: c.fs}

	c.T().Run("Should unmarshal successfully", func(t *testing.T) {
		sampleMap := make(map[string]interface{})
		err := customFs.UnmarshalFromFile(filepath.Join(correctJsonsDir, "1.json"), &sampleMap)
		assert.NoError(t, err)
	})
	c.T().Run("Should exit because invalid file", func(t *testing.T) {
		sampleMap := make(map[string]interface{})
		err := customFs.UnmarshalFromFile(filepath.Join(correctJsonsDir, "2.json"), &sampleMap)
		assert.Contains(t, err.Error(), "does not exist")
	})

}

func (c *FsTestSuite) TestWriteJsonToFile() {

	customFs := customFs{fs: c.fs}

	c.T().Run("Should marshal successfully", func(t *testing.T) {
		err := customFs.WriteJsonToFile("", "some-path")
		assert.NoError(t, err)
	})
	c.T().Run("Should throw error no file permission", func(t *testing.T) {
		customFs.fs = afero.NewReadOnlyFs(customFs.fs)
		err := customFs.WriteJsonToFile("", "some-path")
		assert.Error(t, err)
	})
}

func (c *FsTestSuite) TestReadDirAndUnMarshal() {
	customFs := customFs{fs: c.fs}
	c.T().Run("test should scan the directory with only jsons and return all values", func(t *testing.T) {

		err := customFs.UnMarshalFromDir(correctJsonsDir, &map[string]interface{}{})
		assert.NoError(t, err)
	})
	c.T().Run("test should return error when there is invalid json", func(t *testing.T) {
		err := customFs.UnMarshalFromDir(wrongJsonDir, &map[string]interface{}{})
		assert.Error(t, err)
	})
	c.T().Run("test should return error when a non-existent directory is used", func(t *testing.T) {
		err := customFs.UnMarshalFromDir("some-random-dir", &map[string]interface{}{})
		assert.Contains(t, err.Error(), "does not exist")
	})
	c.T().Run("test should error out when a filepath is passed instead of a directory path", func(t *testing.T) {
		err := customFs.UnMarshalFromDir(filepath.Join(correctJsonsDir, "1.json"), &map[string]interface{}{})
		assert.Error(t, err)
	})
}

func TestFs(t *testing.T) {
	suite.Run(t, new(FsTestSuite))
}
