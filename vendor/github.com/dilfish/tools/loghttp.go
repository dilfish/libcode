package tools


import (
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"strings"
)


var ErrPostOne = errors.New("post one error")
var ErrGetStat = errors.New("get stat error")


// HTTPLogger holds api server's url
type RequestLogger struct {
	PostUrl string
	GetUrl string
}


// ErrInfo
type ErrInfo struct {
	Err int `json:"err"`
	Msg string `json:"msg"`
}


// RequestLoggerStat holds logs from start to end
type RequestLoggerStat struct {
	ErrInfo
	MethodCount map[string]int64 `json:"methodCount"`
	PathCount map[string]int64 `json:"pathCount"`
	ClientIPCount map[string]int64 `json:"clientIPCount"`
}


// RequestInfo retrieves all info in http.Request
type RequestInfo struct {
	Id string `json:"_" bson:"_id"` // mongodb id
	Name string `json:"name" bson:"name"` // log name
	Method string `json:"method" bson:"method"` // request method, get, post, put, head etc.
	Path string `json:"path" bson:"path"` // url.Path
	ClientIP string `json:"clientIP" bson:"clientIP"` // client ip
	Time time.Time `json:"time" bson:"time"` // when does the request fired
}


// RequestToInfo makes info from request
func RequestToInfo(req *http.Request, t time.Time) RequestInfo {
	var ri RequestInfo
	ri.Method = req.Method
	ri.Path = req.URL.Path
	arr := strings.Split(req.RemoteAddr, ":")
	if len(arr) == 2 {
		ri.ClientIP = arr[0]
	} else {
		ri.ClientIP = req.RemoteAddr
	}
	ri.Time = t
	return ri
}


// NewRequestLogger gives a new instance
func NewRequestLogger(post, get string) *RequestLogger {
	return &RequestLogger{PostUrl: post, GetUrl: get}
}


// PostOne post one request log to  server
func (hl *RequestLogger) PostOne(req *http.Request) error {
	ri := RequestToInfo(req, time.Now())
	ret := new(ErrInfo)
	err := DoJsonPost(hl.PostUrl, ri, ret)
	if err != nil {
		return err
	}
	if ret.Err != 0 {
		return ErrPostOne
	}
	return nil
}


// DoJsonPost is a general api for http post
func DoJsonPost(uri string, args, ret interface{}) error {
	bt, err := json.Marshal(args)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(bt)
	resp, err := http.Post(uri, "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bt, &ret)
	if err != nil {
		return err
	}
	return nil
}

// StatRequestInfo just specify start and end time
type StatRequestInfo struct {
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
}

// GetStat get log stat from start to end
func (hl *RequestLogger) GetStat(start, end time.Time) (*RequestLoggerStat, error) {
	var sri StatRequestInfo
	sri.Start = start
	sri.End = end
	ret := new(RequestLoggerStat)
	err := DoJsonPost(hl.GetUrl, sri, ret)
	if err != nil {
		return nil, err
	}
	if ret.Err != 0 {
		return nil, ErrGetStat
	}
	return ret, nil
}


///////////////////// server side ///////////////////////////



// ServeRequestLogger is server side object of request logger
type ServeRequestLogger struct {
	DB *ReqLogDB
	Conf MgoConfig
}


// NewServeRequestLogger create instance of server side logger
func NewServeRequestLogger(conf MgoConfig) *ServeRequestLogger {
	db := OpenReqLogDB(conf)
	if db == nil {
		return nil
	}
	return &ServeRequestLogger{Conf: conf, DB: db}
}


// OneRequest handle one request
// record it into mongodb
func (s *ServeRequestLogger) OneRequest (r *RequestInfo) error {
	return s.DB.InsertOne(r)
}


// GetStat get data from mongodb
// and give back ip info
func (s *ServeRequestLogger) GetStat (start, end time.Time) (*RequestLoggerStat, error) {
	ris, err := s.DB.FindDuration(start, end)
	if err != nil {
		return nil, err
	}
	stat := new(RequestLoggerStat)
	stat.MethodCount = make(map[string]int64)
	stat.PathCount = make(map[string]int64)
	stat.ClientIPCount = make(map[string]int64)
	for _, r := range ris {
		if r.Method != "" {
			stat.MethodCount[r.Method] += 1
		}
		stat.PathCount[r.Path] += 1
		stat.ClientIPCount[r.ClientIP] += 1
	}
	return stat, nil
}


// ReqLogDB to mongodb
type ReqLogDB struct {
	Coll mgo.Collection
	Session *mgo.Session
}

// MgoConfig is copyed
type MgoConfig struct {
	Addrs []string `json:"Addrs"`
	Username string `json:"user"`
	Password string `json:"pass"`
	DB string `json:"db"`
	Coll string `json:"coll"`
}


// OpenReqLogDB opens new db
func OpenReqLogDB(conf MgoConfig) *ReqLogDB {
	conf.Coll = "reqLog"
	info := mgo.DialInfo{
		Addrs:    conf.Addrs,
		Username: conf.Username,
		Password: conf.Password,
		Database: conf.DB,
	}
	session, err := mgo.DialWithInfo(&info)
	if err != nil {
		return nil
	}
	db := session.DB(conf.DB)
	c := db.C(conf.Coll)
	return &ReqLogDB{
		Coll: *c,
		Session: session,
	}
}


// FindDuration find logs from start to end
func (db *ReqLogDB) FindDuration(start, end time.Time) ([]RequestInfo, error) {
	var ri []RequestInfo
	err := db.Coll.Find(bson.M{"time":bson.M{"$gt": start, "$lt": end}}).All(&ri)
	return ri, err
}


// InsertOne insert one log
func (db *ReqLogDB) InsertOne(ri *RequestInfo) error {
	if ri.Id == "" {
		ri.Id = bson.NewObjectId().Hex()
	}
	return db.Coll.Insert(ri)
}
