package productgrabber

import (
	"crypto/md5"
	"fmt"
	"github.com/gocolly/colly"
	"html"
	"log"
	"regexp"
	"strings"
	"universalScraper/product/globaltypes"
	"universalScraper/product/tgbot"
)

type ProductGrabber struct {
	links    []string
	category globaltypes.Category
	Config   globaltypes.ProductConfig
}

func checkRegexp(reg string, orig string) string {
	if reg != "" {
		re := regexp.MustCompile(reg)
		return strings.Join(re.FindAllString(orig, -1), "")
	}
	return orig
}
func checkForAttr(attr string) string {
	if attr != "" {
		return attr
	}
	return "src"
}
func (p *ProductGrabber) Parse(siteName string, storeID int, counter *int) []globaltypes.Product {
	products := make([]globaltypes.Product, 0)
	c := colly.NewCollector()
	c.OnHTML("html", func(e *colly.HTMLElement) {
		product := globaltypes.Product{
			InnerID:         fmt.Sprintf("%x", md5.Sum([]byte(e.Request.URL.String()))),
			URL:             e.Request.URL.String(),
			Name:            checkRegexp(p.Config.NameRegexp, e.ChildText(p.Config.NameSelector)),
			Price:           getNumericPart(checkRegexp(p.Config.PriceRegexp, e.ChildText(p.Config.PriceSelector))),
			Brand:           checkRegexp(p.Config.BrandRegexp, e.ChildText(p.Config.BrandSelector)),
			MainImage:       e.ChildAttr(p.Config.MainImageSelector, checkForAttr(p.Config.MainImageAttr)),
			InnerCategory:   p.category.CategoryName,
			InnerCategoryID: p.category.CategoryID,
			CategoryID:      p.category.CategoryID,
			Description:     stripHTML(e.ChildText(p.Config.DescriptionSelector)),
			PropertyText:    stripHTML(e.ChildText(p.Config.PropertyTextSelector)),
			Properties:      make(map[string][]globaltypes.S, 0),
			StoreID:         storeID,
		}
		if product.Brand == "" {
			product.Brand = "Other"
		}
		if checkRegexp(p.Config.HasExistRegexp, e.ChildText(p.Config.HasExistSelector)) == p.Config.HasExistTrueValue {
			product.HasExist = true
		}
		if product.Price == "0" || product.Price == "" {
			product.Price = "0"
			product.HasExist = false
		}
		if product.MainImage == "" {
			product.MainImage = e.ChildAttr(p.Config.MainImageSelector, "href")
		}
		if product.MainImage == "" {
			e.ForEach(p.Config.ImagesSelector, func(i int, e *colly.HTMLElement) {
				if i == 0 {
					product.MainImage = e.Request.AbsoluteURL(e.Attr(checkForAttr(p.Config.ImagesAttr)))
				} else {
					product.Images = append(product.Images, e.Request.AbsoluteURL(checkForAttr(p.Config.ImagesAttr)))
				}
			})
		} else {
			e.ForEach(p.Config.ImagesSelector, func(i int, e *colly.HTMLElement) {
				product.Images = append(product.Images, e.Request.AbsoluteURL(checkForAttr(p.Config.ImagesAttr)))
			})
		}
		head := ""
		e.ForEach(p.Config.PropertySelector, func(i int, e *colly.HTMLElement) {
			if p.Config.PropertyHeadSelector == "" { //nolint:nestif
				head = "Ընդհանուր Բնութագիր"
			} else {
				if e.ChildText(p.Config.PropertyHeadSelector) != "" {
					head = checkRegexp(p.Config.PropertyHeadRegexp, e.ChildText(p.Config.PropertyHeadSelector))
				} else {
					name := checkRegexp(p.Config.PropertyNameRegexp, e.ChildText(p.Config.PropertyNameSelector))
					value := checkRegexp(p.Config.PropertyValueRegexp, e.ChildText(p.Config.PropertyValueSelector))
					if head != "" && name != "" && value != "" {
						product.Properties[head] = append(product.Properties[head], globaltypes.S{
							Name:  name,
							Value: value,
						})
					}
				}
			}
		})
		products = append(products, product)
		*counter++
		tgbot.Add(tgbot.WebSiteInfo{
			Name:         siteName,
			Id:           storeID,
			ProductCount: *counter,
			State:        1,
		})
		log.Println("Products -> ", len(products))
	})
	for _, link := range p.links {
		if err := c.Visit(link); err != nil {
			log.Printf("error while visiting category %v", err)
		}
	}
	return products
}
func NewProductGrabber() *ProductGrabber {
	return &ProductGrabber{}
}
func (p *ProductGrabber) Init(links []string, category globaltypes.Category, config globaltypes.ProductConfig) {
	p.links = links
	p.category = category
	p.Config = config
}
func stripHTML(s string) string {
	s = strings.ReplaceAll(s, "</p>", "\n")
	s = html.UnescapeString(s)
	r := regexp.MustCompile("<.*?>")
	l := r.ReplaceAllString(s, " ")
	r = regexp.MustCompile("  +")
	s = strings.TrimSpace(r.ReplaceAllString(l, "\n"))
	r = regexp.MustCompile(`<\!--(\n.*)+-->`)
	s = r.ReplaceAllString(s, "")
	p := regexp.MustCompile(`(\n |\n|\r\n)+`)
	return p.ReplaceAllString(s, "\n")
}
func getNumericPart(number string) string {
	reg := regexp.MustCompile("[^0-9]+")
	return reg.ReplaceAllString(number, "")
}
