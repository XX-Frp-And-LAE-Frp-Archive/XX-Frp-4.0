package define

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Traffic  int64  `json:"traffic" db:"traffic"`
	Group    string `json:"group" db:"group"`
	Regtime  int64  `json:"reg_time" db:"regtime"` // 假设注册时间应该是日期格式
	Status   int    `json:"status" db:"status"`
	Token    string `json:"token" db:"token"`
}

// 定义用户返回结构体
type UserRes struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Traffic  int64  `json:"traffic" db:"traffic"`
	Group    string `json:"group" db:"group"`
	RegTime  int64  `json:"reg_time" db:"regtime"` // 假设注册时间应该是日期格式
	Status   int    `json:"status" db:"status"`
	Token    string `json:"token" db:"token"`
	EmailMD5 string `gorm:"-" json:"email_md5"`
	Outbound int64  `json:"outbound"`
	Inbound  int    `json:"inbound"`
	Proxies  int    `json:"proxies"`
}

type Client struct {
	ID      int    `json:"id"`
	Version string `json:"version"`
	Log     string `json:"log"`
	Url     string `json:"url"`
}
type Config struct {
	Server struct {
		Host  string
		Port  int
		Name  string
		Url   string
		Token string
	}
	Debug struct {
		Debug bool
	}
	Mysql struct {
		User     string
		Password string
		Host     string
		Port     string
		Database string
	}
	Smtp struct {
		Addr   string
		Passwd string
		Port   int
		From   string
	}
	Realname struct {
		SecretID  string
		SecretKey string
	}
}

type Setting struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Sponsor struct {
	ID      int    `gorm:"primaryKey" json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Thing   string `json:"thing"`
	Comment string `json:"comment"`
}

type TodayTraffic struct {
	User    string `gorm:"column:user"`
	Traffic int64  `gorm:"column:traffic"`
}

type FrpsTraffic struct {
	Name            string `json:"name"`
	TodayTrafficIn  int64  `json:"today_traffic_in"`
	TodayTrafficOut int64  `json:"today_traffic_out"`
	CurConns        int    `json:"cur_conns"`
	LastStartTime   string `json:"last_start_time"`
	LastCloseTime   string `json:"last_close_time"`
	Status          string `json:"status"`
	ClientVersion   string `json:"client_version,omitempty"` // omitempty 表示如果该字段为空，则在序列化时忽略该字段
}
type FrpsTrafficData struct {
	Proxies []FrpsTraffic `json:"proxies"`
}

type Proxies struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	ProxyName       string `json:"proxy_name"`
	ProxyType       string `json:"proxy_type"`
	LocalIP         string `json:"local_ip"`
	LocalPort       string `json:"local_port"`
	Domain          string `json:"domain"`
	Node            string `json:"node"`
	RemotePort      string `json:"remote_port"`
	RunID           string `json:"run_id"`
	Status          int    `json:"status"`
	Online          string `json:"online"`
	TodayTrafficIn  int64  `json:"today_traffic_in"`
	CurConns        int    `json:"cur_conns"`
	TodayTrafficOut int64  `json:"today_traffic_out"`
	LastStartTime   string `json:"last_start_time"`
	LastCloseTime   string `json:"last_close_time"`
	Lastupdate      int64  `json:"lastupdate"`
	ClientVersion   string `json:"client_version,omitempty"`
}

// stop here
type Proxy struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	ProxyName       string `json:"proxy_name"`
	ProxyType       string `json:"proxy_type"`
	LocalIP         string `json:"local_ip"`
	LocalPort       string `json:"local_port"`
	Domain          string `json:"domain"`
	Node            string `json:"node"`
	NodeName        string `json:"node_name"`
	RunID           string `json:"run_id"`
	NodeHostname    string `json:"node_hostname"`
	NodePort        string `json:"node_port"`
	NodeToken       string `json:"node_token"`
	Online          string `json:"online"`
	CurConns        int    `json:"cur_conns"`
	TodayTrafficIn  int64  `json:"today_traffic_in"`
	TodayTrafficOut int64  `json:"today_traffic_out"`
	RemotePort      string `json:"remote_port"`
	LastStartTime   string `json:"last_start_time"`
	LastCloseTime   string `json:"last_close_time"`
	Status          int    `json:"status"`
	ClientVersion   string `json:"client_version,omitempty"`
}

type Node struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	Description     string `json:"description"`
	Port            string `json:"port"`
	Token           string `json:"token"`
	Group           string `json:"group"`
	AllowPort       string `json:"allow_port"`
	AdminPort       string `json:"admin_port"`
	AdminPass       string `json:"admin_pass"`
	AllowType       string `json:"allow_type"`
	Status          int    `json:"status"`
	TotalTrafficIn  int64  `json:"total_traffic_in"`
	TotalTrafficOut int64  `json:"total_traffic_out"`
	OnlineCount     int64  `json:"online_count"`
	Version         string `json:"version"`
}

type Nodes struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	Description     string `json:"description"`
	Port            string `json:"port"`
	Token           string `json:"token"`
	Group           string `json:"group"`
	AllowPort       string `json:"allow_port"`
	AllowType       string `json:"allow_type"`
	Status          int    `json:"status"`
	TotalTrafficIn  int64  `json:"total_traffic_in"`
	TotalTrafficOut int64  `json:"total_traffic_out"`
	OnlineCount     int64  `json:"online_count"`
	Version         string `json:"version"`
}

type FrpsData struct {
	TotalTrafficIn  int64  `json:"total_traffic_in"`
	TotalTrafficOut int64  `json:"total_traffic_out"`
	Version         string `json:"version"`
	ProxyTypeCount  struct {
		Tcp   int `json:"tcp"`
		Udp   int `json:"udp"`
		Http  int `json:"http"`
		Https int `json:"https"`
	} `json:"proxy_type_count"`
}

type Code struct {
	Email string
	Code  string
	Time  int64
}

type Groups struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `gorm:"varchar(255); not null" json:"name"`
	Traffic    int64
	Outbound   int    `json:"outbound"`
	Inbound    int    `json:"inbound"`
	Proxies    int    `json:"proxies"`
	CreateTime string `gorm:"varchar(255); not null" json:"create_time"`
}

type Limit struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Inbound  int    `json:"inbound"`
	Outbound int    `json:"outbound"`
	Proxies  int    `json:"proxies"`
}

type Sign struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Username     string `gorm:"varchar(255); not null" json:"username"`
	Signdate     int64  `gorm:"varchar(255); not null" json:"signdate"`
	Totalsign    int64  `json:"totalsign"`
	Totaltraffic int64  `json:"totaltraffic"`
}

type Findpass struct {
	Username string `gorm:"primaryKey"`
	Link     string
	Time     int64
}
