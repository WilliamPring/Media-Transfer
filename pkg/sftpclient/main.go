package sftpclient

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/williampring/media-transfer/config"
	"github.com/williampring/media-transfer/pkg/helper"

	"golang.org/x/crypto/ssh"
)

type FilterFiles struct {
	hostFilePath string
	fileName     string
}

func deepLookForContent(path string, contentExt []string) []FilterFiles {
	var a []FilterFiles
	filepath.WalkDir(path, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if helper.Contains(contentExt, filepath.Ext(d.Name())) {
			a = append(a, FilterFiles{hostFilePath: s, fileName: d.Name()})
		}
		return nil
	})
	return a
}

// func checkPath(hostPath string) (error, bool) {
// 	if _, err := os.Stat(hostPath); os.IsNotExist(err) {
// 		// path/to/whatever does not exist
// 		fmt.Println("HOST PATH DOES NOT EXIST!")
// 		return err, false
// 	}
// 	return nil, true
// }

// func filterFiles(path string) []FilterFiles {
// 	pathError, validPath := checkPath(path)
// 	if pathError != nil {
// 		log.Fatal(validPath)
// 		os.Exit(1)
// 	}
// 	hostPathFiles, err := ioutil.ReadDir(path)
// 	if err != nil {
// 		log.Fatal(err)
// 		os.Exit(1)
// 	}
// 	var filterArray []FilterFiles
// 	for _, file := range hostPathFiles {
// 		fmt.Println(file.Name(), file.IsDir())
// 		if !file.IsDir() && !strings.Contains(file.Name(), ".HEIC") {
// 			filterArray = append(filterArray, FilterFiles{hostFilePath: path + "/" + file.Name(), fileName: file.Name()})
// 		}
// 	}
// 	return filterArray
// }

// change to read from env todo
func setupConfig(clientConfig config.Configurations) ssh.ClientConfig {
	var auths []ssh.AuthMethod
	auths = append(auths, ssh.Password(clientConfig.Sftp.Pass))
	sshConfig := ssh.ClientConfig{
		User: clientConfig.Sftp.User,
		Auth: auths,
		// Uncomment to ignore host key check
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConfig
}

func Start(hostPath string, distPath string, clientConfig config.Configurations) {

	//extend to look for mutliple paths
	imageContentMineType := []string{".PNG", ".JPG"}
	filterFiles := deepLookForContent(hostPath, imageContentMineType)
	sshConfig := setupConfig(clientConfig)
	addr := fmt.Sprintf("%s:%d", clientConfig.Sftp.Host, clientConfig.Sftp.Port)

	// Connect to server
	conn, err := ssh.Dial("tcp", addr, &sshConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connecto to [%s]: %v\n", addr, err)
		os.Exit(1)
	}

	defer conn.Close()

	// Create new SFTP client
	sc, err := sftp.NewClient(conn,
		sftp.UseConcurrentReads(true),
		sftp.UseConcurrentWrites(true),
		sftp.MaxConcurrentRequestsPerFile(64),
		sftp.MaxPacketUnchecked(32768),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start SFTP subsystem: %v\n", err)
		os.Exit(1)
	}
	defer sc.Close()
	for _, filterFile := range filterFiles {
		uploadFile(sc, filterFile.hostFilePath, "./media/images/"+filterFile.fileName)
	}
}

// Upload file to sftp server
func uploadFile(sc *sftp.Client, localFile, remoteFile string) (err error) {
	fmt.Fprintf(os.Stdout, "Uploading [%s] to [%s] ...\n", localFile, remoteFile)

	srcFile, err := os.Open(localFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open local file: %v\n", err)
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		sc.Mkdir(path)
	}

	// Note: SFTP To Go doesn't support O_RDWR mode
	dstFile, err := sc.OpenFile(remoteFile, (os.O_WRONLY | os.O_CREATE | os.O_TRUNC))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
		return
	}
	defer dstFile.Close()

	bytes, err := io.Copy(dstFile, srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to upload local file: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d bytes copied\n", bytes)

	return
}
