package download

import (
    "../common"
    "../utils"
    "path/filepath"
    urllib "../url"
    log "github.com/cihub/seelog"

)
func Worker() {
    for {
        ts := <- common.TsQueue
        log.Debug("tslink: ", ts.Url,", tsCookie: ", ts.Cookie, " tsname: ", ts.Name)
        if utils.IsFileOrPathExist(filepath.Join(common.Cfg.News.SaveDir, ts.Name)) {
            return
        }
        ret, tsContent, _ := urllib.GetUrl(ts.Url, "GET", ts.Cookie, true)
        if !ret {
            log.Warnf("ts link[%s] download failed", ts.Url)
        } else {
            log.Infof("ts link[%s] download success", ts.Url)
            err := utils.SaveFileDisk(common.Cfg.News.SaveDir, ts.Name, tsContent, false)
            if err != nil {
                log.Critical("ts file: ", filepath.Join(common.Cfg.News.SaveDir, ts.Name), " save failed, err: ", err)
            } else {
                log.Info("ts file: ", filepath.Join(common.Cfg.News.SaveDir, ts.Name), " save success")
            }
        }
    }
}
