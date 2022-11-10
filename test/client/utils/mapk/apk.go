package mapk

import (
	"errors"
	"fmt"
	"github.com/shogo82148/androidbinary/apk"
	"io"
	"os"
)

func GetNewApkInfo(updateUrl string) (err error) {
	savePath := "./tt.apk"

	//resp, err := http.Get(updateUrl)
	//if err != nil {
	//	err = errors.New("http get mapk error: " + err.Error())
	//	return
	//}
	//defer resp.Body.Close()
	//// 保存文件
	//err = saveFile(savePath, resp.Body)
	//if err != nil {
	//	err = errors.New("saveFile error: " + err.Error())
	//	return
	//}

	// 解包，获取包信息
	err = unpack(savePath)
	if err != nil {
		fmt.Println("unpack err: ", err)
	}
	return
}

func unpack(filePath string) error {
	pkg, err := apk.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer pkg.Close()
	mf := pkg.Manifest()

	code, _ := mf.VersionCode.Int32()
	versionName, _ := mf.VersionName.String()
	pkgName, _ := mf.Package.String()

	fmt.Println(code)
	fmt.Println(versionName)
	fmt.Println(pkgName)
	fmt.Println(mf.SDK.Target.MustInt32())
	fmt.Println(mf.SDK.Min.MustInt32())
	fmt.Println(mf.SDK.Max.MustInt32())
	fmt.Println(mf.Instrument.Name.String())
	fmt.Println(mf.Instrument.Target.String())

	//for _, val := range mf.UsesPermissions {
	//	fmt.Println(val.Name.String())
	//}

	fmt.Println(mf.App.Name.String())
	fmt.Println(mf.App.Description.String())
	fmt.Println(mf.App.Label.String())
	fmt.Println(mf.App.Icon.String())
	fmt.Println(mf.App.Logo.String())
	fmt.Println(mf.App.BackupAgent.String())
	fmt.Println(mf.App.ManageSpaceActivity.String())
	fmt.Println(mf.App.Process.String())
	fmt.Println(mf.App.RequiredAccountType.String())
	fmt.Println(mf.App.RestrictedAccountType.String())
	fmt.Println(mf.App.TaskAffinity.String())
	fmt.Println(mf.App.Theme.String())

	for _, val := range mf.App.MetaData {

		fmt.Print(val.Name.String())
		fmt.Println(val.Value.String())
	}

	return nil
}

func saveFile(path string, saveInfo io.Reader) error {
	localFile, err := os.Create(path)
	defer localFile.Close()
	if err != nil {
		err = errors.New("os.create() err: " + err.Error())
		return err
	}

	_, err = io.Copy(localFile, saveInfo)
	if err != nil {
		err = errors.New("io.Copy() err: " + err.Error())
	}
	return err
}
