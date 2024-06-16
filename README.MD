Running the command::
```bash
go run main.go -path README.MD
```

Testing with the command:
```bash
grpcurl -plaintext -d '{}' localhost:50051 proto.DeviceService/ListDevices
```

Output:
```json
{
  "device": [
    {
      "devicePath": "/dev/mouse",
      "vendorId": 1234,
      "productId": 4321,
      "deviceType": "MOUSE"
    },
    {
      "devicePath": "/dev/keyboard",
      "vendorId": 4321,
      "productId": 1234,
      "deviceType": "KEYBOARD"
    },
    {
      "devicePath": "/dev/camera",
      "vendorId": 5678,
      "productId": 8765,
      "deviceType": "CAMERA"
    },
    {
      "devicePath": "/dev/microphone",
      "vendorId": 8765,
      "productId": 5678,
      "deviceType": "MICROPHONE"
    },
    {
      "devicePath": "/dev/speaker",
      "vendorId": 9876,
      "productId": 6789,
      "deviceType": "SPEAKER"
    }
  ]
}

```

Code generation is performed in the proto directory:
```bash
protoc --go_out=. --go-grpc_out=. service.proto  
```
