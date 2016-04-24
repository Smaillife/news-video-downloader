package utils

import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
    log "github.com/cihub/seelog"
    "strings"
    "regexp"
    "strconv"
)


func Init() {
    logger, err := log.LoggerFromConfigAsFile("/Users/fengguanjin/project/news-video-downloader/conf/log.xml")
    if err != nil {
        fmt.Println("log init failed, programme exit", err)
        os.Exit(1)
    }
    log.ReplaceLogger(logger)
}

func SaveFileDisk(rootPath string, fileName string, content []byte, forceDel bool) error {
    fullPath := filepath.Join(rootPath, fileName)
    if !IsFileOrPathExist(rootPath) {
        err := os.Mkdir(rootPath, 0777)
        if err != nil {
            fmt.Println(err, " path create failed")
        }
    }
    if IsFileOrPathExist(fullPath) && !forceDel {
        return nil
    } else if IsFileOrPathExist(fullPath)  && checkSeq(string(content), fullPath) {
        DeleteFile(fullPath)
        err := ioutil.WriteFile(fullPath, content, 0644)
        return err
    }

    if !IsFileOrPathExist(fullPath) {
        err := ioutil.WriteFile(fullPath, content, 0644)
        return err
    }

    return nil
}

/**
 * Removes the file if it exists.
 *
 * @param filename - string - The file to remove
 */
func DeleteFile(filename string) {
    //Check if the file is present by opening it
    dataFile, _ := os.Open(filename)
    if dataFile != nil {
        //the file exists.
        //Delete the file
        dataFile.Close()
        os.Remove(filename)
    }
}

func IsFileOrPathExist(item string) bool {
    _, err := os.Stat(item)
    return err == nil || os.IsExist(err)
}
func checkSeq(cur string, oldPath string) bool {
    curSeq := extractSeq(cur)
    oldContent, err := ioutil.ReadFile(oldPath)
    if err != nil {
        log.Warn("old m3u8 file read failed, err: ", oldPath)
        return true
    }
    oldSeq := extractSeq(string(oldContent))
    if curSeq >= oldSeq {
        return true
    }

    log.Warnf("seq was delayed, old[%d], new[%d]", oldSeq, curSeq)
    return false
}

func extractSeq(str string) (int) {
    ret := strings.Split(str, "\n")
    var digitsRegexp = regexp.MustCompile(`(\d+)`)
    for _, line := range ret {
        if strings.Contains(line, "SEQUENCE") {
            seqRet := digitsRegexp.FindStringSubmatch(line)
            if seqRet[0] != "" {
                index, err := strconv.Atoi(seqRet[0])
                if err != nil {
                    log.Warn("extract seq failed, seq: ", seqRet[0])
                    return 0
                } else {
                    return index
                }

            }
        }
    }
    return 0
}