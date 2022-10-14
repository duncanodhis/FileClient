package main

import (
    "fmt"
     "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "regexp"
    "time"
    "runtime"
    "strings"
    "os"
)

const (
    url     = "http://localhost:8080/"

    singleChunkBufferLenght = 10000//set the size of a single chunk
     maxOffsetAsChunks  = 10000//sets the size of max chunk
)


type Part struct {
    Data  []byte
    Index int
}
//responsible for error checking
func check(e error) {
    if e != nil {
        panic(e)
    }
}
//fetches name of the file to be downloaded
func getFileNames(ur string) []string{

    files :=[]string{}
    client := &http.Client{}
    resp, err := client.Get(ur)
    if err != nil {
        // code to restart download
        return files
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    str1 := string(body)//converts body to string
    re := regexp.MustCompile(`<a.*?>(.*)</a>`)//removes <a> tag,in my case the files are in a tag
    submatchall := re.FindAllStringSubmatch(str1, -1)
    
    for _, element := range submatchall {
        // fmt.Println("At index", index, "value is", element[1])
        files=append(files,element[1])
        
    }
    // fmt.Println(files[:])
    return files

}
//gets full url name
func fullURL(files []string ,ur string) []string{
    fullurl :=[]string{}
    var str1 string
    str1 = ur
    for _, file := range files {

        fullurl =append(fullurl,(str1+file))

        
    }
    // fmt.Println(fullurl)
    return fullurl
}


// TotalCost returns download duration.
func  TotalCost() int64 {
    start := time.Now()
// Code to measure
    duration := time.Since(start)
    return duration.Nanoseconds()
}
//Will be passed through a go routine to ensure that files are downloaded in chunks
func download(index, size int, c chan Part,ur string,names string,numberOfChunksToRead int ) {


    client := &http.Client{}

    start := index * size
    dataRange := fmt.Sprintf("bytes=%d-%d", start, start+size-1)
    // var numberOfChunksToRead = fileSize/size

    if index == numberOfChunksToRead - 1 {
        dataRange = fmt.Sprintf("bytes=%d-", start)
    }

    log.Println(dataRange)

    req, err := http.NewRequest("GET", ur, nil)

    if err != nil {
        // code to restart download
        return
    }

    req.Header.Add("Range", dataRange)//gets data range

    resp, err := client.Do(req)

    if err != nil {
        // code to restart download
        return
    }

    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    // fmt.Println(string(body))

    if err != nil {
        // code to restart download
        return
    }

    c <- Part{Index: index, Data: body}



}
//gets default concurrency value
func getDefaultConcurrency() int {

    c := int(runtime.NumCPU() * 3)

    // Set default max concurrency to 20.
    if c > 20 {
        c = 20
    }

    // Set default min concurrency to 4.
    if c <= 2 {
        c = 4
    }

    return c
}
//gets default chunk size
func getDefaultChunkSize(totalSize, min, max, concurrency int) int {

    cs := totalSize / concurrency

    // if chunk size >= 102400000 bytes set default to (ChunkSize / 2)
    if cs >= 102400000 {
        cs = cs / 2
    }

    // Set default min chunk size to 2m, or file size / 2
    if min == 0 {

        min = 2097152

        if min >= totalSize {
            min = totalSize / 2
        }
    }

    // if Chunk size < Min size set chunk size to min.
    if cs < min {
        cs = min
    }

    // Change ChunkSize if MaxChunkSize are set and ChunkSize > Max size
    if max > 0 && cs > max {
        cs = max
    }

    // When chunk size > total file size, divide chunk / 2
    if cs >= totalSize {
        cs = totalSize / 2
    }

    return cs
}
//client side ,responsible for downloading files
func fileClient(ur string ,filename string){
    

    client := &http.Client{}

    req, err := http.NewRequest("HEAD", ur, nil)

    if err != nil {
        log.Fatal(err)
    }

    resp, err := client.Do(req)

    if err != nil {
        log.Fatal(err)
    }


    log.Println("Headers : ", resp.Header["Content-Length"])

    var singleChunkBufferLenght = 10000
    var maxOffsetAsChunks = 10000
    
    header, ok := resp.Header["Content-Length"]
    if !ok {
        fmt.Println("error")
    }
    fileSize, err := strconv.Atoi(header[0])
    check(err)
    // fmt.Println(maxOffsetAsChunks)
    var numberOfChunksToRead = getDefaultChunkSize(fileSize,singleChunkBufferLenght,maxOffsetAsChunks,getDefaultConcurrency())

    var size int
    size = fileSize / numberOfChunksToRead
    results := make(chan Part, numberOfChunksToRead)
    parts := make([][]byte, numberOfChunksToRead)

    for i :=0 ;i<numberOfChunksToRead ; i++ {
       parts[i] = make([]byte, numberOfChunksToRead)
    }
    

    for i := 0; i < numberOfChunksToRead; i++ {
    
        go download(i, size, results,ur,filename,numberOfChunksToRead)
      
    }

    counter := 0

    for part := range results {
        counter++

        parts[part.Index] = part.Data
        if counter == numberOfChunksToRead {
            // fmt.Println("taking a break")
            break
        }
    }

    file := []byte{}

    for _, part := range parts {
        file = append(file, part...)
    }

    // Set permissions accordingly, 0700 may not
    // be the best choice
    err = ioutil.WriteFile(filename, file, 0700)

    if err != nil {
        log.Fatal(err)
    }

    

   
}
//find the smallest positive integer in the array,-1 means A does not exist in the file
func findMinElement(arr []int) int {
   min_num := arr[0]

   for i:=0; i<len(arr); i++{
    if arr[i] > -1 {

    if arr[i] < min_num {
         min_num = arr[i]
      }

    }
      
   }
   return min_num
}
//removes unwanted files
func removeFile(value int,names []string){
 
   for _, n := range names {
      // fmt.Println("reading ",n)
        b, err := ioutil.ReadFile(n) // just pass the file name
        
        if err != nil {
            fmt.Print(err)
             }

        str := string(b) // convert content to a 'string'
        res1 := strings.Index(str, "A") //gets earliest position
       
        if res1 != value{
        e := os.Remove(n)//removes the files with no earliest position
        if e != nil {
            log.Fatal(e)
        }
     }else{
        fmt.Println("File downloaded is: ",n)
     }
   }
}

func main() {

        var l = len(getFileNames(url))//gets the number off urls
        names := make([]string, l)
        names = getFileNames(url)//names of the urls
        urls := make([]string, l)
        urls=fullURL(names ,url)
        // fmt.Println(urls)
        // fmt.Println(names)
        array := make([]int, l)
        

        for i := 0; i < len(urls) && i < len(names);i++ {
            fmt.Println(names[i])
            fileClient(urls[i] ,names[i])
            fmt.Println("Duration in Nanoseconds = ",TotalCost())
           
        }
        for i, n := range names {
        // fmt.Println("reading ",n)
        b, err := ioutil.ReadFile(n) //  pass the file name
        
        if err != nil {
            fmt.Print(err)
             }

        str := string(b) // convert content to a 'string'
        res1 := strings.Index(str, "A") 
       
        array[i]=res1

        }
        removeFile(findMinElement(array),names)    
}




    

    

   