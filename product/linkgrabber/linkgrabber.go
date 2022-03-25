package linkgrabber

import (
	"github.com/gocolly/colly"
	"log"
	"universalScraper/product/globaltypes"
	"universalScraper/product/tgbot"
)

type LinkGrabber struct {
	category      globaltypes.Category
	grabberConfig globaltypes.GrabberConfig
}

func NewLinkGrabber() *LinkGrabber {
	return &LinkGrabber{}
}
func (lg *LinkGrabber) Init(categoryID string, categoryName string, categoryLink string, withPagination bool, productSelector string, nextPageSelector string) {
	lg.category.CategoryID = categoryID
	lg.category.CategoryName = categoryName
	lg.category.CategoryLink = categoryLink
	lg.grabberConfig.WithPagination = withPagination
	lg.grabberConfig.ProductSelector = productSelector
	lg.grabberConfig.NextPageSelector = nextPageSelector
}

func (lg *LinkGrabber) GrabLinks(siteName string, storeID int, estimationCounter *int) ([]string, globaltypes.Category) {
	c := colly.NewCollector()
	links := make([]string, 0)
	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(lg.grabberConfig.ProductSelector, func(_ int, e *colly.HTMLElement) {
			links = append(links, e.Request.AbsoluteURL(e.Attr("href")))
			*estimationCounter++
			tgbot.Add(tgbot.WebSiteInfo{
				Name:         siteName,
				Id:           storeID,
				ProductCount: *estimationCounter,
				State:        3,
			})
			log.Println("Link count:", len(links))
		})
		if lg.grabberConfig.WithPagination {
			if lg.grabberConfig.NextPageSign != "" {
				if e.ChildAttr(lg.grabberConfig.NextPageSelector, "href") != "" && e.ChildText(lg.grabberConfig.NextPageSelector) == lg.grabberConfig.NextPageSign {
					if err := e.Request.Visit(e.Request.AbsoluteURL(e.ChildAttr(lg.grabberConfig.NextPageSelector, "href"))); err != nil {
						log.Printf("error while visiting next page %v", err)
					}
				}
			} else {
				if e.ChildAttr(lg.grabberConfig.NextPageSelector, "href") != "" {
					if err := e.Request.Visit(e.Request.AbsoluteURL(e.ChildAttr(lg.grabberConfig.NextPageSelector, "href"))); err != nil {
						log.Printf("error while visiting next page %v", err)
					}
				}
			}
		}
	})
	if err := c.Visit(lg.category.CategoryLink); err != nil {
		log.Printf("error while visiting category %v", err)
	}
	return links, lg.category
}
