# FileClient
Downloads Files from a file server
This a client which downloads files using the mentioned server.
It downloads a file containing char 'A' on earlier position than other files.
In case several files have the 'A' char on the same the earliest position it downloads all of them.

**Strategy**

My strategy starts with an HTTP request to the download URL using the verb HEAD. One of the headers returned by some servers is Content-Length. This header specifies the file size in bytes. Once the file size is known, launch multiple Goroutines, each with its own data range to download. To begin the download, the goroutine will send an HTTP request to the URL using the GET verb.
The request's header will be Range. This header specifies how much of the file should be returned to the client. Once it has finished downloading the data, a Goroutine will send it back through the channel. Once all of the goroutines have completed, the data is combined and saved to a file.

Files in to be downloaded
![alt text](https://github.com/duncanodhis/FileClient/blob/3241fe341d076a2395e2f21bc8bd0fde5514ea53/Screenshot%20from%202022-10-14%2013-14-09.png)
