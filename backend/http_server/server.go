package http_server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/golang/protobuf/jsonpb"

	"github.com/gorilla/mux"
	todomvc "github.com/oinume/todomvc-example/proto-gen/go/proto/todomvc"
)

type server struct {
}

func New() *server {
	return &server{}
}

func (s *server) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/todos", s.CreateTodo).Methods("POST")
	//r.HandleFunc("/todos", s.fetcher).Methods("GET")
	//r.HandleFunc("/todos/{id}", s.fetcher).Methods("Put")
	return r
}

func (s *server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	um := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	req := &todomvc.CreateTodoItemRequest{}
	if err := um.Unmarshal(r.Body, req); err != nil {
		internalServerError(w, err)
		return
	}

	_, _ = fmt.Fprintf(w, "req = %+v\n", req)
	_, _ = fmt.Fprintf(w, "id = %+v\n", req.TodoItem.Id)
}

func internalServerError(w http.ResponseWriter, err error) {
	//switch _ := errors.Cause(err).(type) { // TODO:
	//default:
	// unknown error
	//sUserID := ""
	//if userID == 0 {
	//	sUserID = fmt.Sprint(userID)
	//}
	//util.SendErrorToRollbar(err, sUserID)
	//fields := []zapcore.Field{
	//	zap.Error(err),
	//}
	//if e, ok := err.(errors.StackTracer); ok {
	//	b := &bytes.Buffer{}
	//	for _, f := range e.StackTrace() {
	//		fmt.Fprintf(b, "%+v\n", f)
	//	}
	//	fields = append(fields, zap.String("stacktrace", b.String()))
	//}
	//if appLogger != nil {
	//	appLogger.Error("internalServerError", fields...)
	//}

	http.Error(w, fmt.Sprintf("Internal Server Error\n\n%v", err), http.StatusInternalServerError)
	//if !config.IsProductionEnv() {
	//	fmt.Fprintf(w, "----- stacktrace -----\n")
	//	if e, ok := err.(errors.StackTracer); ok {
	//		for _, f := range e.StackTrace() {
	//			fmt.Fprintf(w, "%+v\n", f)
	//		}
	//	}
	//}
}

func writeJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		http.Error(w, `{ "status": "Failed to Encode as writeJSON" }`, http.StatusInternalServerError)
	}
}

//func writeHTML(w http.ResponseWriter, code int, body string) {
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	w.WriteHeader(code)
//	if _, err := fmt.Fprint(w, body); err != nil {
//		http.Error(w, "Failed to write HTML", http.StatusInternalServerError)
//	}
//}

func writeHTMLWithTemplate(
	w http.ResponseWriter,
	code int,
	t *template.Template,
	data interface{},
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if err := t.Execute(w, data); err != nil {
		internalServerError(w, err)
	}
}
