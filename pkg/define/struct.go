package define

type User struct {
	ID       int64  `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Email    string `json:"email" db:"email"`
	Traffic  int64  `json:"traffic" db:"traffic"`
	Group    string `json:"group" db:"group"`
	Regtime  int64  `json:"reg_time" db:"regtime"` // 假设注册时间应该是日期格式
	Status   bool   `json:"status" db:"status"`
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
	Status   bool   `json:"status" db:"status"`
	Token    string `json:"token" db:"token"`
	EmailMD5 string `gorm:"-" json:"email_md5"`
	Outbound int64  `json:"outbound"`
	Inbound  int    `json:"inbound"`
	Proxies  int
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

type Proxies struct {
	ID         int    `json:"id"`
	Username   string `json:"username"`
	ProxyName  string `json:"proxy_name"`
	ProxyType  string `json:"proxy_type"`
	LocalIP    string `json:"local_ip"`
	LocalPort  string `json:"local_port"`
	Domain     string `json:"domain"`
	Node       string `json:"node"`
	RemotePort string `json:"remote_port"`
	RunID      string
	Status     bool
	Lastupdate int64 `json:"lastupdate"`
}

// stop here
type Proxy struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	ProxyName    string `json:"proxy_name"`
	ProxyType    string `json:"proxy_type"`
	LocalIP      string `json:"local_ip"`
	LocalPort    string `json:"local_port"`
	Domain       string `json:"domain"`
	Node         string `json:"node"`
	NodeName     string `json:"node_name"`
	RunID        string
	NodeHostname string `json:"node_hostname"`
	NodePort     string `json:"node_port"`
	NodeToken    string `json:"node_token"`
	RemotePort   string `json:"remote_port"`
}

type Node struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Token     string `json:"token"`
	Group     string `json:"group"`
	AllowPort string `json:"allow_port"`
	AllowType string `json:"allow_type"`
	Status    int
	KumaId    int
	Health24  float64
}

type KumaData struct {
	HeartbeatList map[string][]struct {
		Status int
		Time   string
	}
	UptimeList map[string]float64 `json:"uptimeList"`
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
