# Video Morpher

WebServer written in Go that allows users to upload video files and apply effects. The server uses FFMPEG to handle video operations.

It can be accessed on **http://15.188.106.46:8080/**

## Prerequisites

- To run the WebServer, [FFmpeg](https://ffmpeg.org/) is required to be installed on the machine.

## Running the Web Server from source code

1. Clone the repository:
    ```sh
    git clone https://github.com/davidpalves06/VideoMorpher.git
    cd VideoMorpher
    ```

2. Build the executable:
    ```sh
    go build cmd/videomorpher/videomorpher.go
    ```

3. Run the executable:
    ```
    ./videomorpher
    ```

4. Access web server from browser:
   ```
    http://localhost:8080/
    ```

## Running the Web Server from docker image
1. Clone the repository:
    ```sh
    git clone https://github.com/davidpalves06/VideoMorpher.git
    cd VideoMorpher
    ```

2. Build the docker image:
    ```sh
   docker build . -t videomorpher:alpha -f docker/Dockerfile 
    ```

3. Run the docker image:
    ```sh
   docker run --net=host videomorpher:alpha
    ```
    
4. Access web server from browser:
   ```
    http://localhost:8080/
    ```

## Configuration
By default, the server will look for the configuration on **_config.json_**.  If you want to change this, you can pass -configFile as an argument on server startup. The server will look for the configuration on the file defined after the flag. Just be aware that the configuration is always on json format.

### Configuration parameters

#### Configuration Example
```
{
    "server": {
        "host": "0.0.0.0",
        "port": 8080
    },
    "log" :{
        "level": 0
    },
    "uploadDir": "./uploads/"
}
```

 #### Server configurations
 With this configurations, the user can set the host and the port in which the socket should listen to requests.

 #### Log configurations 
 Here the user can define the minimum level of logs that the server should generate.  
 - 0 - DEBUG
 - 1 - INFO
 - 2 - WARN
 - 3 - ERROR  

 If not specified, the default log level is 1 (INFO)

 #### Upload Directory configuration
 This configuration always the user to change the folder where uploads are stored.
 If not specified, the default upload directory is _uploads_.

## Next Steps
- Add new effects
- HTTPS
