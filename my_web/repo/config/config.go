package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"reflect"
)

func (cf *configReal) InitConfig(x interface{}, filepath string) error {
	if reflect.ValueOf(x).Kind() != reflect.Ptr {
		return errors.New("first param must be a pointer")
	}
	conFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(conFile, x)
	if err != nil {
		return err
	}

	cf.filepath = filepath
	cf.svrConfig = x
	return nil
}

// 使用 ExternalService 接口而不关心其实现
func (cf *configReal) GetConfig(x ConfigInter) interface{} {
	return cf.svrConfig
}

// func MoniterConfig(x ConfigInter, updateConfigChan chan<- *Config) {
// 	x.moniterConfigFile(updateConfigChan)
// }

// // 实时监控config文件
// func (cf *configReal) moniterConfigFile(updateConfigChan chan<- *Config) {
// 	go cf.watchConfigFile(cf.ConfigFilePath, updateConfigChan)
// }

// // 动态更新config文件
// func (cf *configReal) watchConfigFile(configPath string, updateConfigChan chan<- *Config) {
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer watcher.Close()

// 	done := make(chan bool)

// 	go func() {
// 		for {
// 			time.Sleep(1 * time.Second)
// 			select {
// 			case event, ok := <-watcher.Events:
// 				if !ok {
// 					return
// 				}
// 				if event.Op&fsnotify.Write == fsnotify.Write {
// 					err := cf.initConfig()
// 					if err != nil {
// 						log.Println("Error loading config:", err)
// 						continue
// 					}
// 					updateConfigChan <- cf.svrConfig
// 				}
// 			case err, ok := <-watcher.Errors:
// 				if !ok {
// 					return
// 				}
// 				log.Println("Error:", err)
// 			}
// 		}
// 	}()

// 	err = watcher.Add(configPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	<-done
// }
