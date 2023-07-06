# python爬虫
import random
import re
import requests
from bs4 import BeautifulSoup
import time
'''
 in : 一级url
 out : 二级url https://www.gsmchoice.com/zh-cn/catalogue/nec/
'''
def craw_lev1(base_url, url):
    li = []
    req_headers= dict()
    user_agent_list = ["Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
                       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
                       "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
                       "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
                       "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36",
                       "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)",
                       "Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10.5; en-US; rv:1.9.2.15) Gecko/20110303 Firefox/3.6.15",
                       ]
    req_headers['User-Agent'] = random.choice(user_agent_list)
    req_obj = requests.get(url,headers=req_headers)
    bresp = BeautifulSoup(req_obj.text,'lxml')
    CatalogueBrands = bresp.find(id='CatalogueBrands')
    a = CatalogueBrands.find_all('a')
    for item in a:
        if ("https" in item['href']):
            # 确认没有重复的框架内href没有重复的,一层url不去重直接追加
            li.append(item['href'])
        else:
            li.append(base_url + item['href'])
    return li

'''
 in : 二级url
 out : 三级url https://www.gsmchoice.com/zh-cn/catalogue/nec/mediasxn06e/
'''
def craw_lev2(url):
    soup_a = []
    base_url3 = []
    base_url = "https://www.gsmchoice.com/"

    factory = url.split('/')[-3]
    reg_key = 'href="/zh-cn/catalogue/' + factory + '/\w*'

    req_obj = requests.get(url)
    soup = BeautifulSoup(req_obj.text,'html.parser')
    soup_len = len(soup.find_all('div',class_='phone-container phone-container--left'))
    if soup_len == 2:
        soup_a = soup.find_all('div',class_='phone-container phone-container--left')[0].find_all('a')+soup.find_all('div',class_='phone-container phone-container--left')[1].find_all('a')
    else:
        soup_a = soup.find_all('div',class_='phone-container phone-container--left')[0].find_all('a')
    for i in soup_a:
        reg = re.compile(reg_key)
        x = reg.findall(str(i))[0]
        base_url3.append(base_url + str(x).split('"/')[1])
    return base_url3

def page_num(u):
    req_obj = requests.get(u)
    soup = BeautifulSoup(req_obj.text,'html.parser')
    b = soup.find_all('b')
    num = re.findall("\d+",str(b[-3]))[0]
    return num

if __name__ == '__main__':
    base_url = "https://www.gsmchoice.com"
    url_lev1 = "https://www.gsmchoice.com/zh-cn/catalogue/"
    #458个品牌
    url_lev2 = craw_lev1(base_url,url_lev1)
    #     #check每一二级页面的手机个数
    #     print (craw_lev1(base_url,url)[i],page_num(craw_lev1(base_url,url)[i]))
    #拿二级（手机品牌）分页 取三级（手机品牌-手机型号）
    
#     for iu in url_lev2:
#         print(iu) 

    with open("/Users/xuefeima/Desktop/python代码/craw_results.txt",'a' ,encoding="utf-8") as file:
#         url_lev2 = ["https://www.gsmchoice.com/zh-cn/catalogue/redmi/", "https://www.gsmchoice.com/zh-cn/catalogue/honor/"]
        cnt = 0
        for iu in url_lev2:
            url_lev3 = []
            cnt = cnt + 1
            print("开始处理第"+str(cnt)+"个品牌的手机")
            # real_url = https://www.gsmchoice.com/zh-cn/catalogue/huawei/models/80
            i = 0
            while 1:                
                real_url = iu + "models/" + str(i*40)
                staus_code = requests.get(real_url).status_code
                if(staus_code == 404):
                    i = 0
                    break
                url_lev3 += craw_lev2(real_url)
                i = i + 1
#                 print(str(staus_code)+"-成功爬取:"+real_url)

                time.sleep(1)
            print("开始写入"+iu.split("/")[-2]+"的数据")
            for m in url_lev3:
                file.write(m+"\n")
            print(iu.split("/")[-2]+"数据写入完成")
    print("程序执行结束")
