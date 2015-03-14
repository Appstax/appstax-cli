package hosting

import (
	"appstax-cli/appstax/apiclient"
	"appstax-cli/appstax/fail"
	"archive/tar"
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func UploadStatic(archivePath string, progressWriter io.Writer) error {
	_, _, err := apiclient.PostFile(apiclient.Url("/appstax/hosting/static"), archivePath, progressWriter)
	return err
}

func PrepareArchive(rootPath string) (string, int64) {
	file, err := ioutil.TempFile("", "")
	fail.Handle(err)
	defer file.Close()
	fileWriter := bufio.NewWriter(file)
	defer fileWriter.Flush()
	gzipWriter, err := gzip.NewWriterLevel(fileWriter, gzip.BestCompression)
	fail.Handle(err)
	defer gzipWriter.Close()
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	fullRootPath, err := filepath.Abs(rootPath)
	fail.Handle(err)
	addAllToArchive(fullRootPath, tarWriter)

	tarWriter.Close()
	gzipWriter.Close()
	fileWriter.Flush()
	file.Close()

	fileInfo, err := os.Stat(file.Name())
	fail.Handle(err)
	return file.Name(), fileInfo.Size()
}

func addAllToArchive(fullRootPath string, tarWriter *tar.Writer) {
	filepath.Walk(fullRootPath, func(path string, fileInfo os.FileInfo, err error) error {
		fail.Handle(err)
		if !fileInfo.IsDir() && fileInfo.Name()[:1] != "." {
			addFileToArchive(path, path[len(fullRootPath+"/"):], tarWriter, fileInfo)
		}
		return nil
	})
}

func addFileToArchive(filePath string, addPath string, tarWriter *tar.Writer, fileInfo os.FileInfo) {
	//println("Adding " + addPath + " (" + filePath + ")")
	fileReader, err := os.Open(filePath)
	fail.Handle(err)
	defer fileReader.Close()

	header := new(tar.Header)
	header.Name = addPath
	header.Size = fileInfo.Size()
	header.Mode = int64(fileInfo.Mode())
	header.ModTime = fileInfo.ModTime()

	err = tarWriter.WriteHeader(header)
	fail.Handle(err)
	_, err = io.Copy(tarWriter, fileReader)
	fail.Handle(err)
}
