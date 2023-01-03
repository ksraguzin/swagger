// This file is safe to edit. Once it exists it will not be overwritten

package handler

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"service-pdf-compose/internal/handler/operations"
	"service-pdf-compose/pkg/composer"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

//go:generate swagger generate server --target ..\..\..\server --name ServicePdfCompose --spec ..\..\api\api.yml --server-package internal/handler --principal interface{}

func configureFlags(api *operations.ServicePdfComposeAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ServicePdfComposeAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.MultipartformConsumer = runtime.DiscardConsumer

	api.BinProducer = runtime.ByteStreamProducer()

	// You may change here the memory limit for this multipart form parser. Below is the default (32 MB).
	// operations.PostSendMaxParseMemory = 32 << 20
	// if api.PostSendHandler == nil {
	api.PostSendHandler = operations.PostSendHandlerFunc(func(params operations.PostSendParams) middleware.Responder {
		// return middleware.NotImplemented("operation operations.PostSend has not yet been implemented")
		if params.Upfile1 == nil || params.Upfile2 == nil || params.Upfile3 == nil {
			return middleware.Error(404, fmt.Errorf("no file provided"))
		}
		defer func() {
			_ = params.Upfile1.Close()
			_ = params.Upfile2.Close()
			_ = params.Upfile3.Close()
		}()
		fs := []io.ReadCloser{params.Upfile1, params.Upfile2, params.Upfile3}
		files, err := composer.ComposeFromFiles(fs)
		if err != nil {
			return operations.NewPostSendInternalServerError()
			return middleware.Error(500, fmt.Errorf("could not compose PDF-file on server"))
		}
		// f, err := os.Open("D:/Яндекс/дебаг/raguzinkc.pdf")
		// if err != nil {
		// 	log.Println("error open file")
		// 	return middleware.Error(404, fmt.Errorf("no file provided"))
		// }
		// defer f.Close()
		// log.Println("error handler")
		// fileEx := ioutil.NopCloser(f)
		// return operations.NewPostSendOK().WithPayload(fileEx)
		return operations.NewPostSendOK().WithPayload(files)

	})
	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
