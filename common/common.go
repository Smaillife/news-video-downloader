package common

//初始化config结构体
type Config struct {
    News struct {
        TokenUrl string
        SaveDir string
        M3u8FileName string
        MaxThreads int
        RefreshInterval int
        FakeIp string
        FakeUA string
        CTimeout int
        RWTimeout int

    }
}

type FileLink struct {
    Url string
    Cookie []string
    Name string
}

var Cfg Config
var TsQueue = make(chan *FileLink, 1)
