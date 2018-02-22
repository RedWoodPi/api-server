package weather


import (
    "github.com/opesun/goquery"
    "fmt"
    "strings"
    "encoding/json"
    "bytes"
    client "api/redis"
    "github.com/garyburd/redigo/redis"
)
type CitySlice struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Forecasts []Forecast `json:"forecasts"`
}
type Forecast struct {
    Date string `json:"date"`
    Weather [2]string `json:"weather"`
    Temprature [2]string `json:"temprature"`
}
func Weather(name string) (strs string){
    url, city := city(name)
    if url == "" {
        return `{"id":404}`
    }
    //首先读取缓存
    rc := client.RedisClient.Get()
    defer rc.Close()
    result, err := redis.String(rc.Do("GET", city))
    if err != nil {
        fmt.Println(err)
    }
    if result != "" {
        return result
    }
    
    content, err := goquery.ParseUrl("http://www.nmc.cn"+url)
    if err != nil {
        fmt.Println("获取目标网站出错")
    }
    a := content.Find(".forecast")
    td := a.Find(".today")
    d := a.Find(".day")
    date := d.Find(".date")
    weath := td.Find(".wdesc")
    temp := td.Find(".temp")
    
    //对时间，天气，温度字符串进行清理
    dateArr := text2arr(date)
    weathArr := text2arr(weath)
    tempArr := text2arr(temp)
    var s CitySlice
    s.Id = 1
    s.Name = city
    //构造
    var forecasts [7]Forecast
    for i:=0; i<7; i++ {
        if i==0 {
            forecasts[i].Date = dateArr[i]
            forecasts[i].Weather[0] = weathArr[i]
            forecasts[i].Weather[1] = ""
            forecasts[i].Temprature[0] = tempArr[i]
            forecasts[i].Temprature[1] = ""
        } else {
            forecasts[i].Date = dateArr[i]
            forecasts[i].Weather[0] = weathArr[2*i-1]
            forecasts[i].Weather[1] = weathArr[2*i]
            forecasts[i].Temprature[0] = tempArr[2*i-1]
            forecasts[i].Temprature[1] = tempArr[2*i]
        }
        s.Forecasts = append(s.Forecasts, forecasts[i])
    }
    //转json
    b, err := json.Marshal(s)
    rc.Do("SET", city, string(b))
    rc.Do("EXPIRE", city, 2*3600)
    return string(b)
}

