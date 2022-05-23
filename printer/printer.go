package printer

import (
	"fmt"
	"strings"

	"github.com/xgadget-lab/nexttrace/trace"

	"github.com/xgadget-lab/nexttrace/ipgeo"
)

var dataOrigin string

func TraceroutePrinter(res *trace.Result) {
	for i, hop := range res.Hops {
		fmt.Print(i + 1)
		for _, h := range hop {
			HopPrinter(h)
		}
	}
}

func HopPrinter(h trace.Hop) {
	if h.Address == nil {
		fmt.Println("\t*")
	} else {
		txt := "\t"

		if h.Hostname == "" {
			txt += fmt.Sprint(h.Address, " ", fmt.Sprintf("%.2f", h.RTT.Seconds()*1000), "ms")
		} else {
			txt += fmt.Sprint(h.Hostname, " (", h.Address, ") ", fmt.Sprintf("%.2f", h.RTT.Seconds()*1000), "ms")
		}

		if h.Geo != nil {
			txt += " " + formatIpGeoData(h.Address.String(), h.Geo)
		}

		fmt.Println(txt)
	}
}

func formatIpGeoData(ip string, data *ipgeo.IPGeoData) string {
	var res = make([]string, 0, 10)

	if data.Asnumber == "" {
		res = append(res, "*")
	} else {
		res = append(res, "AS"+data.Asnumber)
	}

	// TODO: 判断阿里云和腾讯云内网，数据不足，有待进一步完善
	// TODO: 移动IDC判断到Hop.fetchIPData函数，减少API调用
	if strings.HasPrefix(ip, "9.") {
		res = append(res, "局域网", "腾讯云")
	} else if strings.HasPrefix(ip, "11.") {
		res = append(res, "局域网", "阿里云")
	} else if data.Country == "" {
		res = append(res, "局域网")
	} else {
		// 有些IP的归属信息为空，这个时候将ISP的信息填入
		if data.Owner == "" {
			data.Owner = data.Isp
		}
		if data.District != "" {
			data.City = data.City + ", " + data.District
		}
		if data.Prov == "" && data.City == "" {
			// anyCast或是骨干网数据不应该有国家信息
			data.Owner = data.Owner + ", " + data.Owner
		} else {
			// 非骨干网正常填入IP的国家信息数据
			res = append(res, data.Country)
		}

		if data.Prov != "" {
			res = append(res, data.Prov)
		}
		if data.City != "" {
			res = append(res, data.City)
		}

		if data.Owner != "" {
			res = append(res, data.Owner)
		}
	}

	return strings.Join(res, ", ")
}