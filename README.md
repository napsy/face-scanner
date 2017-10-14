# Face Scanner

This is far from complete but is fun to use.

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

this Scans PNG images inside ``src-dir`` and puts found human faces inside ``dest-dir`` as PNG images.

Or, if you want to browse faces from the browser:

```
./face-scanner -web src-dir
```

and open ``http://localhost:4000``.


