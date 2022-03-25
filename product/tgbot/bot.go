package tgbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"universalScraper/product/tgbot/config"
)

type bot struct {
	url          string
	token        string
	id           int64
	reply        int
	sendMessage  func(message string) error
	sendDocument func(filepath string) error
	editMessage  func(message string) (*http.Response, error)
}
type Response struct {
	Result struct {
		Message_id int
	}
}
type WebSiteInfo struct {
	Name         string
	Id           int
	ProductCount interface{}
	State        int
}

var websites = sync.Map{} //map[int]WebSiteInfo{}
var b bot
var startTime = time.Now()
var BOTMESSAGE string = "Price City PRODUCT updater started at \n" + startTime.Format("01-02-2006 15:04:05 Monday")

func NewBotAPI(token string, s *bot) {
	b := s
	b.url = "https://api.telegram.org/bot"
	b.token = token
	if config.IsDebug() {
		b.id = config.DebugID
	} else {
		b.id = config.StandardID
	}
	b.sendMessage = func(message string) error {
		prpurl := fmt.Sprintf(b.url+b.token+"/sendMessage?parse_mode=markdown&chat_id=%d&text=%s", b.id, url.QueryEscape(message))
		res, err := http.Get(prpurl)
		if err != nil {
			return err
		}
		var result Response
		json.NewDecoder(res.Body).Decode(&result)
		b.reply = result.Result.Message_id
		return nil
	}
	b.editMessage = func(message string) (*http.Response, error) {
		prpurl := fmt.Sprintf(b.url+b.token+"/editMessageText?chat_id=%d&text=%s&parse_mode=markdown&message_id=%d", b.id, url.QueryEscape(message), b.reply)
		res, err := http.Get(prpurl)
		return res, err
	}
	b.sendDocument = func(filepath string) error {
		println("Sending Document...")
		prpurl := fmt.Sprintf(b.url + b.token + "/sendDocument")
		params := map[string]string{
			"chat_id": fmt.Sprintf("%d", b.id),
		}
		s, err := newfileUploadRequest(prpurl, params, "document", filepath)
		client := &http.Client{}
		client.Do(s)
		if err != nil {
			return err
		}
		return nil
	}
}
func SendDocument(filepath string) error {
	return b.sendDocument(filepath)
}
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
func Add(ws WebSiteInfo) {
	websites.Store(ws.Id+ws.State, ws)
}
func Remove(id int) {
	websites.Delete(id)
}
func Init() {
	NewBotAPI(config.Token, &b)
	err := b.sendMessage(BOTMESSAGE)
	if err != nil {
		log.Println(err)
	}
}

var quit chan struct{}

func StartPolling() {
	quit = make(chan struct{})
	go poll()
}
func poll() {

	for {
		display := [][]string{}
		str := ""
		// str := "\nID\tNAME\tPRODUCTS\tSTATE"
		display = append(display, []string{
			"ID",
			"NAME",
			"PRODUCTS",
			"STATE",
		})
		select {
		case <-quit:
			return
		default:
			websites.Range(func(_, value interface{}) bool {
				v := value.(WebSiteInfo)
				var st string
				switch v.State {
				case 1:
					st = "Running"
				case 2:
					st = "ReadyToScrapers"
				case 3:
					st = "Estimating"
				case 4:
					st = "Updated"
				case 5:
					st = "InScrappers"
				case 6:
					st = "UpdStarted..."
				}
				fm := "%d"
				if fmt.Sprintf("%T", v.ProductCount) == "string" {
					fm = "%s"
				}
				display = append(display, []string{
					fmt.Sprintf("[%d]", v.Id),
					v.Name,
					fmt.Sprintf(fm, v.ProductCount),
					st,
				})
				return true
			})
			//for _, v := range websites {
			//
			//	// str += fmt.Sprintf("\n[%d]\t%s\t%d\t%s", v.Id, v.Name, v.ProductCount, st)
			//}
			celllen := len(display[0])
			rowlen := len(display)
			for i := 0; i < celllen; i++ {
				maxlength := 0
				for j := 0; j < rowlen; j++ {
					if len(display[j][i]) > maxlength {
						maxlength = len(display[j][i])
					}
				}
				for k := 0; k < rowlen; k++ {
					if len(display[k][i]) < maxlength {
						display[k][i] = display[k][i] + strings.Repeat(" ", maxlength-len(display[k][i]))
					}

				}
			}
			head := display[0]
			others := display[1:]
			// fmt.Println("head", head)
			// fmt.Println("body", others)
			sort.Slice(others, func(p, q int) bool {
				return others[p][2] < others[q][2]
			})
			display = append(append([][]string{}, head), others...)
			for i := 0; i < len(display); i++ {
				str += strings.Join(display[i], " | ") + "\n"
			}
			dt := time.Now()
			datestr := dt.Format("01-02-2006 15:04:05 Monday")

			b.editMessage(BOTMESSAGE + "``` \n" + str + "\n " + datestr + "\n RUNNING TIME " + fmt.Sprintf("%f", dt.Sub(startTime).Minutes()) + " ```")
			time.Sleep(3 * time.Second)
		}
	}
}
func EndPolling() {
	quit <- struct{}{}
}
