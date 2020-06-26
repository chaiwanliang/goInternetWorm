package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HttpGet(url string)(result string,err error){
	resp,err1 :=http.Get(url)
	if err1!=nil{
		// 将封装函数内部的错误,传出给调用者
		err=err1
		return
	}
	defer resp.Body.Close()

	// 循环读取网页数据,传出给调用者
	buf :=make([]byte,4096)
	for{
		n,err2 :=resp.Body.Read(buf)
		if n==0{
			break
		}
		if err2!=nil && err2!=io.EOF{
			err =err2
			return
		}
		// 累加每一次循环读到的buf数据,存入result，一次性返回
		result +=string(buf[:n])
	}
	return
}

// 单个页面爬取操作
func SpiderPage(index int,page chan int){
	url :="https://tieba.baidu.com/f?kw=%E7%BB%9D%E5%9C%B0%E6%B1%82%E7%94%9F&ie=utf-8&pn="+strconv.Itoa((index-1)*50)
	result,err :=HttpGet(url)
	if err!=nil{
		fmt.Println("HttpGet err：",err)
		return
	}
	// fmt.Println("result :",result)
	// 将读到的整网页数据，保存成一个文件
	file,err :=os.Create("data/第"+strconv.Itoa(index)+" 页"+".html")
	if err!=nil{
		fmt.Println("Create err:",err)
		return
	}
	// 写文件
	file.WriteString(result)
	// 保存好一个文件,关闭一个文件
	file.Close()

	page<-index  // 与主go程完成同步
}

// 爬取页面操作
func working(start,end int)  {

	fmt.Printf("正在爬取第%d页到%d页...\n",start,end)
	page :=make(chan int)
	// 循环爬取每一页的数据
	for i:=start;i<=end;i++{
		go SpiderPage(i,page)
	}

	for i:=start;i<=end;i++{
		fmt.Printf("爬取第 %d 页完成\n",<-page)
	}
}

func main(){
	// 指定爬取起始、终止页
	var start,end int
	fmt.Println("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Println("请输入爬取的终止页(>=start):")
	fmt.Scan(&end)

	// 获取当前目录
	// dir,err :=os.Getwd()
	// if err!=nil{
	// 	fmt.Println("os.Getwd() err:",err)
	// }
	// fmt.Println("当前目录：",dir)

	// 创建数据文件夹
	if err :=os.Mkdir("百度贴吧爬虫/data",os.ModePerm);err !=nil{
		fmt.Println("os.Mkdir err :",err)
	}

	working(start,end)
}