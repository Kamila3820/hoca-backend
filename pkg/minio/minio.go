package minio

import (
	"log"

	mod "github.com/Kamila3820/hoca-backend/pkg"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func Init() {
	endpoint := "minio.bsthun.com"
	accessKeyID := "oP2LF7ioFyN9j27uTFpN"
	secretAccessKey := "ybBIrJnrxCRs3xV37zaC2nElB34wvUJk9PKTm4zr"
	useSSL := true

	//Initialize minio client object
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
	}

	mod.Minio = minioClient
}
