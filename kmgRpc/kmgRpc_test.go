package kmgRpc

import (
	"github.com/bronze1man/kmg/kmgCmd"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgRpc/testPackage"
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func TestReflectToTplConfig(t *testing.T) {
	conf := reflectToTplConfig(
		GenerateRequest{
			Object:               &testPackage.Demo{},
			ObjectName:           "Demo",
			OutFilePath:          "testPackage/generated.go",
			OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
		},
	)
	kmgTest.Equal(len(conf.ApiList), 2)
	kmgTest.Equal(conf.ApiList[0].Name, "DemoFunc2")
	kmgTest.Equal(conf.ApiList[1].Name, "PostScoreInt")
}

func TestMustGenerateCode(t *testing.T) {
	kmgFile.MustDeleteFile("testPackage/generated.go")
	MustGenerateCode(GenerateRequest{
		Object:               &testPackage.Demo{},
		ObjectName:           "Demo",
		OutFilePath:          "testPackage/generated.go",
		OutPackageImportPath: "github.com/bronze1man/kmg/kmgRpc/testPackage",
	})
	kmgCmd.CmdString("kmg go test").SetDir("testPackage").Run()

}
func TestTplGenerateCode(t *testing.T) {
	out := tplGenerateCode(tplConfig{
		OutPackageName: "testPackage",
		OutKeyBase64:   "wwbo0EGSB6IVKFEy4dH6my1DIaxCCtzPUM9vfx2Hbog=",
		ObjectName:     "Demo",
		ObjectTypeStr:  "*Demo",
		ApiList: []Api{
			{
				Name: "PostScoreInt",
				InArgsList: []ArgumentNameTypePair{
					{
						Name:          "LbId",
						ObjectTypeStr: "string",
					},
					{
						Name:          "Score",
						ObjectTypeStr: "int",
					},
				},
				OutArgsList: []ArgumentNameTypePair{
					{
						Name:          "Info",
						ObjectTypeStr: "string",
					},
					{
						Name:          "err",
						ObjectTypeStr: "error",
					},
				},
			},
		},
		ImportPathMap: map[string]bool{
			"encoding/json": true,
			"errors":        true,
			"fmt":           true,
			"github.com/bronze1man/kmg/encoding/kmgBase64": true,
			"github.com/bronze1man/kmg/kmgCrypto":          true,
			"github.com/bronze1man/kmg/kmgLog":             true,
			"github.com/bronze1man/kmg/kmgNet/kmgHttp":     true,
			"net/http": true,
			"bytes":    true,
		},
	})
	kmgFile.MustDeleteFile("testPackage/generated.go")
	kmgFile.MustWriteFileWithMkdir("testPackage/generated.go", out)
	kmgCmd.CmdString("kmg go test").SetDir("testPackage").Run()
}
