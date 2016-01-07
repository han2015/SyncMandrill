
// sync_mandrill_templates.go syncs templats from mailchimp account to our test/dev mandrill accounts.
// Production templates could updated by providing two environment variables PROD_KEY.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mattbaird/gochimp"
)

type account struct {
	Email  string //mandrill account
	APIKey string //mandrill APIkey
}

type Config struct {
	MailChimp struct {
		APIKey string //mailChimp APIkey
	}
	Official account           //project's official account of mandrill.
	Accounts []account         //mandrill
	Slugs    map[string]string //key is mailchimp template name and value is mandrill slug value.
}

var prod bool

func main() {
	flag.BoolVar(&prod, "o", false, ":sync the MailChimp templates to the official account of Mandrill.")
	flag.Parse()

	start := time.Now()

	var config Config
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(content, &config); err != nil {
		panic(err)
	}
	fmt.Println(config)
	return

	key := config.MailChimp.APIKey
	if key == "" {
		log.Println("Please specify PROD_MAILCHIMP_KEY")
		os.Exit(1)
	}
	mailchimp := gochimp.NewChimp(key, true)

	var accounts []account
	if prod {
		log.Println("Updating Production Templates")
		accounts = append(accounts, config.Official)
	} else {
		log.Println("Updating Dev/Test Templates")
		accounts = config.Accounts
	}

	slugs := config.Slugs

	chimpList, err := mailchimp.TemplatesList(gochimp.TemplatesList{
		Types:   gochimp.TemplateListType{User: true, Gallery: true, Base: true},
		Filters: gochimp.TemplateListFilter{IncludeDragAndDrop: true},
	})
	if err != nil {
		panic(err)
	}

	for _, tmpl := range chimpList.User {
		slug, ok := slugs[tmpl.Name]
		if !ok {
			continue
		}
		log.Println("Slug:", slug)

		info, err := mailchimp.TemplatesInfo(gochimp.TemplateInfo{
			TemplateID: tmpl.Id,
			Type:       "user",
		})
		if err != nil {
			panic(err)
		}

		for _, account := range accounts {
			log.Println("Account:", account.Email)
			mandril, err := gochimp.NewMandrill(account.APIKey)
			if err != nil {
				panic(err)
			}

			var exists bool
			_, err = mandril.TemplateInfo(slug)
			if err != nil {
				if merr, ok := err.(gochimp.MandrillError); !ok && merr.Name != "Unknown_Template" {
					panic(err)
				}
			} else {
				exists = true
			}

			if exists {
				log.Println("Action: Update")
				_, err = mandril.TemplateUpdate(slug, info.Source, true)
			} else {
				log.Println("Action: Add")
				_, err = mandril.TemplateAdd(slug, info.Source, true)
			}
			if err != nil {
				panic(err)
			}
		}
	}
	log.Println("Took", time.Now().Sub(start))
}
![Build Status](https://travis-ci.org/lib/pq.png?branch=master)](https://travis-ci.org/lib/pq)
