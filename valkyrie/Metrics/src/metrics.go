package main

import (
	"os"
	"fmt"
	"time"
	"net/http"
	"strconv"
	"strings"
//	"cipherize"

	"math/rand"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/dgraph-io/badger"
)

type ActionMessage struct {
        Host string
        Action string
        ActionLevel string
        Contact string
        MiscParams string
        OnRemote string
        Image string
	GitRepo string
	Epoch string
}

var db *badger.DB

var license string
var business string

func handleWhoAreYou(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Valkyrie Metrics")
}

func handlePing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func handleDescription(w http.ResponseWriter, r *http.Request) {
	html := "Valkyrie Metrics - Tool for handling reporting metrics for actions taken\n"
	fmt.Fprintf(w, html)
}

func handleRecord(w http.ResponseWriter, r *http.Request) {
        Host := r.URL.Query().Get("host")
        Action := r.URL.Query().Get("action")
        ActionLevel := r.URL.Query().Get("actionlevel")
        Contact := r.URL.Query().Get("contact")
        MiscParams := r.URL.Query().Get("miscparams")
        onremote := r.URL.Query().Get("onremote")
        Image := r.URL.Query().Get("image")
	GitRepo := r.URL.Query().Get("gitrepo")

        if len(Host) == 0 {
                fmt.Fprintf(w, "failed to capture: missing host data")
                return
        }

        if len(Action) == 0 {
                fmt.Fprintf(w, "failed to capture: missing action data")
                return
        }

        if len(ActionLevel) == 0 {
                fmt.Fprintf(w, "failed to capture: missing action level")
                return
        }

        if len(Contact) == 0 {
                fmt.Fprintf(w, "failed to capture: missing contact")
		return
        }

        if len(MiscParams) == 0 {
                fmt.Fprintf(w, "failed to capture: missing miscparams")
		return
        }

        if len(onremote) == 0 {
                fmt.Fprintf(w, "failed to capture: missing onremote")
		return
        }

        if len(Image) == 0 {
                fmt.Fprintf(w, "failed to capture: missing image")
		return
        }

	if len(GitRepo) == 0 {
		fmt.Fprintf(w, "failed to capture: missing gitrepo")
	}

	now := time.Now().Unix()
	snow := strconv.FormatInt(now, 10)

        action := ActionMessage{}
        action.Host = Host
        action.Action = Action
        action.ActionLevel = ActionLevel
        action.Contact = Contact
        action.MiscParams = MiscParams
        action.OnRemote = onremote
        action.Image = Image
	action.GitRepo = GitRepo
	action.Epoch = snow 

	encoded, err := json.Marshal(action)

	db.Update(func(tx *badger.Txn) error {
		if err := tx.Set([]byte(snow), encoded); err != nil {
			fmt.Fprintf(w, "error inserting metric: " + err.Error())
			return err
		}

		return err 
	})

	fmt.Fprintf(w, "successfully logged metric")
}

func handleMetrics(w http.ResponseWriter, r *http.Request) {
        actions := []ActionMessage{}

	db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var val []byte	
			item := it.Item()
			err := item.Value(func(v []byte) {
				val = v
			})
			if err != nil {
				return err 
			}

			tmpaction := ActionMessage{}
			err = json.Unmarshal(val, &tmpaction)
			actions = append(actions, tmpaction)
		}

		return nil
	})

	jsn, err := json.Marshal(actions)
	if err != nil {
		fmt.Fprintf(w, "failed to marshal json")
		return
	}

	fmt.Fprintf(w, string(jsn))
}

func handleViewSpecific(w http.ResponseWriter, r *http.Request) {
        actions := []ActionMessage{}

	db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var val []byte	
			item := it.Item()
			err := item.Value(func(v []byte) {
				val = v
			})
			if err != nil {
				return err 
			}

			tmpaction := ActionMessage{}
			err = json.Unmarshal(val, &tmpaction)
			actions = append(actions, tmpaction)
		}

		return nil
	})

	jsn, err := json.Marshal(actions)
	if err != nil {
		fmt.Fprintf(w, "failed to marshal json")
		return
	}

	fmt.Fprintf(w, string(jsn))
}

