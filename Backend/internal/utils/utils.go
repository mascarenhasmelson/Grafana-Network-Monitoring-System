package utils

type Conf struct {
	Host string
	Port string
    Ifaces  string
	ArpStrs []string
	Timeout int
	HistInDB bool
	TrimHist int
}

// Host - one host
type Host struct {
	ID    int    `db:"ID"`
	Name  string `db:"NAME"`
	DNS   string `db:"DNS"`
	IP    string `db:"IP"`
	Iface string `db:"IFACE"`
	Mac   string `db:"MAC"`
	Hw    string `db:"HW"`
	Date  string `db:"DATE"`
	Known int    `db:"KNOWN"`
	Now   int    `db:"NOW"`
}
