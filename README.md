# Identicon

A simple implementation to generate Iedenticons

## Install

```sh
go get -u github.com/hambletor/identicon
```
## Examples

To create a simple 250 x 250 pixel 5 x 5 identicon 

```go
 i, err := identicon.New("filename")
 if err != nil{
     log.Printf("unable to create identicon %v\n")
 }
 i.SavePNG()
 ```

 Now let's change the background color to something other than white

 ```go
 i, err := identicon.New("filename", identicon.WithBackgroundColor(color.Black))
 if err != nil{
     log.Printf("unable to create identicon %v\n")
 }
 i.SavePNG()
 ```

 Now for a fully customized identicon

 ```go
 fg := color.RGBA{R: 122, G: 67, B: 210, A: 255}
 i, err := identicon.New("filename",
     identicon.WithBackgroundColor(color.Black),
     identicon.WithForegroundColor(fg),
     identicon.WithPixels(350),
     identicon.WithSize(7)
     )
 if err != nil{
     log.Printf("unable to create identicon %v\n")
 }
 i.SaveJPEG()
 ```