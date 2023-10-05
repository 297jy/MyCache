package config

import (
	"bufio"
	"gomemory/lib/logger"
	"gomemory/lib/utils"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type CacheServerProperties struct {
	RunID        string `cfg:"runid"` // runID always different at every exec.
	Bind         string `cfg:"bind"`
	Port         int    `cfg:"port"`
	Dir          string `cfg:"dir"`
	AnnounceHost string `cfg:"announce-host"`

	CfPath        string `cfg:"cf,omitempty"`
	ClusterEnable bool   `cfg:"cluster-enable"`

	ClusterAsSeed bool   `cfg:"cluster-as-seed"`
	ClusterSeed   string `cfg:"cluster-seed"`

	ClusterEnabled string `cfg:"cluster-enabled"` // Not used at present.
	Self           string `cfg:"self"`
}

type ServerInfo struct {
	StartUpTime time.Time
}

var CacheProperties *CacheServerProperties

var EachTimeServerInfo *ServerInfo

func (p *CacheServerProperties) AnnounceAddress() string {
	return p.AnnounceHost + ":" + strconv.Itoa(p.Port)
}

func SetupCacheConfig(configFileName string) {
	file, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	CacheProperties = &CacheServerProperties{}
	CacheProperties = parse[*CacheServerProperties](file, CacheProperties)
	CacheProperties.RunID = utils.RandString(40)

	configFilePath, err := filepath.Abs(configFileName)
	if err != nil {
		panic(err)
	}
	CacheProperties.CfPath = configFilePath
	if CacheProperties.Dir == "" {
		CacheProperties.Dir = "."
	}
}

func readRawMap(src io.Reader) map[string]string {
	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		// # 代表配置文件中的注释
		if len(line) > 0 && strings.TrimLeft(line, " ")[0] == '#' {
			continue
		}
		// 将每行配置解析成 键值对
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 {
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}
	return rawMap
}

func parse[T *CacheServerProperties](src io.Reader, config T) T {
	//config := &ServerProperties{}
	rawMap := readRawMap(src)

	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		// 动态获取类型标签
		key, ok := field.Tag.Lookup("cfg")
		if !ok || strings.TrimLeft(key, " ") == "" {
			key = field.Name
		}
		value, ok := rawMap[strings.ToLower(key)]
		if ok {
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == value
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					slice := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(slice))
				}
			}
		}

	}

	return config
}

func GetTempDir() string {
	return CacheProperties.Dir + "/tmp"
}
