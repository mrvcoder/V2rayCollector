package main

import (
	"fmt"
	"os"

	"github.com/mrvcoder/V2rayCollector/collector"
)

// channels
var defaultChannelsConfig = []string{
	// best
	"https://t.me/s/v2rayng_org",
	"https://t.me/s/v2rayngvpn",
	"https://t.me/s/flyv2ray",
	"https://t.me/s/v2ray_outlineir",
	"https://t.me/s/v2_vmess",
	"https://t.me/s/FreeVlessVpn",
	"https://t.me/s/freeland8",
	"https://t.me/s/vmess_vless_v2rayng",
	"https://t.me/s/PrivateVPNs",
	"https://t.me/s/vmessiran",
	"https://t.me/s/Outline_Vpn",
	"https://t.me/s/vmessq",
	"https://t.me/s/WeePeeN",
	"https://t.me/s/V2rayNG3",
	"https://t.me/s/ShadowsocksM",
	"https://t.me/s/shadowsocksshop",
	"https://t.me/s/v2rayan",
	"https://t.me/s/ShadowSocks_s",
	"https://t.me/s/VmessProtocol",
	"https://t.me/s/napsternetv_config",
	"https://t.me/s/Easy_Free_VPN",
	"https://t.me/s/V2Ray_FreedomIran",
	"https://t.me/s/V2RAY_VMESS_free",
	"https://t.me/s/v2ray_for_free",
	"https://t.me/s/V2rayN_Free",
	"https://t.me/s/free4allVPN",
	"https://t.me/s/vpn_ocean",
	"https://t.me/s/configV2rayForFree",
	"https://t.me/s/FreeV2rays{all_messages}",
	"https://t.me/s/DigiV2ray",
	"https://t.me/s/v2rayNG_VPN",
	"https://t.me/s/freev2rayssr",
	"https://t.me/s/v2rayn_server",
	"https://t.me/s/Shadowlinkserverr",
	"https://t.me/s/iranvpnet",
	"https://t.me/s/vmess_iran",
	"https://t.me/s/mahsaamoon1",
	"https://t.me/s/V2RAY_NEW",
	"https://t.me/s/v2RayChannel",
	"https://t.me/s/configV2rayNG{all_messages}",
	"https://t.me/s/config_v2ray",
	"https://t.me/s/vpn_proxy_custom",
	"https://t.me/s/vpnmasi{all_messages}",
	"https://t.me/s/v2ray_custom",
	"https://t.me/s/VPNCUSTOMIZE",
	"https://t.me/s/HTTPCustomLand",
	"https://t.me/s/vpn_proxy_custom",
	"https://t.me/s/ViPVpn_v2ray",
	"https://t.me/s/FreeNet1500",
	"https://t.me/s/v2ray_ar{all_messages}",
	"https://t.me/s/beta_v2ray",
	"https://t.me/s/vip_vpn_2022",
	"https://t.me/s/FOX_VPN66",
	"https://t.me/s/VorTexIRN",
	"https://t.me/s/YtTe3la",
	"https://t.me/s/V2RayOxygen",
	"https://t.me/s/Network_442",
	"https://t.me/s/VPN_443",
	"https://t.me/s/v2rayng_v",
	"https://t.me/s/ultrasurf_12",
	"https://t.me/s/iSeqaro{all_messages}",
	"https://t.me/s/frev2rayng",
	"https://t.me/s/frev2ray",
	"https://t.me/s/FreakConfig",
	"https://t.me/s/Awlix_ir",
	"https://t.me/s/v2rayngvpn",
	"https://t.me/s/God_CONFIG{all_messages}",
	"https://t.me/s/Configforvpn01",
	"https://t.me/s/polproxy",
	"https://t.me/s/v2rayvpnchannel",
	"https://t.me/s/proxy_mtm",
	"https://t.me/s/vpn_ioss{all_messages}",
	"https://t.me/s/V2Ray_FreedomIran",
	"https://t.me/s/v2rayfree1",
	"https://t.me/s/free_v2rayyy",
	"https://t.me/s/nx_v2ray",
	"https://t.me/s/nufilter",
	"https://t.me/s/Free_HTTPCustom",
	"https://t.me/s/customv2ray",
	"https://t.me/s/vpn_Nv1",
	"https://t.me/s/AliAlma_GSM{all_messages}",
	"https://t.me/s/reality_daily{all_messages}",
	"https://t.me/s/shopingv2ray",
	"https://t.me/s/v2rayng_vpnrog",
	"https://t.me/s/ServerNett",
	"https://t.me/s/MT_TEAM_IRAN",
	"https://t.me/s/V2ray_Team",
	"https://t.me/s/VpnProsecc",
	"https://t.me/s/ConfigsHUB",
	"https://t.me/s/melov2ray",
	"https://t.me/s/V2pedia",
	"https://t.me/s/FalconPolV2rayNG",
	"https://t.me/s/ShadowProxy66",
	"https://t.me/s/VPNCUSTOMIZE",
	"https://t.me/s/prrofile_purple",
	"https://t.me/s/MsV2ray",
	"https://t.me/s/VlessConfig",
	"https://t.me/s/vless_vmess",
	"https://t.me/s/MehradLearn",
	"https://t.me/s/kingofilter",
	"https://t.me/s/IRN_VPN",
	"https://t.me/s/V2raysFree",
	"https://t.me/s/SvnTeam",
	"https://t.me/s/flyv2ray",
	"https://t.me/s/free1_vpn",
	"https://t.me/s/UnlimitedDev",
	"https://t.me/s/vpn_xw",
	"https://t.me/s/V2RayTz",
	"https://t.me/s/ipV2Ray",
	"https://t.me/s/OutlineVpnOfficial",
	"https://t.me/s/mehrosaboran",
	"https://t.me/s/mftizi",
	"https://t.me/s/https_config_injector",
	"https://t.me/s/Hope_Net",
	"https://t.me/s/V2rayng_Fast",
	"https://t.me/s/DailyV2RY",
	"https://t.me/s/shh_proxy",
	"https://t.me/s/forwardv2ray",
	"https://t.me/s/Lockey_vpn",

	// test
	"https://t.me/s/oneclickvpnkeys",
	"https://t.me/s/DirectVPN",
	"https://t.me/s/Parsashonam",
	"https://t.me/s/V2rayNGmat",
	"https://t.me/s/fnet00",
	"https://t.me/s/Outline_Vpn",
	"https://t.me/s/azadi_az_inja_migzare",
	"https://t.me/s/vmess_vless_v2rayng",
	"https://t.me/s/v2ray_vpn_ir",
	"https://t.me/s/BestV2rang",
	"https://t.me/s/v2logy",
	"https://t.me/s/Awlix_V2RAY",
	"https://t.me/s/reality_daily",
	"https://t.me/s/DeamNet_Proxy",
	"https://t.me/s/vpn_go67",
	"https://t.me/s/SafeNet_Server",
	"https://t.me/s/mahvarehnewssat",
	"https://t.me/s/sinabigo",
	"https://t.me/s/FAKEOFTVC",
	"https://t.me/s/customv2ray",
	"https://t.me/s/VPNCUSTOMIZE",
	"https://t.me/s/rxv2ray",
	"https://t.me/s/MTConfig",
	"https://t.me/s/MehradLearn",
	"https://t.me/s/v2Line",
	"https://t.me/s/v2rayNG_Matsuri",
	"https://t.me/s/Helix_Servers",
	"https://t.me/s/EUT_VPN",
	"https://t.me/s/proxy_mtm",
	"https://t.me/s/eliya_chiter0",
	"https://t.me/s/melov2ray",
	"https://t.me/s/servermomo",
	"https://t.me/s/ARv2ray",
	"https://t.me/s/vless_vmess",
	"https://t.me/s/GozargahVPN",
	"https://t.me/s/V2rayCollector",
	"https://t.me/s/v2rayng_config_amin",
	"https://t.me/s/VPNCLOP",
	"https://t.me/s/DarkTeam_VPN",
	"https://t.me/s/ProxyForOpeta",
	"https://t.me/s/Outlinev2rayNG",
	"https://t.me/s/v2_team",
	"https://t.me/s/VpnFreeSec",
	"https://t.me/s/freeconfigv2",
	"https://t.me/s/hcv2ray",
	"https://t.me/s/NIM_VPN_ir",
	"https://t.me/s/Capital_NET",
	"https://t.me/s/v2ray_swhil",
	"https://t.me/s/XsV2ray",
	"https://t.me/s/V2parsin",
	"https://t.me/s/EliV2ray",
	"https://t.me/s/proxyymeliii",
	"https://t.me/s/V2rayCollectorDonate",
	"https://t.me/s/DigiV2ray",
	"https://t.me/s/free_v2rayyy",
	"https://t.me/s/yaney_01",
	"https://t.me/s/ShadowProxy66",
	"https://t.me/s/MrV2Ray",
	"https://t.me/s/V2rayNGn",
	"https://t.me/s/V2pedia",
	"https://t.me/s/FalconPolV2rayNG",
	"https://t.me/s/CUSTOMVPNSERVER",
	"https://t.me/s/MsV2ray",
	"https://t.me/s/ServerNett",
	"https://t.me/s/lrnbymaa",
	"https://t.me/s/Proxy_PJ",
	"https://t.me/s/vmessorg",
	"https://t.me/s/polproxy",
	"https://t.me/s/v2rayng_vpnrog",
	"https://t.me/s/lightning6",
	"https://t.me/s/frev2ray",
	"https://t.me/s/proxy_kafee",
	"https://t.me/s/Qv2rayDONATED",
	"https://t.me/s/Capoit",
	"https://t.me/s/PrivateVPNs",
	"https://t.me/s/Cov2ray",
	"https://t.me/s/kiava",
	"https://t.me/s/prrofile_purple",
	"https://t.me/s/nofiltering2",
}

func main() {
	configs := collector.GetConfigs(defaultChannelsConfig)
	WriteToFile(configs, "iran.txt")
}

func WriteToFile(fileContent string, filePath string) {

	// Check if the file exists
	if _, err := os.Stat(filePath); err == nil {
		// If the file exists, clear its content
		err = os.WriteFile(filePath, []byte{}, 0644)
		if err != nil {
			fmt.Println("Error clearing file:", err)
			return
		}
	} else if os.IsNotExist(err) {
		// If the file does not exist, create it
		_, err = os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	} else {
		// If there was some other error, print it and return
		fmt.Println("Error checking file:", err)
		return
	}

	// Write the new content to the file
	err := os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("File written successfully -->", filePath)
}

// func reverse(lines []string) []string {
// 	for i := 0; i < len(lines)/2; i++ {
// 		j := len(lines) - i - 1
// 		lines[i], lines[j] = lines[j], lines[i]
// 	}
// 	return lines
// }
