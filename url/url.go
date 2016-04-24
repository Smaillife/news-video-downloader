package url

/**
 * The following package contains a simple url caller function.
 *
 * @author Fahad Zia Syed <fzia@folio3.com>
 * @edit fengguanjin <412816322@qq.com>
 */
import (
    "io/ioutil"
    "net/http"
    "encoding/json"
    "strings"
    "../common"
    "net"
    "time"
    urllib "net/url"
    utils "../utils"
    log "github.com/cihub/seelog"
    "path/filepath"
)

// The generic dispatch function which call the API url
// and returns the parsed response
func Dispatch(url string, method string, retry bool) {
    status, vdoLink, _ := GetUrl(url, method, nil, true)
    if !status {
        log.Warn("get vdo link failed")
        return
    }
    ExtractUrlFromJson(vdoLink)
}


func GetUrl(url string, method string, cookie []string, retry bool) (bool, []byte, []string) {
    defer log.Flush()
    client := &http.Client{
        Transport: &http.Transport{
            Dial: func(netw, addr string) (net.Conn, error) {
                deadline := time.Now().Add(time.Duration(common.Cfg.News.RWTimeout + common.Cfg.News.CTimeout) * time.Second)
                c, err := net.DialTimeout(netw, addr, time.Second * time.Duration(common.Cfg.News.CTimeout))
                if err != nil {
                    return nil, err
                }
                c.SetDeadline(deadline)
                return c, nil
            },
        },
    }
    request, err := http.NewRequest(method, url, nil)
    if err != nil {
        log.Warn("http client init failed, err:", err)
    }
    if cookie != nil {
        cookieS := ""
        for _, cookieE := range cookie {
            cookieS = cookieS + cookieE + " "
        }
        request.Header.Add("Cookie", cookieS)
    }
    request.Header.Add("CLIENTIP", common.Cfg.News.FakeIp)
    request.Header.Add("User-Agent", common.Cfg.News.FakeUA)

    resp, err := client.Do(request)

    //if there was error in response
    //retry. If retry is true send to failed
    if err != nil {
        if retry == true {
            //retry only once
            return GetUrl(url, method, cookie, false)
        }

        //if failed after retry. send the customer id to failed channel
        log.Warn("fetch url:", url , " failed, err: ", err)
        return false, nil, nil
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    rCookie := resp.Header["Set-Cookie"]

    returnBool := true
    if err != nil {

        log.Warn("read content failed, url:", url , " failed, err: ", err)
        returnBool = false
    }

    return returnBool, body, rCookie
}

func ExtractUrlFromJson(str []byte) []byte {
    type TokenJson struct {
        Url string `json:url`
    }
    var token TokenJson

    json.Unmarshal(str, &token)
    ret := strings.Replace(token.Url, "\\/", "/", -1)
    getVdoLinks(ret)
    return []byte(ret)

}

func getVdoLinks(vdoUrl string) {
    status, ret, cookie := GetUrl(vdoUrl, "GET", nil, true)
    urlInfo, err := urllib.Parse(vdoUrl)
    if err != nil {
        log.Critical("Video chunklist list: ", urlInfo, " fetch content failed",
            "vdolink: ", vdoUrl)
        return
    }
    if status == true {
        m3u8BitRate := strings.Split(string(ret), "\n")
        m3u8File, err := urlInfo.Parse(m3u8BitRate[7])
        if err != nil {
            log.Critical("Video m3u8 link: ", m3u8File, " fetch content failed",
                "M3u8 chunklist: ", m3u8BitRate)
            return
        }
        tsStatus, tsList, _ := GetUrl(m3u8File.String(), "GET", cookie, true)
        if !tsStatus {
            log.Critical("Video ts video list fetch failed")
            return
        }
        err = utils.SaveFileDisk(common.Cfg.News.SaveDir, common.Cfg.News.M3u8FileName, tsList, true)
        if err != nil {
            log.Error("url: ", common.Cfg.News.TokenUrl,
                " M3u8: ", filepath.Join(common.Cfg.News.SaveDir, common.Cfg.News.M3u8FileName), " create failed, err: ", err)
        } else {
            log.Info("url: ", common.Cfg.News.TokenUrl,
                " M3u8: ", filepath.Join(common.Cfg.News.SaveDir, common.Cfg.News.M3u8FileName), " create success")
        }
        tsListArr := strings.Split(string(tsList), "\n")
        for _, line := range tsListArr {
            line = strings.Trim(line, " ")
            if len(line) > 1 && line[0:1] != "#" {
                tsLink, err := urlInfo.Parse(line)
                if err != nil {
                    log.Critical("tslink: ", line, " generate failed, err: ", err)
                }
                log.Info("tslink: ", tsLink, " is ready  inserted into tsqueue")
                common.TsQueue <- &common.FileLink{tsLink.String(), cookie, line}
                log.Info("tslink: ", tsLink, " has inserted into tsqueue")
            }
        }
    }
}
