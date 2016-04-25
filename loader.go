package main

import (
    "time"
    "flag"
    "runtime"
    "./common"
    tool "./utils"
    log "github.com/cihub/seelog"
    gcfg "gopkg.in/gcfg.v1"
    worker "./download"
    url "./url"
    //"fmt"
    "os"
)


var refreshTimer = time.NewTicker(time.Second * 1)
var confFile = flag.String("c", tool.GetCurrentDirectory() + "../conf/config.ini", "set the config path")


func main() {
    flag.Parse()
    tool.Init()
    err := gcfg.ReadFileInto(&common.Cfg, *confFile)
    if err != nil {
        log.Critical("config file load failed, err: ", err)
        os.Exit(1)
    }
    runtime.GOMAXPROCS(common.Cfg.News.MaxThreads)

    refreshTimer = time.NewTicker(time.Second * time.Duration(common.Cfg.News.RefreshInterval))
    go func() {
        for {
            select {
            case <-refreshTimer.C:
                log.Infof("Video M3U8 Url[%s] refresh...", common.Cfg.News.TokenUrl)
                url.Dispatch(common.Cfg.News.TokenUrl, "GET", true)

            }
        }
    }()

    for i := 0; i < common.Cfg.News.MaxThreads; i++ {
        go worker.Worker()
    }

    for {
        log.Flush()
    }

    defer log.Flush()
}
