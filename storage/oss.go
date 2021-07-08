package storage

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"log"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
)

var client *oss.Client

func InitOssClient() {

	er := godotenv.Load()
	if er != nil {
		log.Fatalf("加载oss参数失败:%s", er.Error())
	}

	endpoint := os.Getenv("endpoint")
	accessKeyId := os.Getenv("accessKeyId")
	accessKeySecret := os.Getenv("accessKeySecret")

	var err error
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("连接OSS失败:%s", err.Error())
	}
}

func Test018() *oss.Bucket {
	buc, err := client.Bucket("test018")
	if err != nil {
		log.Fatalf(err.Error())
	}

	return buc
}

func handleError(err error) {
	fmt.Println("Error:", err)
	os.Exit(-1)
}

type myReader struct {
	rder io.Reader

	hasher hash.Hash
	size   int64
}

func newMyReader(rder io.Reader) *myReader {
	return &myReader{
		rder:   rder,
		hasher: sha256.New(),
		size:   0,
	}
}

func (m *myReader) Read(p []byte) (n int, err error) {
	n, err = m.rder.Read(p)
	if err != nil {
		return
	}

	m.size += int64(n)
	m.hasher.Write(p)
	return
}

func (m *myReader) getSize() int64 {
	return m.size
}

func (m *myReader) getHash() string {
	sum := m.hasher.Sum(nil)
	return fmt.Sprintf("%x", sum)
}

func Put(objectName string, file io.Reader) error {
	return Test018().PutObject(objectName, file)
}

func Get(objName string) (io.ReadCloser, error) {
	return Test018().GetObject(objName)
}

func put(objectName string, localFileName string) error {
	bucket, err := client.Bucket("test018")
	if err != nil {
		handleError(err)
	}

	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return err
	}
	return nil
}

func dwn(objectName string, downloadedFileName string) error {
	return Test018().GetObjectToFile(objectName, downloadedFileName)
}

func list() error {
	bucket, err := client.Bucket("test018")
	if err != nil {
		handleError(err)
	}
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}
		for _, object := range lsRes.Objects {
			fmt.Println("BUcket:", object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return nil
}

func delete(objectName string) error {
	bucket, err := client.Bucket("test018")
	if err != nil {
		handleError(err)
	}

	err = bucket.DeleteObject(objectName)
	if err != nil {
		return err
	}
	return nil
}
