package jsonsaver

import (
	"encoding/json"
	"io/ioutil"
	"universalScraper/product/globaltypes"
)

func MakeJsonAndSave(a interface{}, fileName string) error {
	jsonData, err := json.MarshalIndent(a, "", "\t")
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(fileName, jsonData, 0600); err != nil {
		return err
	}
	return nil
}
func CopyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile(dst, data, 0600); err != nil {
		return err
	}
	return nil
}
func GenerateCategories(allProducts []globaltypes.Product) (allCategories []globaltypes.ProductCategory) {
	allCategories = []globaltypes.ProductCategory{}
	categories := map[string]globaltypes.ProductCategory{}
	for _, v := range allProducts {
		categories[v.InnerCategoryID] = globaltypes.ProductCategory{
			InnerID: v.InnerCategoryID,
			Name:    v.InnerCategory,
			StoreID: v.StoreID,
		}
	}

	for _, v := range categories {
		allCategories = append(allCategories, v)
	}
	return
}
