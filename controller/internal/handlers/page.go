package handlers

import (
	"bufio"
	"context"
	"controller-service/client/operations"
	"fmt"
	"html/template"
	"net/http"

	"time"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	apiclient "github.com/go-swagger/go-swagger/examples/cli/client"
)

type Handler struct {
	HttpClient *http.Client
}

func (h *Handler) Send(w http.ResponseWriter, r *http.Request) {

	transport := httptransport.New(apiclient.DefaultHost+":8090", apiclient.DefaultBasePath, []string{"http"})
	transport.Consumers["application/pdf"] = runtime.ByteStreamConsumer()

	c := operations.New(transport, strfmt.Default)

	f1, _, _ := r.FormFile("upfile1")
	f2, _, _ := r.FormFile("upfile2")
	f3, _, _ := r.FormFile("upfile3")

	upfile1 := runtime.NamedReader("upfile1", f1)
	upfile2 := runtime.NamedReader("upfile2", f2)
	upfile3 := runtime.NamedReader("upfile3", f3)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	params := &operations.PostSendParams{
		Upfile1:    upfile1,
		Upfile2:    upfile2,
		Upfile3:    upfile3,
		Context:    ctx,
		HTTPClient: h.HttpClient,
	}

	Ok, err := c.PostSend(params, w)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	if Ok == nil {
		w.Write([]byte("error PostSend nil"))
		return
	}
	if Ok.Payload != nil {
		w.Write([]byte(Ok.Error()))
		return
	}
	w.Header().Set("Content-Type", "pdf-compose-service")
	w.WriteHeader(http.StatusOK)
	buf := bufio.NewWriter(Ok.GetPayload())
	w.Write(buf.AvailableBuffer())
	// w.Write([]byte("Ok"))
}

func (h *Handler) Web(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("web/templates/form.html")
	if err != nil {
		fmt.Println(err)
	}

	t.Execute(w, nil)
}
