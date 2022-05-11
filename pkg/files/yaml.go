/*
Copyright © 2022 Ivan Pedrazas <ipedrazas@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package files

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v2"
)

func SaveAsYaml(path string, out interface{}) error {
	d, err := yaml.Marshal(&out)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Load(configFile string, out interface{}) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("unable to load configuration file: %w", err)
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("unable to parse configuration file: %w", err)
	}
	return nil
}

func LoadTOML(configFile string, out interface{}) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("unable to load configuration file: %w", err)
	}
	if err := toml.Unmarshal(data, out); err != nil {
		return fmt.Errorf("unable to parse configuration file: %w", err)
	}
	return nil
}

func Copy(source string, target string) error {
	bytesRead, err := ioutil.ReadFile(source)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(target, bytesRead, 0644)

	if err != nil {
		return err
	}
	return nil
}