func handleViewDay(w http.ResponseWriter, r *http.Request) {
        actions := []ActionMessage{}

        now := time.Now().Unix()

	db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var val []byte	
			epoch := it.Key()
			item := it.Item()
			err := item.Value(func(v []byte) {
				val = v
			})
			if err != nil {
				return err 
			}

			tmpaction := ActionMessage{}
			err = json.Unmarshal(val, &tmpaction)
			then, err := strconv.ParseInt(epoch, 10, 64)
			if err != nil {
				fmt.Fprintf(w, "failed,couldn't convert time")
				return
			}

			

			actions = append(actions, tmpaction)
		}

		return nil
	})

	jsn, err := json.Marshal(actions)
	if err != nil {
		fmt.Fprintf(w, "failed to marshal json")
		return
	}

	fmt.Fprintf(w, string(jsn))
}

func handleViewMonth(w http.ResponseWriter, r *http.Request) {
        actions := []ActionMessage{}

	db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var val []byte	
			item := it.Item()
			err := item.Value(func(v []byte) {
				val = v
			})
			if err != nil {
				return err 
			}

			tmpaction := ActionMessage{}
			err = json.Unmarshal(val, &tmpaction)
			actions = append(actions, tmpaction)
		}

		return nil
	})

	jsn, err := json.Marshal(actions)
	if err != nil {
		fmt.Fprintf(w, "failed to marshal json")
		return
	}

	fmt.Fprintf(w, string(jsn))
}

func handleViewYear(w http.ResponseWriter, r *http.Request) {
        actions := []ActionMessage{}

	db.View(func(tx *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := tx.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			var val []byte	
			item := it.Item()
			err := item.Value(func(v []byte) {
				val = v
			})
			if err != nil {
				return err 
			}

			tmpaction := ActionMessage{}
			err = json.Unmarshal(val, &tmpaction)
			actions = append(actions, tmpaction)
		}

		return nil
	})

	jsn, err := json.Marshal(actions)
	if err != nil {
		fmt.Fprintf(w, "failed to marshal json")
		return
	}

	fmt.Fprintf(w, string(jsn))
}

func VerifyLicense(key string, business string) bool {
	secretkey := business
	padding := "nullsoftllcvalkyrie"

	if len(secretkey) < 16 {
		mdiff := 16 - len(secretkey)
		newchars := padding[0:mdiff]
		secretkey += newchars
	}

	if len(secretkey) > 16 {
		secretkey = secretkey[0:16]
	}


	nkey := strings.Replace(key, "-", "", -1)
	if len(nkey) != 24 {
		return false
	}

	now := time.Now().Unix()

	group1 := nkey[0:8]
	group2 := nkey[8:16]
	group3 := nkey[16:24]

	ogeneratedon, err := strconv.ParseInt("0x" + group1, 0, 64)
	if err != nil {
		return false
	}

	omid, err := strconv.ParseInt("0x" + group2, 0, 64)
	if err != nil {
		return false
	}

	oexptime, err := strconv.ParseInt("0x" + group3, 0, 64)
	if err != nil {
		return false
	}

	generatedon := ogeneratedon ^ omid
	exptime := oexptime ^ ogeneratedon
	diff := exptime - generatedon

	rand.Seed(generatedon)
	mid := rand.Int63n(generatedon - 5878423) + 5878423

	if mid != omid {
		return false
	}

	if diff > 31556952 {
		return false
	}

	if now > exptime {
		return false
	}

	return true
}

func main() {
	var err error

	license = os.Getenv("license")
	business = os.Getenv("business")
	dbpath := os.Getenv("dbpath")

	if len(license) == 0 {
		fmt.Println("Missing License Key")
		return
	}

	if len(business) == 0 {
		fmt.Println("Missing Business Name")
		return
	}

	if len(dbpath) == 0 {
		dbpath = "/tmp"
	}

	if ! VerifyLicense(license, business) {
		fmt.Println("License Is Invalid")
		return
	}

        opts := badger.DefaultOptions
        opts.Dir = dbpath + "/metrics"
        opts.ValueDir = dbpath + "/metrics"
        db, err = badger.Open(opts)
	if err != nil {
		fmt.Println("failed to open metrics database: " + err.Error())
		return
	}
        defer db.Close()

	router := mux.NewRouter()
        router.HandleFunc("/whoareyou", handleWhoAreYou)
        router.HandleFunc("/ping", handlePing)
	router.HandleFunc("/description", handleDescription)
	router.HandleFunc("/record", handleRecord)
        router.HandleFunc("/view", handleMetrics)
	router.HandleFunc("/viewspecific", handleViewSpecific)
	router.HandleFunc("/viewday", handleViewDay)
	router.HandleFunc("/viewmonth", handleViewMonth)
	router.HandleFunc("/viewyear", handleViewYear)

        err = http.ListenAndServe(":8095", router)
        if err != nil {
                fmt.Println("ListenAndServe: ", err)
	}
}
