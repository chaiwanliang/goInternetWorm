package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

//
func HttpGetDB(url string)(result string,err error){
	resp,err1 :=http.Get(url)
	if err1!=nil{
		err=err1
		return
	}
	defer resp.Body.Close()

	buf :=make([]byte,4096)
	for{
		n,err2 :=resp.Body.Read(buf)
		if n==0{
			break
		}
		if err2!=nil && err2!=io.EOF {
			err=err2
			return
		}
		result+=string(buf[:n])
	}
	return
}

// 保存提取的有用信息
func SaveFile(index int,filName,fileScore,peopleName [][]string){
	file,err :=os.Create("豆瓣电影爬虫/data/第"+strconv.Itoa(index)+"页.txt")
	if err!=nil{
		fmt.Println("os.Create err:",err)
		return
	}
	defer file.Close()

	n :=len(filName) // 得到条目数
	// 先打印 抬头 电影名称 评分 评分人数
	file.WriteString("电影名称"+"\t\t\t"+"评分"+"\t\t"+"评分人数"+"\n")
	for i:=0;i<=n;i++{
		file.WriteString(filName[i][1]+"\t\t\t"+fileScore[i][1]+"\t\t"+peopleName[i][1]+"\n")
	}
}

// 抓取单个页面操作
func SpiderPageDB(index int,page chan int){
	url :="https://movie.douban.com/top250?start="+strconv.Itoa((index-1)*25)+"&filter="
	fmt.Println(url)
	result,err :=HttpGetDB(url)
	if err!=nil{
		fmt.Println("HttpGetDB err:",err)
		return
	}
	fmt.Println("result= ",result)
	// 解析、编译正则表达式 —— —— 电影名称
	ret1 :=regexp.MustCompile(`<img width="100" alt="(.*?)"`)
	// 提取需要信息
	filName :=ret1.FindAllStringSubmatch(result,-1)

	// 解析、编译正则表达式 —— —— 分数
	ret2 :=regexp.MustCompile(`<span class="rating_num" property="v:average">(?s:(.*?))</span>`)
	// 提取需要信息
	fileScore :=ret2.FindAllStringSubmatch(result,-1)

	// 解析、编译正则表达式 —— —— 评分人数
	ret3 :=regexp.MustCompile(`<span>(.*?)人评价</span>`)
	// 提取需要信息
	peopleName :=ret3.FindAllStringSubmatch(result,-1)

	// 提取到有用的信息 封装到文件中
	SaveFile(index,filName,fileScore,peopleName)

	page<-index
}


func toWork(start,end int){

	fmt.Printf("正在爬取第%d页到%d页...\n",start,end)

	page :=make(chan int)
	for i:=start;i<=end;i++{
		go SpiderPageDB(i,page)
	}
	for i:=start;i<=end;i++{
		fmt.Printf("第%d页爬取完后",<-page)
	}

}

func main(){
	// 指定爬取起始、终止页
	var start,end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止页(>=start):")
	fmt.Scan(&end)

	// 创建数据文件夹
	if err :=os.Mkdir("豆瓣电影爬虫/data",os.ModePerm);err !=nil{
		fmt.Println("os.Mkdir err :",err)
	}

	toWork(start,end)
}