//取消空格和回车
func text2arr(a goquery.Nodes)(arr []string)  {
    
    for i:=0;i < 20; i++ {
        if a.Eq(i).Text() =="" {break}
        st := strings.Replace(a.Eq(i).Text(), " ","",-1)
        st = strings.Replace(st, "\n","",-1)
        arr = append(arr,st)
    }
    return arr
}
//城市链接
func city(name string)(urlstr string, arr string)  {
    //城市链接map
    citys := map[string]string{"巢湖":"/publish/forecast/AAH/chaohu.html", "杭州":"/publish/forecast/AZJ/hangzhou.html", "西安":"/publish/forecast/ASN/xian.html", "阳泉":"/publish/forecast/ASX/yangquan.html", "二连浩特":"/publish/forecast/ANM/erlianhaote.html", "漠河":"/publish/forecast/AHL/mohe.html", "铁岭":"/publish/forecast/ALN/tieling.html", "铜川":"/publish/forecast/ASN/tongchuan.html", "洞头":"/publish/forecast/AZJ/dongtou.html", "北京":"/publish/forecast/ABJ/beijing.html", "宁海":"/publish/forecast/AZJ/ninghai.html", "临高":"/publish/forecast/AHI/lingao.html", "滨海新区":"/publish/forecast/ATJ/binhaixinqu.html", "沈阳":"/publish/forecast/ALN/shenyang.html", "贵港":"/publish/forecast/AGX/guigang.html", "武都":"/publish/forecast/AGS/wudou.html", "内江":"/publish/forecast/ASC/neijiang.html", "莱芜":"/publish/forecast/ASD/laiwu.html", "吉首":"/publish/forecast/AHN/jishou.html", "临夏":"/publish/forecast/AGS/linxia.html", "酒泉":"/publish/forecast/AGS/jiuquan.html", "广安":"/publish/forecast/ASC/guangan.html", "东台":"/publish/forecast/AJS/dongtai.html", "六安":"/publish/forecast/AAH/luan.html", "潢川":"/publish/forecast/AHA/huangchuan.html", "保定":"/publish/forecast/AHE/baoding.html", "杜蒙":"/publish/forecast/AHL/dumeng.html", "泽当":"/publish/forecast/AXZ/zedang.html", "安康":"/publish/forecast/ASN/ankang.html", "天门":"/publish/forecast/AHB/tianmen.html", "淮南":"/publish/forecast/AAH/huainan.html", "博乐":"/publish/forecast/AXJ/bole.html", "自贡":"/publish/forecast/ASC/zigong.html", "眉山":"/publish/forecast/ASC/meishan.html", "驻马店市":"/publish/forecast/AHA/zhumadianshi.html", "贵阳":"/publish/forecast/AGZ/guiyang.html", "淄博":"/publish/forecast/ASD/zibo.html", "黎城":"/publish/forecast/ASX/licheng.html", "上饶":"/publish/forecast/AJX/shangrao.html", "许昌":"/publish/forecast/AHA/xuchang.html", "荆门":"/publish/forecast/AHB/jingmen.html", "宝鸡":"/publish/forecast/ASN/baoji.html", "海拉尔":"/publish/forecast/ANM/hailaer.html", "贵溪":"/publish/forecast/AJX/guixi.html", "深圳":"/publish/forecast/AGD/shenzhen.html", "渭南":"/publish/forecast/ASN/weinan.html", "澳门":"/publish/forecast/AAM/aomen.html", "唐山":"/publish/forecast/AHE/tangshan.html", "阿拉善右旗":"/publish/forecast/ANM/alashanyouqi.html", "瓦房店":"/publish/forecast/ALN/wafangdian.html", "保山":"/publish/forecast/AYN/baoshan.html", "来宾":"/publish/forecast/AGX/laibin.html", "克拉玛依":"/publish/forecast/AXJ/kelamayi.html", "娄底":"/publish/forecast/AHN/loudi.html", "汕尾":"/publish/forecast/AGD/shanwei.html", "拉萨":"/publish/forecast/AXZ/lasa.html", "衡水":"/publish/forecast/AHE/hengshui.html", "高邮":"/publish/forecast/AJS/gaoyou.html", "赣州":"/publish/forecast/AJX/ganzhou.html", "珠海":"/publish/forecast/AGD/zhuhai.html", "乐山":"/publish/forecast/ASC/leshan.html", "北海":"/publish/forecast/AGX/beihai.html", "儋州":"/publish/forecast/AHI/danzhou.html", "绵阳":"/publish/forecast/ASC/mianyang.html", "衢州":"/publish/forecast/AZJ/quzhou.html", "焦作市":"/publish/forecast/AHA/jiaozuoshi.html", "桦甸":"/publish/forecast/AJL/huadian.html", "刚察":"/publish/forecast/AQH/gangcha.html", "昆明":"/publish/forecast/AYN/kunming.html", "沧州":"/publish/forecast/AHE/cangzhou.html", "阳江":"/publish/forecast/AGD/yangjiang.html", "白沙":"/publish/forecast/AHI/baisha.html", "营口":"/publish/forecast/ALN/yingkou.html", "中沙":"/publish/forecast/AHI/zhongsha.html", "佛冈":"/publish/forecast/AGD/fogang.html", "崇左":"/publish/forecast/AGX/chongzuo.html", "三亚":"/publish/forecast/AHI/sanya.html", "临汾":"/publish/forecast/ASX/linfen.html", "陵水":"/publish/forecast/AHI/lingshui.html", "新宾":"/publish/forecast/ALN/xinbin.html", "吉安":"/publish/forecast/AJX/jian.html", "南阳市":"/publish/forecast/AHA/nanyangshi.html", "湛江":"/publish/forecast/AGD/zhanjiang.html", "包头":"/publish/forecast/ANM/baotou.html", "神农架":"/publish/forecast/AHB/shennongjia.html", "佛山":"/publish/forecast/AGD/foshan.html", "云浮":"/publish/forecast/AGD/yunfu.html", "五台山":"/publish/forecast/ASX/wutaishan.html", "芜湖":"/publish/forecast/AAH/wuhu.html", "宣城":"/publish/forecast/AAH/xuancheng.html", "无锡":"/publish/forecast/AJS/wuxi.html", "汉中":"/publish/forecast/ASN/hanzhong.html", "库尔勒":"/publish/forecast/AXJ/kuerle.html", "石家庄":"/publish/forecast/AHE/shijiazhuang.html", "文水":"/publish/forecast/ASX/wenshui.html", "商洛":"/publish/forecast/ASN/shangluo.html", "集宁":"/publish/forecast/ANM/jining.html", "荆州":"/publish/forecast/AHB/jingzhou.html", "平湖":"/publish/forecast/AZJ/pinghu.html", "衡阳市":"/publish/forecast/AHN/hengyangshi.html", "德令哈":"/publish/forecast/AQH/delingha.html", "乐东":"/publish/forecast/AHI/ledong.html", "钓鱼岛":"/publish/forecast/AFJ/diaoyudao.html", "厦门":"/publish/forecast/AFJ/xiamen.html", "泸县":"/publish/forecast/ASC/luxian.html", "中山":"/publish/forecast/AGD/zhongshan.html", "琼中":"/publish/forecast/AHI/qiongzhong.html", "顺德":"/publish/forecast/AGD/shunde.html", "兰州":"/publish/forecast/AGS/lanzhou.html", "镇江":"/publish/forecast/AJS/zhenjiang.html", "哈密":"/publish/forecast/AXJ/hami.html", "峨眉山":"/publish/forecast/ASC/emeishan.html", "临河":"/publish/forecast/ANM/linhe.html", "临海":"/publish/forecast/AZJ/linhai.html", "嘉峪关":"/publish/forecast/AGS/jiayuguan.html", "黔江":"/publish/forecast/ACQ/qianjiang.html", "安顺":"/publish/forecast/AGZ/anshun.html", "河源":"/publish/forecast/AGD/heyuan.html", "白银":"/publish/forecast/AGS/baiyin.html", "益阳":"/publish/forecast/AHN/yiyang.html", "清远":"/publish/forecast/AGD/qingyuan.html", "马鞍山":"/publish/forecast/AAH/maanshan.html", "潞西":"/publish/forecast/AYN/luxi.html", "南宁城区":"/publish/forecast/AGX/nanningchengqu.html", "贡山":"/publish/forecast/AYN/gongshan.html", "集安":"/publish/forecast/AJL/jian.html", "鸡西":"/publish/forecast/AHL/jixi.html", "普兰":"/publish/forecast/AXZ/pulan.html", "西沙永兴岛":"/publish/forecast/AHI/xishayongxingdao.html", "海口":"/publish/forecast/AHI/haikou.html", "都匀":"/publish/forecast/AGZ/douyun.html", "玉山":"/publish/forecast/AJX/yushan.html", "饶平":"/publish/forecast/AGD/raoping.html", "加格达奇":"/publish/forecast/AHL/jiagedaqi.html", "湖州":"/publish/forecast/AZJ/huzhou.html", "三明":"/publish/forecast/AFJ/sanming.html", "廊坊":"/publish/forecast/AHE/langfang.html", "黑河":"/publish/forecast/AHL/heihe.html", "徐州":"/publish/forecast/AJS/xuzhou.html", "奉节":"/publish/forecast/ACQ/fengjie.html", "思茅":"/publish/forecast/AYN/simao.html", "忻州":"/publish/forecast/ASX/xinzhou.html", "沅陵":"/publish/forecast/AHN/yuanling.html", "西宁":"/publish/forecast/AQH/xining.html", "雅安":"/publish/forecast/ASC/yaan.html", "郑州":"/publish/forecast/AHA/zhengzhou.html", "咸宁":"/publish/forecast/AHB/xianning.html", "南岳":"/publish/forecast/AHN/nanyue.html", "海晏":"/publish/forecast/AQH/haiyan.html", "大同":"/publish/forecast/ASX/datong.html", "大庆":"/publish/forecast/AHL/daqing.html", "庐山":"/publish/forecast/AJX/lushan.html", "随州":"/publish/forecast/AHB/suizhou.html", "甘孜":"/publish/forecast/ASC/ganzi.html", "南京":"/publish/forecast/AJS/nanjing.html", "淮北":"/publish/forecast/AAH/huaibei.html", "瑞丽":"/publish/forecast/AYN/ruili.html", "烟台":"/publish/forecast/ASD/yantai.html", "佳木斯":"/publish/forecast/AHL/jiamusi.html", "抚州":"/publish/forecast/AJX/fuzhou.html", "河南":"/publish/forecast/AQH/henan.html", "资阳":"/publish/forecast/ASC/ziyang.html", "泰州":"/publish/forecast/AJS/taizhou.html", "文山":"/publish/forecast/AYN/wenshan.html", "定安":"/publish/forecast/AHI/dingan.html", "广昌":"/publish/forecast/AJX/guangchang.html", "宜春":"/publish/forecast/AJX/yichun.html", "和田":"/publish/forecast/AXJ/hetian.html", "遂宁":"/publish/forecast/ASC/suining.html", "温州":"/publish/forecast/AZJ/wenzhou.html", "新余":"/publish/forecast/AJX/xinyu.html", "南通":"/publish/forecast/AJS/nantong.html", "日照":"/publish/forecast/ASD/rizhao.html", "青岛":"/publish/forecast/ASD/qingdao.html", "七台河":"/publish/forecast/AHL/qitaihe.html", "绥化":"/publish/forecast/AHL/suihua.html", "屯昌":"/publish/forecast/AHI/tunchang.html", "菏泽":"/publish/forecast/ASD/heze.html", "黄山市":"/publish/forecast/AAH/huangshanshi.html", "凯里":"/publish/forecast/AGZ/kaili.html", "临沂":"/publish/forecast/ASD/linyi.html", "高要":"/publish/forecast/AGD/gaoyao.html", "万宁":"/publish/forecast/AHI/wanning.html", "昌吉":"/publish/forecast/AXJ/changji.html", "梧州":"/publish/forecast/AGX/wuzhou.html", "宜昌":"/publish/forecast/AHB/yichang.html", "十堰":"/publish/forecast/AHB/shiyan.html", "崆峒":"/publish/forecast/AGS/kongtong.html", "乌鲁木齐":"/publish/forecast/AXJ/wulumuqi.html", "南昌":"/publish/forecast/AJX/nanchang.html", "潮州":"/publish/forecast/AGD/chaozhou.html", "林芝":"/publish/forecast/AXZ/linzhi.html", "朝阳":"/publish/forecast/ALN/chaoyang.html", "楚雄":"/publish/forecast/AYN/chuxiong.html", "吴忠":"/publish/forecast/ANX/wuzhong.html", "阿拉尔":"/publish/forecast/AXJ/alaer.html", "临沧":"/publish/forecast/AYN/lincang.html", "邢台":"/publish/forecast/AHE/xingtai.html", "阿克苏":"/publish/forecast/AXJ/akesu.html", "鄂州":"/publish/forecast/AHB/ezhou.html", "澄迈":"/publish/forecast/AHI/chengmai.html", "格尔木":"/publish/forecast/AQH/geermu.html", "大武口":"/publish/forecast/ANX/dawukou.html", "承德":"/publish/forecast/AHE/chengde.html", "赣榆":"/publish/forecast/AJS/ganyu.html", "孝感":"/publish/forecast/AHB/xiaogan.html", "武汉":"/publish/forecast/AHB/wuhan.html", "南雄":"/publish/forecast/AGD/nanxiong.html", "达州":"/publish/forecast/ASC/dazhou.html", "本溪":"/publish/forecast/ALN/benxi.html", "金华":"/publish/forecast/AZJ/jinhua.html", "莆田":"/publish/forecast/AFJ/putian.html", "临安":"/publish/forecast/AZJ/linan.html", "仙桃":"/publish/forecast/AHB/xiantao.html", "锦州":"/publish/forecast/ALN/jinzhou.html", "盘山":"/publish/forecast/ALN/panshan.html", "珲春":"/publish/forecast/AJL/hunchun.html", "滨州":"/publish/forecast/ASD/binzhou.html", "玛沁":"/publish/forecast/AQH/maqin.html", "铜仁":"/publish/forecast/AGZ/tongren.html", "锡林浩特":"/publish/forecast/ANM/xilinhaote.html", "吉林":"/publish/forecast/AJL/jilin.html", "安庆":"/publish/forecast/AAH/anqing.html", "塔城":"/publish/forecast/AXJ/tacheng.html", "台北":"/publish/forecast/ATW/taibei.html", "齐齐哈尔":"/publish/forecast/AHL/qiqihaer.html", "牡丹江":"/publish/forecast/AHL/mudanjiang.html", "永州":"/publish/forecast/AHN/yongzhou.html", "那曲":"/publish/forecast/AXZ/naqu.html", "景洪":"/publish/forecast/AYN/jinghong.html", "长沙":"/publish/forecast/AHN/changsha.html", "长治":"/publish/forecast/ASX/changzhi.html", "盐城":"/publish/forecast/AJS/yancheng.html", "漯河市":"/publish/forecast/AHA/luoheshi.html", "白山":"/publish/forecast/AJL/baishan.html", "日喀则":"/publish/forecast/AXZ/rikaze.html", "浦城":"/publish/forecast/AFJ/pucheng.html", "文昌":"/publish/forecast/AHI/wenchang.html", "济南":"/publish/forecast/ASD/jinan.html", "德州":"/publish/forecast/ASD/dezhou.html", "济源":"/publish/forecast/AHA/jiyuan.html", "阿拉善左旗":"/publish/forecast/ANM/alashanzuoqi.html", "辽源":"/publish/forecast/AJL/liaoyuan.html", "蚌埠":"/publish/forecast/AAH/bengbu.html", "盱眙":"/publish/forecast/AJS/xuyi.html", "韶关":"/publish/forecast/AGD/shaoguan.html", "邯郸":"/publish/forecast/AHE/handan.html", "保亭":"/publish/forecast/AHI/baoting.html", "武威":"/publish/forecast/AGS/wuwei.html", "桑植":"/publish/forecast/AHN/sangzhi.html", "苏州":"/publish/forecast/AJS/suzhou.html", "枣庄":"/publish/forecast/ASD/zaozhuang.html", "黄石":"/publish/forecast/AHB/huangshi.html", "中卫":"/publish/forecast/ANX/zhongwei.html", "四平":"/publish/forecast/AJL/siping.html", "福州":"/publish/forecast/AFJ/fuzhou.html", "襄阳":"/publish/forecast/AHB/xiangyang.html", "郴州":"/publish/forecast/AHN/chenzhou.html", "平安":"/publish/forecast/AQH/pingan.html", "固原":"/publish/forecast/ANX/guyuan.html", "宿州":"/publish/forecast/AAH/suzhou.html", "哈尔滨":"/publish/forecast/AHL/haerbin.html", "铜陵":"/publish/forecast/AAH/tongling.html", "太原":"/publish/forecast/ASX/taiyuan.html", "常德":"/publish/forecast/AHN/changde.html", "兴义":"/publish/forecast/AGZ/xingyi.html", "楚州":"/publish/forecast/AJS/chuzhou.html", "丽水":"/publish/forecast/AZJ/lishui.html", "钦州":"/publish/forecast/AGX/qinzhou.html", "通辽":"/publish/forecast/ANM/tongliao.html", "榆林":"/publish/forecast/ASN/yulin.html", "共和":"/publish/forecast/AQH/gonghe.html", "重庆":"/publish/forecast/ACQ/chongqing.html", "萍乡":"/publish/forecast/AJX/pingxiang.html", "玉树":"/publish/forecast/AQH/yushu.html", "丹东":"/publish/forecast/ALN/dandong.html", "常州":"/publish/forecast/AJS/changzhou.html", "宁波":"/publish/forecast/AZJ/ningbo.html", "石河子":"/publish/forecast/AXJ/shihezi.html", "昭通":"/publish/forecast/AYN/zhaotong.html", "延吉":"/publish/forecast/AJL/yanji.html", "百色":"/publish/forecast/AGX/baise.html", "玉溪":"/publish/forecast/AYN/yuxi.html", "东莞":"/publish/forecast/AGD/dongguan.html", "茂名":"/publish/forecast/AGD/maoming.html", "天水":"/publish/forecast/AGS/tianshui.html", "涪陵":"/publish/forecast/ACQ/fuling.html", "遵义":"/publish/forecast/AGZ/zunyi.html", "张掖":"/publish/forecast/AGS/zhangye.html", "梅县":"/publish/forecast/AGD/meixian.html", "防城港":"/publish/forecast/AGX/fangchenggang.html", "南平":"/publish/forecast/AFJ/nanping.html", "黄冈":"/publish/forecast/AHB/huanggang.html", "红河":"/publish/forecast/AYN/honghe.html", "景德镇":"/publish/forecast/AJX/jingdezhen.html", "济宁":"/publish/forecast/ASD/jining.html", "汕头":"/publish/forecast/AGD/shantou.html", "德阳":"/publish/forecast/ASC/deyang.html", "葫芦岛":"/publish/forecast/ALN/huludao.html", "惠州":"/publish/forecast/AGD/huizhou.html", "勃利":"/publish/forecast/AHL/boli.html", "阿勒泰":"/publish/forecast/AXJ/aletai.html", "连云港":"/publish/forecast/AJS/lianyungang.html", "九江":"/publish/forecast/AJX/jiujiang.html", "聊城":"/publish/forecast/ASD/liaocheng.html", "喀什":"/publish/forecast/AXJ/kashen.html", "宜宾":"/publish/forecast/ASC/yibin.html", "潜江":"/publish/forecast/AHB/qianjiang.html", "东胜":"/publish/forecast/ANM/dongsheng.html", "鞍山":"/publish/forecast/ALN/anshan.html", "信阳市":"/publish/forecast/AHA/xinyangshi.html", "湘潭":"/publish/forecast/AHN/xiangtan.html", "邵阳市":"/publish/forecast/AHN/shaoyangshi.html", "庆城":"/publish/forecast/AGS/qingcheng.html", "奇台":"/publish/forecast/AXJ/qitai.html", "通化":"/publish/forecast/AJL/tonghua.html", "亳州":"/publish/forecast/AAH/bozhou.html", "潍坊":"/publish/forecast/ASD/weifang.html", "伊宁市":"/publish/forecast/AXJ/yiningshi.html", "黄岩岛":"/publish/forecast/AHI/huangyandao.html", "濮阳":"/publish/forecast/AHA/puyangxian.html", "桂林":"/publish/forecast/AGX/guilin.html", "伊春":"/publish/forecast/AHL/yichun.html", "绍兴":"/publish/forecast/AZJ/shaoxing.html", "双鸭山":"/publish/forecast/AHL/shuangyashan.html", "滁州":"/publish/forecast/AAH/chuzhou.html", "恩施":"/publish/forecast/AHB/enshi.html", "延安":"/publish/forecast/ASN/yanan.html", "朔州":"/publish/forecast/ASX/shuozhou.html", "合肥":"/publish/forecast/AAH/hefei.html", "阿坝":"/publish/forecast/ASC/aba.html", "临武":"/publish/forecast/AHN/linwu.html", "柳州":"/publish/forecast/AGX/liuzhou.html", "大连":"/publish/forecast/ALN/dalian.html", "阜新":"/publish/forecast/ALN/fuxin.html", "商丘":"/publish/forecast/AHA/shangqiu.html", "琼海":"/publish/forecast/AHI/qionghai.html", "西昌":"/publish/forecast/ASC/xichang.html", "任丘":"/publish/forecast/AHE/renqiu.html", "龙岩":"/publish/forecast/AFJ/longyan.html", "宜城":"/publish/forecast/AHB/yicheng.html", "吕梁":"/publish/forecast/ASX/lvliang.html", "绥芬河":"/publish/forecast/AHL/suifenhe.html", "昌都":"/publish/forecast/AXZ/changdou.html", "鹤岗":"/publish/forecast/AHL/hegang.html", "定海":"/publish/forecast/AZJ/dinghai.html", "岳阳":"/publish/forecast/AHN/yueyang.html", "广州":"/publish/forecast/AGD/guangzhou.html", "巴中":"/publish/forecast/ASC/bazhong.html", "丽江":"/publish/forecast/AYN/lijiang.html", "阿图什":"/publish/forecast/AXJ/atushen.html", "银川":"/publish/forecast/ANX/yinchuan.html", "南充":"/publish/forecast/ASC/nanchong.html", "张家口":"/publish/forecast/AHE/zhangjiakou.html", "晋江":"/publish/forecast/AFJ/jinjiang.html", "沙湾":"/publish/forecast/AXJ/shawan.html", "鹤山":"/publish/forecast/AGD/heshan.html", "攀枝花":"/publish/forecast/ASC/panzhihua.html", "怀化":"/publish/forecast/AHN/huaihua.html", "洛阳市":"/publish/forecast/AHA/luoyangshi.html", "张家界":"/publish/forecast/AHN/zhangjiajie.html", "玉林":"/publish/forecast/AGX/yulin.html", "晋城":"/publish/forecast/ASX/jincheng.html", "赤峰":"/publish/forecast/ANM/chifeng.html", "宁德":"/publish/forecast/AFJ/ningde.html", "河池":"/publish/forecast/AGX/hechi.html", "鹤壁":"/publish/forecast/AHA/hebi.html", "三门峡":"/publish/forecast/AHA/sanmenxia.html", "安定":"/publish/forecast/AGS/anding.html", "运城":"/publish/forecast/ASX/yuncheng.html", "新乡市":"/publish/forecast/AHA/xinxiangshi.html", "吐鲁番":"/publish/forecast/AXJ/tulufan.html", "广元":"/publish/forecast/ASC/guangyuan.html", "漳州":"/publish/forecast/AFJ/zhangzhou.html", "威海":"/publish/forecast/ASD/weihai.html", "毕节":"/publish/forecast/AGZ/bijie.html", "宿迁":"/publish/forecast/AJS/suqian.html", "周口":"/publish/forecast/AHA/zhoukou.html", "电白":"/publish/forecast/AGD/dianbai.html", "鹰潭":"/publish/forecast/AJX/yingtan.html", "安阳市":"/publish/forecast/AHA/anyangshi.html", "贺州":"/publish/forecast/AGX/hezhou.html", "平遥":"/publish/forecast/ASX/pingyao.html", "株洲":"/publish/forecast/AHN/zhuzhou.html", "秦皇岛":"/publish/forecast/AHE/qinhuangdao.html", "扬州":"/publish/forecast/AJS/yangzhou.html", "辽阳":"/publish/forecast/ALN/liaoyang.html", "开封市":"/publish/forecast/AHA/kaifengshi.html", "香港":"/publish/forecast/AXG/xianggang.html", "乌海":"/publish/forecast/ANM/wuhai.html", "长春":"/publish/forecast/AJL/changchun.html", "水城":"/publish/forecast/AGZ/shuicheng.html", "嵊州":"/publish/forecast/AZJ/shengzhou.html", "阜阳":"/publish/forecast/AAH/fuyang.html", "平顶山":"/publish/forecast/AHA/pingdingshan.html", "嘉兴":"/publish/forecast/AZJ/jiaxing.html", "池州":"/publish/forecast/AAH/chizhou.html", "甘南":"/publish/forecast/AHL/gannan.html", "东营":"/publish/forecast/ASD/dongying.html", "枣阳":"/publish/forecast/AHB/zaoyang.html", "西青":"/publish/forecast/ATJ/xiqing.html", "白城":"/publish/forecast/AJL/baicheng.html", "泰安":"/publish/forecast/ASD/taian.html", "黔西":"/publish/forecast/AGZ/qianxi.html", "天津":"/publish/forecast/ATJ/tianjin.html", "巴林右旗":"/publish/forecast/ANM/balinyouqi.html", "梅河口":"/publish/forecast/AJL/meihekou.html", "咸阳":"/publish/forecast/ASN/xianyang.html", "揭阳":"/publish/forecast/AGD/jieyang.html", "田阳":"/publish/forecast/AGX/tianyang.html", "呼和浩特":"/publish/forecast/ANM/huhehaote.html", "乌兰浩特":"/publish/forecast/ANM/wulanhaote.html", "麻城":"/publish/forecast/AHB/macheng.html", "东方":"/publish/forecast/AHI/dongfang.html", "成都":"/publish/forecast/ASC/chengdu.html", "大理":"/publish/forecast/AYN/dali.html", "松原":"/publish/forecast/AJL/songyuan.html", "上海":"/publish/forecast/ASH/shanghai.html", "香格里拉":"/publish/forecast/AYN/xianggelila.html", "泸州":"/publish/forecast/ASC/luzhou.html",}
    url := citys[name]
    if url == "" {
        if strings.Contains(name, "市") {
            name = strings.Replace(name, "市", "", -1)
        } else {
            var buf bytes.Buffer
            buf.WriteString(name)
            buf.WriteString("市")
            name = buf.String()
        }
    }
    return citys[name], name
}