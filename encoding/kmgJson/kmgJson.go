package kmgJson

import (
	"encoding/json"
	"github.com/bronze1man/kmg/encoding/kmgYaml"
	"github.com/bronze1man/kmg/typeTransform"
	"io/ioutil"
	"os"
)

func ReadFile(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

//读取json文件,并修正json的类型问题(map key 必须是string的问题)
func ReadFileTypeFix(path string, obj interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var iobj interface{}
	err = json.Unmarshal(b, &iobj)
	if err != nil {
		return err
	}
	return typeTransform.Transform(iobj, obj)
}

func WriteFile(path string, obj interface{}) (err error) {
	out, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

//写入json文件,并修正json的类型问题(map key 必须是string的问题)
func WriteFileTypeFix(path string, obj interface{}) (err error) {
	//a simple work around
	yamlBytes, err := kmgYaml.Marshal(obj)
	if err != nil {
		return
	}
	var yamlData interface{}
	err = kmgYaml.Unmarshal(yamlBytes, &yamlData)
	if err != nil {
		return
	}
	yamlData, err = kmgYaml.Yaml2JsonTransformData(yamlData)
	if err != nil {
		return
	}
	out, err := json.Marshal(yamlData)
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, out, os.FileMode(0777))
}

func UnmarshalNoType(r []byte) (interface{}, error) {
	var obj interface{}
	err := json.Unmarshal(r, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
