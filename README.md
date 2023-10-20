# tarolas 

Lightweight file storage service written in Go.

## Installation

Download sources from this repository and build **tarolas** binary.

```shell
$ go build wisbery.com/tarolas
```
   
## Running

Tarolas server requires single configuration file given in command line:

```shell
$ tarolas --config my_tarolas_config.json
``` 

Configuration options are described in details in **Configuration** section.

The simplest configuration file may look like the following:

(Windows)
```json
{
  "serverPort": 15000,
  "rootDirectory": "C:\\Temp\\tarolas"
}
```
(Linux)

```json
{
  "serverPort": 15000,
  "rootDirectory": "/home/john/tarolas"
}
```    

You can find example configuration file in **config** directory of this repository.

After starting tarolas server you will see the similar messages in the console:

(Windows)

```shell
tarolas v0.0.8
  >>    server port : 15000
  >> root directory : C:\Temp\tarolas
```

(Linux)

```shell
tarolas v0.0.8
  >>    server port : 15000
  >> root directory : /home/john/tarolas
```      

The **tarolas** file server is ready to serve your files!

## Configuration

- port number the sever will be listening on
- root directory for storing files
- public key for JWT token verification

## Functionality

### Directories

- list directory tree,
- list directory content,
- create directory,
- delete directory,
- move directory,
- copy directory.

### Files

- create file,
- append file,
- delete file,
- move file,
- copy file,
- share file,
- read file.

## Security

Directories and files may be accessed without any restrictions.
Usually this is not a good idea. That's why any operation may be restricted 
only to user who has required rights granted.

## License

Licensed under either of

- [MIT license](https://opensource.org/licenses/MIT) ([LICENSE-MIT](https://github.com/wisbery/tarolas/blob/main/LICENSE-MIT)), or
- [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0) ([LICENSE-APACHE](https://github.com/wisbery/tarolas/blob/main/LICENSE-APACHE))

at your option.

## Contribution

All contributions intentionally submitted for inclusion in the work by you,
shall be dual licensed as above, without any additional terms or conditions.
