# FastShare



#### FastShareはnode.jsを使ってプライベートネットワーク上に指定したディレクトリを公開します。



##### ディレクトリの公開

```
node fastShare.js -d .\my_favorite_cats\
	2222/2/2  2:22:22: Publish: .\my_favorite_cats\	
	ipv4: 192.168.xxx.yyy
```

`-d`オプションで公開するディレクトリを指定します。



##### ディレクトリのダウンロード

`fastShareClient.exe http://192.168.xxx.yyy -d my_favorite_cats`

```
fastShareClient.exe http://192.168.xxx.yyy -d my_favorite_cats

[STATS]:Requesting directory information
[STATS]:Creating each directory
[STATS]:Download
[STATS]:Done
```

> `my_favorite_cats`をダウンロード

-dで指定したディレクトリをダウンロードします。

