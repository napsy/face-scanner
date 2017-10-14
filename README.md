# Face Scanner

Deps:

```
go get github.com/napsy/go-opencv
go get github.com/disintegration/imaging
go get github.com/gorilla/websocket
```

Usage:

```
./face-scanner src-dir dest-dir
```

Scans PNG images inside ``src-dir`` and puts found human faces inside ``dest-dir`` as PNG images.
