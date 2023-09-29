package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/samber/lo"
)

const pageTemplatePath = "./resources/home.tpl.html"

// myfiles is a very fancy data store as you can tell
var myfiles = [...]file{
	{
		Owner:          "spiffe://example.com/dataexporter",
		Name:           "customers_export.csv",
		Location:       "Belgium",
		Classification: "CUSTOMERS_PII",
	},
	{
		Owner:       "spiffe://example.com/test",
		Name:        "testdata.csv",
		Location:    "Belgium",
		Environment: "staging",
	},
	{
		Owner:          "hr@example.com",
		Name:           "payslip_e9999_123.pdf",
		Location:       "Belgium",
		Classification: "HR_PII",
		EmployeeID:     "e9999",
	},
	{
		Owner:          "hr@example.com",
		Name:           "payslip_e8888_123.pdf",
		Location:       "USA",
		Classification: "HR_PII",
		EmployeeID:     "e8888",
	},
}

func main() {
	ctx := context.Background()

	services := lo.Must(createServices(ctx))

	ro := httprouter.New()
	ro.GET("/", handleErrors(services.homeHandle))

	const addr = "127.0.0.1:8080"

	srv := http.Server{
		Addr:    addr,
		Handler: ro,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	log.Println("serving on", "http://"+addr)

	lo.Must0(srv.ListenAndServe())
}

func (s services) homeHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	t, err := template.ParseFiles(pageTemplatePath)
	if err != nil {
		return err
	}

	user := s.user.get()
	input := page{
		User:  struct2Map(user),
		Files: make([]displayFile, len(myfiles)),
	}

	json.NewEncoder(os.Stdout).Encode(&input)

	for i, fileToCheck := range myfiles {
		authz, err := s.policy.eval(r.Context(), policyRequest{
			Resource: fileToCheck,
			Subject:  user,
		})
		if err != nil {
			return fmt.Errorf("policy evaluation failed with unexpected error %w", err)
		}

		input.Files[i] = displayFile{
			File:  struct2Map(fileToCheck),
			Authz: authz,
		}
	}

	w.Header().Set("Content-Type", "text/html")
	return t.ExecuteTemplate(w, filepath.Base(pageTemplatePath), &input)
}
