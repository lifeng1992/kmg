package kmgQiniu

import (
	"os"
	"path/filepath"

	"github.com/bronze1man/kmg/kmgErr"
	"github.com/bronze1man/kmg/kmgFile"
	"github.com/bronze1man/kmg/kmgTask"
)

var Download = DownloadDir

//下载单个文件,到本地,会覆盖本地已经存在的文件,会创建所有父级目录,会使用hash检查文件是否存在.
// @deprecated
func DownloadFileWithHash(ctx *Context, remotePath string, localPath string, hash string) (err error) {
	fhash, err := ComputeHashFromFile(localPath)
	if err == nil && fhash == hash {
		return
	}
	if err != nil && !os.IsNotExist(err) {
		kmgErr.LogErrorWithStack(err)
		return
	}
	if fhash == hash {
		return
	}
	return DownloadFile(ctx, remotePath, localPath)
}

// @deprecated
// TODO 这个会多建一个空目录,这个算安全隐患.
func DownloadFile(ctx *Context, remotePath string, localPath string) (err error) {
	err = kmgFile.MkdirForFile(localPath)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	f, err := os.Create(localPath)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return
	}
	defer f.Close()
	return ctx.DownloadToWriter(remotePath, f)
}

// @deprecated
func DownloadDir(ctx *Context, remoteRoot string, localRoot string) (err error) {
	entries, err := ListPrefix(ctx, remoteRoot)
	if err != nil {
		kmgErr.LogErrorWithStack(err)
		return err
	}
	if len(entries) == 0 {
		return ErrNoFile
	}
	// TODO 多线程提高性能.
	tm := kmgTask.NewLimitThreadErrorHandleTaskManager(ThreadNum, 3)
	defer tm.Close()
	for _, entry := range entries {
		entry := entry
		refPath, err := filepath.Rel(remoteRoot, entry.Key)
		if err != nil {
			kmgErr.LogErrorWithStack(err)
			return err
		}
		tm.AddTask(func() (err error) {
			err = DownloadFileWithHash(ctx, entry.Key, filepath.Join(localRoot, refPath), entry.Hash)
			return err
		})
	}
	tm.Wait()
	return nil
}
