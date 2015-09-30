package kmgViewResource

import (
	"github.com/bronze1man/kmg/encoding/kmgJson"
	"github.com/bronze1man/kmg/kmgCache"
	"github.com/bronze1man/kmg/kmgConfig"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgNet/kmgHttp"
	"github.com/bronze1man/kmg/kmgXss"
	"strings"
	"sync"
)

type Generated struct {
	Name                 string //拿去做缓存控制的.
	GeneratedJsFileUrl   string
	GeneratedCssFileUrl  string
	GeneratedJsFileName  string
	GeneratedCssFileName string
	GeneratedUrlPrefix   string // 末尾不包含 /
	RequestImportList    []string

	locker     sync.Mutex
	cachedInfo resourceBuildToDirResponse
}

type htmlTplData struct {
	urlPrefix  string
	jsFileUrl  string
	cssFileUrl string
}

func (g *Generated) HtmlRender() string {
	data := g.GetHtmlTplData()
	return `<script>
		function getResourceUrlPrefix(){
			return ` + kmgXss.Jsonv(data.urlPrefix) + `
		}
	</script>
	<link rel="stylesheet" href="` + kmgXss.H(data.cssFileUrl) + `">
	<script src="` + kmgXss.H(data.jsFileUrl) + `"></script>`
}

func (g *Generated) HeaderHtml() string {
	data := g.GetHtmlTplData()
	return `<script>
		function getResourceUrlPrefix(){
			return ` + kmgXss.Jsonv(data.urlPrefix) + `
		}
	</script>
	<link rel="stylesheet" href="` + kmgXss.H(data.cssFileUrl) + `">`
}

func (g *Generated) FooterHtml() string {
	data := g.GetHtmlTplData()
	return `<script src="` + kmgXss.H(data.jsFileUrl) + `"></script>`
}

func (g *Generated) GetHtmlTplData() htmlTplData {
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		g.locker.Lock()
		cachedInfo := g.cachedInfo
		g.locker.Unlock()
		urlPrefix := "/kmgViewResource." + g.Name
		return htmlTplData{
			urlPrefix:  urlPrefix,
			jsFileUrl:  urlPrefix + "/" + cachedInfo.JsFileName,
			cssFileUrl: urlPrefix + "/" + cachedInfo.CssFileName,
		}
	} else {
		urlPrefix := "/kmgViewResource." + g.Name
		return htmlTplData{
			urlPrefix:  urlPrefix,
			jsFileUrl:  urlPrefix + "/" + g.GeneratedJsFileName,
			cssFileUrl: urlPrefix + "/" + g.GeneratedCssFileName,
		}
	}
}

// 返回url前缀,末尾不包含 '/'
func (g *Generated) GetUrlPrefix() string {
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		return "/kmgViewResource." + g.Name
	} else {
		return "/kmgViewResource." + g.Name
	}
}

// TODO 这个init的体验不好.考虑其他实现方式
func (g *Generated) Init() {
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		kmgHttp.MustAddFileToHttpPathToDefaultServer("/kmgViewResource."+g.Name+"/",
			kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/"+g.Name))
	} else {
		// 默认使用反向代理方式提供数据.
		kmgHttp.MustAddUriProxyRefToUriToDefaultServer("/kmgViewResource."+g.Name+"/", g.GeneratedUrlPrefix)
	}
}

// 获取某个资源文件的内容
func (g *Generated) GetContentByName(name string) (b []byte, err error) {
	name = strings.TrimPrefix(name, "/")
	if kmgConfig.HasDefaultEnv() {
		g.recheckAndReloadCache()
		path := kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/" + g.Name + "/" + name)
		return kmgFile.ReadFile(path)
	} else {
		// 默认使用反向代理方式提供数据.
		return kmgHttp.UrlGetContent(g.GeneratedUrlPrefix + "/" + name)
	}

}

func (g *Generated) recheckAndReloadCache() {
	// 加载缓存文件,确定有哪些文件需要检查.
	cachedInfo := &resourceBuildToDirResponse{}
	err := kmgJson.ReadFile(kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_meta/"+g.Name), &cachedInfo)
	if err != nil || cachedInfo.JsFileName == "" {
		g.reloadCache()
		return
	}
	g.locker.Lock()
	g.cachedInfo = *cachedInfo
	g.locker.Unlock()
	kmgCache.MustMd5FileChangeCache("kmgViewResource_"+g.Name, cachedInfo.NeedCachePathList, g.reloadCache)
}

func (g *Generated) reloadCache() {
	debugBuildPath := kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_debug/" + g.Name)
	response := resourceBuildToDir(g.RequestImportList, debugBuildPath)
	response.NeedCachePathList = append(response.NeedCachePathList, debugBuildPath)
	kmgJson.MustWriteFileIndent(kmgConfig.DefaultEnv().PathInProject("tmp/kmgViewResource_meta/"+g.Name), response)
	g.locker.Lock()
	g.cachedInfo = response
	g.locker.Unlock()
}