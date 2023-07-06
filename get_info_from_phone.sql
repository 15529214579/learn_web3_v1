#coding:utf-8
import random
import re
from multiprocessing.pool import ThreadPool

import requests
import unicodedata
from bs4 import BeautifulSoup
from requests.packages import urllib3
urllib3.disable_warnings()
import pandas as pd

# 对机型的详细信息获取，采用pandas写入excel中

datalist = []

def pd_toExcel(data, fileName):  # pandas库储存数据到excel
    cpu = []
    freq = []
    corenum = []
    gpu = []
    for i in range(len(data)):
        cpu.append(data[i]["cpu"])
        freq.append(data[i]["freq"])
        corenum.append(data[i]["corenum"])
        gpu.append(data[i]["gpu"])
 
    dfData = {  # 用字典设置DataFrame所需数据
        'cpu':cpu,
        'freq': freq,
        'corenum': corenum,
        'gpu': gpu,
    }
    df = pd.DataFrame(dfData)  # 创建DataFrame
    df.to_excel(fileName, index=False)  # 存表，去除原始索引列（0,1,2...）

def get_soup(url_lev3):
    soup_one = "null"
    soup_two = "null"
    real_sout_li = []
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
    #源html提取 [soup_1,soup_2]
    req_obj = requests.get(url_lev3,headers=req_headers)
    req_obj.encoding = req_obj.apparent_encoding
    soup = BeautifulSoup(req_obj.text,'lxml')
    soup_t = soup.find_all(class_='PhoneData YesDict')
    real_sout_li.append(str(soup_t[0]))
    for sou in soup_t:
        html_text = str(sou)
        if '加速度计' in html_text:
            real_sout_li.append(html_text)
    #html转soup对
    soup_one = BeautifulSoup(real_sout_li[0],'lxml')
    if len(real_sout_li) == 1:
        return [soup_one,'null']
    soup_two = BeautifulSoup(real_sout_li[1],'lxml')

    return [soup_one,soup_two]

def craw_cell1(soup_1):
    #分别 正则提取 替换
    #key1
    item = re.sub(r'\t*|\n*|\[|\]','',unicodedata.normalize('NFKC', str(soup_1.find_all(class_='phoneCategoryName')).replace('\xa0','')))
    key_item = str(item).replace('<th class="phoneCategoryName">','').replace('</th>','')
    key_li = key_item.split(', ')
    #value1
    item_v = re.sub(r'\t*|\n*|\[|\]','',str(soup_1.find_all(class_='phoneCategoryValue')))
    item_v_li = item_v.split('</td>, <td class=')
    for i in item_v_li:
        if i.find("procSection") != -1:
            cpu = i.split("procSection")[1].split(">")[1].split("<")[0]
            freq = i.split("procSection")[1].split("Processor clock: ")[1].split("<")[0]
            corenum = i.split("procSection")[1].split("芯的数目: ")[1].split("<")[0] 
            gpu = i.split("procSection")[1].split("GPU: ")[1].split("<")[0]            
            datalist.append({"cpu":cpu, "freq":freq, "corenum":corenum,"gpu":gpu}) 
            
            print(i)
            print("CPU型号:",end = "")
            print(cpu)
            print("主频:",end = "")
            print(freq)
            print("核心数:",end = "")
            print(corenum)
            print("GPU型号:",end = "")
            print("gpu")
    


if __name__ == '__main__':
    path_ = '/Users/xuefeima/Desktop/python代码/info20230706.xlsx'
    _path = '/Users/xuefeima/Desktop/python代码/craw_results.txt'
    with open(_path,'r' ,encoding="utf-8") as _file, open(path_, 'a', encoding="utf-8") as file_:
        _file = ["https://www.gsmchoice.com/zh-cn/catalogue/samsung/gtc3322lafleur/","https://www.gsmchoice.com/zh-cn/catalogue/samsung/galaxys23/"]
        list_size = 0
        for url in _file:
            if 200 == requests.get(url).status_code:
                print('开始爬取: '+url)
                r_a = craw_cell1(get_soup(url)[0])
                list_size = list_size + 1
                if list_size > 0:
                    pd_toExcel(datalist, path_)
                    datalist = []
    #                 result = dict(list(r_a.items()) + list(r_b.items()))
                    
                print('成功爬取: '+url)
#             break
    file_.close()
    print('结束爬取，写入文件完成！: '+path_)
# #  采用python内部线程池加速，但是要解决datalist写入问题，先不采用，且并发爬数据ip有被封的风险
# if __name__ == '__main__':
#     _path = '/Users/xuefeima/Desktop/python代码/craw_results.txt'
#     #设置线程并行
#     #遍历url 爬取
#     urls = []
#     with open(_path,'r' ,encoding="utf-8") as _file:
#         for url in _file:
#             urls.append(url)
#     _file.close()
    
#     pool = Pool(processes=10)
#     result = pool.map(thread_job, urls)
#     pool.close()        # 关闭进程池，不再接受新的进程
#     pool.join()    