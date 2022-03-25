package globaltypes

type Category struct {
	CategoryID   string `yaml:"category_id"`
	CategoryName string `yaml:"category_name"`
	CategoryLink string `yaml:"category_link"`
}
type Product struct {
	InnerID         string         `json:"inner_id"`
	URL             string         `json:"url"`
	Name            string         `json:"name"`
	Price           string         `json:"price"`
	Brand           string         `json:"brand"`
	HasExist        bool           `json:"has_exist"`
	MainImage       string         `json:"main_image"`
	InnerCategory   string         `json:"inner_category"`
	InnerCategoryID string         `json:"inner_category_id"`
	CategoryID      string         `json:"category_id"`
	Images          []string       `json:"images"`
	Description     string         `json:"desc"`
	PropertyText    string         `json:"property_text"`
	Properties      map[string][]S `json:"properties"`
	StoreID         int            `json:"store_id"`
}
type GrabberConfig struct {
	WithPagination   bool   `yaml:"with_pagination"`
	ProductSelector  string `yaml:"product_selector"`
	NextPageSelector string `yaml:"next_page_selector"`
	NextPageSign     string `yaml:"next_page_sign"`
}
type ProductConfig struct {
	NameSelector          string `yaml:"name_selector"`
	NameRegexp            string `yaml:"name_regexp"`
	PriceSelector         string `yaml:"price_selector"`
	PriceRegexp           string `yaml:"price_regexp"`
	BrandSelector         string `yaml:"brand_selector"`
	BrandRegexp           string `yaml:"brand_regexp"`
	HasExistSelector      string `yaml:"has_exist_selector"`
	HasExistRegexp        string `yaml:"has_exist_regexp"`
	HasExistTrueValue     string `yaml:"has_exist_true_value"`
	MainImageSelector     string `yaml:"main_image_selector"`
	MainImageAttr         string `yaml:"main_image_attr"`
	ImagesSelector        string `yaml:"images_selector"`
	ImagesAttr            string `yaml:"images_attr"`
	DescriptionSelector   string `yaml:"desc_selector"`
	PropertyTextSelector  string `yaml:"property_text_selector"`
	PropertySelector      string `yaml:"property_selector"`
	PropertyHeadSelector  string `yaml:"property_head_selector"`
	PropertyHeadRegexp    string `yaml:"property_head_regexp"`
	PropertyNameSelector  string `yaml:"property_name_selector"`
	PropertyNameRegexp    string `yaml:"property_name_regexp"`
	PropertyValueSelector string `yaml:"property_value_selector"`
	PropertyValueRegexp   string `yaml:"property_value_regexp"`
}
type S struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type ProductCategory struct {
	InnerID    string `json:"inner_id"`
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	StoreID    int    `json:"store_id"`
}
