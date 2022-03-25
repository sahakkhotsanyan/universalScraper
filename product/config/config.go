package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
	"universalScraper/product/globaltypes"
)

type Config struct {
	configPath       string
	SiteName         string                    `yaml:"siteName"`
	RepositoryFolder string                    `yaml:"repositoryFolder"`
	StoreID          int                       `yaml:"storeID"`
	Grabber          globaltypes.GrabberConfig `yaml:"grabber"`
	Categories       []globaltypes.Category
	Product          globaltypes.ProductConfig `yaml:"product"`
}

func (c *Config) MakeConfigFile() error {
	marshal, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.configPath+"/config.yaml", marshal, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.configPath+"/categories.scat", []byte{}, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) getConfigFile() error {
	cnf, err := ioutil.ReadFile(c.configPath + "/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(cnf, c)
	if err != nil {
		return err
	}
	return nil
}
func (c *Config) getCategoryFile() error {
	cnf, err := ioutil.ReadFile(c.configPath + "/categories.scat")
	if err != nil {
		return err
	}
	for _, s := range strings.Split(string(cnf), "\n") {
		s2 := strings.Split(s, "{}{}")
		c.Categories = append(c.Categories,
			globaltypes.Category{
				CategoryID:   s2[0],
				CategoryName: s2[2],
				CategoryLink: s2[1]},
		)
	}
	return nil
}
func (c *Config) GetConfig() Config {
	if c.getConfigFile() != nil {
		return Config{}
	}
	if c.getCategoryFile() != nil {
		return Config{}
	}
	return *c
}
func NewConfig(configPath string) *Config {
	return &Config{configPath: configPath}
}
