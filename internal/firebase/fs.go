package firebase

import (
	"encoding/json"
	"errors"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type customFs struct{
	fs afero.Fs

}

func (f *customFs) UnMarshalFromDir(dirName string, data interface{}) error{
	dir, err := f.fs.Open(dirName)
	if err != nil{
		return err
	}
	defer dir.Close()
	fileInfoList, err := dir.Readdir(0)
	if err!= nil{
		return err
	}
	var errs []error
	for _, fileInfo := range fileInfoList{
		if fileInfo.IsDir(){
			continue
		}
		err = f.UnmarshalFromFile(filepath.Join(dirName,fileInfo.Name()), data)
		if err!= nil{
			errs = append(errs, err)
		}
	}
	if len(errs) != 0{
		sb := strings.Builder{}
		for i := range errs{
			sb.WriteString("\n\t" + errs[i].Error())
		}
		return errors.New(sb.String())
	}
	return nil

}
func (f *customFs) UnmarshalFromFile(fileName string, data interface{}) error{
	file, err := f.fs.Open(fileName)
	if err!= nil{
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(data)
}
func (f *customFs)WriteJsonToFile(data interface{}, filePath string)error{
	f.fs.MkdirAll(filepath.Dir(filePath), 0744)
	file, err := f.fs.OpenFile(filePath, os.O_RDWR| os.O_CREATE| os.O_TRUNC,0644 )
	if err!= nil{
		return err
	}
	defer file.Close()
	d := json.NewEncoder(file)
	d.SetIndent("", "\t")
	d.SetEscapeHTML(false)
	return d.Encode(data)
}