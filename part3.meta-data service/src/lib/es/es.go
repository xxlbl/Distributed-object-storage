package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

func getMetadata(name string, versionId int) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		os.Getenv("ES_SERVER"), name, versionId)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

func SearchLatestVersion(name string) (meta Metadata, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		os.Getenv("ES_SERVER"), url.PathEscape(name))
	// 按照version降序查询
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

//GetMetadata 函数的功能类似 getMetadata，输入对象的名字和版本号返回对象，
//区别在于当 version 为。0时，会调用 SearchLatestVersion 获取当前最新的版本。
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

//PutMetadata 函数用于向 ES 服务上传一个新的元数据。它的4个输入参数对应元数据的4个属性，函数会将它们拼成一个 ES 文档，
//一个ES的文档相当于数据库的一条记录。用 PUT 方法把这个文档上传到 metadata 索引的 objeets 类型，
//且文档 id 由元数据的 name 和 version 拼成，方便我们 GET。
//使用了 op_ type-create 参数，如果同时有多个客户端上传同一个元数据，结果会发生冲突，只有第一个文档被成功创建。之后的 PUT 请求，ES会返回 409 Conflict。
//此时，我们的函数会让版本号加1并递归调用自身继续上传。
func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	//
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

// 版本+1，put
func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e
	}
	return PutMetadata(name, version.Version+1, size, hash)
}

// 返回metadata切片，全部对象+版本
func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		os.Getenv("ES_SERVER"), from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}
	return metas, nil
}

func DelMetadata(name string, version int) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d",
		os.Getenv("ES_SERVER"), name, version)
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

type Bucket struct {
	Key         string
	Doc_count   int
	Min_version struct {
		Value float32
	}
}

type aggregateResult struct {
	Aggregations struct {
		Group_by_name struct {
			Buckets []Bucket
		}
	}
}

func SearchVersionStatus(min_doc_count int) ([]Bucket, error) {
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_search", os.Getenv("ES_SERVER"))
	body := fmt.Sprintf(`
        {
          "size": 0,
          "aggs": {
            "group_by_name": {
              "terms": {
                "field": "name",
                "min_doc_count": %d
              },
              "aggs": {
                "min_version": {
                  "min": {
                    "field": "version"
                  }
                }
              }
            }
          }
        }`, min_doc_count)
	request, _ := http.NewRequest("GET", url, strings.NewReader(body))
	r, e := client.Do(request)
	if e != nil {
		return nil, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var ar aggregateResult
	json.Unmarshal(b, &ar)
	return ar.Aggregations.Group_by_name.Buckets, nil
}

func HasHash(hash string) (bool, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=0", os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return false, e
	}
	b, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(b, &sr)
	return sr.Hits.Total != 0, nil
}

func SearchHashSize(hash string) (size int64, e error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=hash:%s&size=1",
		os.Getenv("ES_SERVER"), hash)
	r, e := http.Get(url)
	if e != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search hash size: %d", r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		size = sr.Hits.Hits[0].Source.Size
	}
	return
}
