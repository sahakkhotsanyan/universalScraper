package runner

import (
	"log"
	"os"
	"reflect"
	"sync"
	"time"
	"universalScraper/product/config"
	"universalScraper/product/globaltypes"
	"universalScraper/product/jsonsaver"
	"universalScraper/product/linkgrabber"
	"universalScraper/product/productgrabber"
	"universalScraper/product/tgbot"
)

var (
	globalWg sync.WaitGroup
)

func Init() {
	dir, err := os.ReadDir("configs")
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range dir {
		if !entry.IsDir() {
			continue
		}
		if len(os.Args) > 1 {
			if os.Args[1] != entry.Name() {
				continue
			}
		}
		cnf := config.NewConfig("configs/" + entry.Name())
		globalWg.Add(1)
		go func(cnf *config.Config, globalWg *sync.WaitGroup) {
			defer globalWg.Done()
			var counter = 0
			var estimationCounter = 0
			var allProducts = make([]globaltypes.Product, 0)
			if conf := cnf.GetConfig(); reflect.DeepEqual(conf, config.Config{}) {
				log.Println("Config is empty, making default config")
				err := cnf.MakeConfigFile()
				if err != nil {
					return
				}
				os.Exit(1)
			}
			lg := linkgrabber.NewLinkGrabber()
			var wg sync.WaitGroup
			for _, category := range cnf.Categories {
				lg.Init(category.CategoryID, category.CategoryName, category.CategoryLink, cnf.Grabber.WithPagination,
					cnf.Grabber.ProductSelector, cnf.Grabber.NextPageSelector)
				links, ctg := lg.GrabLinks(cnf.SiteName, cnf.StoreID, &estimationCounter)
				wg.Add(1)
				go func(links []string, ctg globaltypes.Category, cnf *config.Config, counter *int, allProducts *[]globaltypes.Product, wg *sync.WaitGroup) {
					defer wg.Done()
					pg := productgrabber.NewProductGrabber()
					pg.Init(links, ctg, cnf.Product)
					productsFromCategory := pg.Parse(cnf.SiteName, cnf.StoreID, counter)
					*allProducts = append(*allProducts, productsFromCategory...)
				}(links, ctg, cnf, &counter, &allProducts, &wg)
			}
			wg.Wait()
			_, err := os.Stat(cnf.RepositoryFolder + "/" + cnf.SiteName)
			if err != nil {
				if err := os.Mkdir(cnf.RepositoryFolder+"/"+cnf.SiteName, 0777); err != nil {
					log.Println(err)
				}
			}
			if err := jsonsaver.MakeJsonAndSave(allProducts,
				cnf.RepositoryFolder+"/"+cnf.SiteName+"/"+cnf.SiteName+"_products.json"); err != nil {
				log.Fatal(err)
			}
			if err := jsonsaver.MakeJsonAndSave(jsonsaver.GenerateCategories(allProducts),
				cnf.RepositoryFolder+"/"+cnf.SiteName+"/"+cnf.SiteName+"_categories.json"); err != nil {
				log.Fatal(err)
			}
			tgbot.Add(tgbot.WebSiteInfo{Name: cnf.SiteName, Id: cnf.StoreID, ProductCount: counter, State: 2})
		}(cnf, &globalWg)
	}
	globalWg.Wait()
	time.Sleep(time.Second * 5)
}
