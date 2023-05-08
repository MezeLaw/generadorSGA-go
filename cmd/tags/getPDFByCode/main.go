package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	pdfFilename := request.PathParameters["pdf"]
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	// nombre del bucket de S3 y ruta al archivo PDF
	bucketName := "pdfs-genrador-sga"
	objectKey := pdfFilename + ".pdf"

	println("El key es>>>", objectKey)

	// descarga el archivo PDF del bucket de S3
	resp, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &objectKey,
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read object: %v", err)
	}

	//headers := map[string]string{
	//	"Content-Type":        "application/pdf",
	//	"Content-Disposition": fmt.Sprintf("attachment; filename=%s", objectKey),
	//}

	response := events.APIGatewayProxyResponse{
		Body:       string(body),
		Headers:    make(map[string]string),
		StatusCode: http.StatusOK,
	}

	response.Headers["Content-Type"] = "application/pdf"
	response.Headers["Content-Disposition"] = "attachment; filename=archivo.pdf"

	return &response, nil
}
