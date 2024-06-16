package main

import (
	"context"
	"fmt"
	"github.com/google/gousb"
	proto "github.com/zmey56/list-devices-protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os/exec"
	"strings"
)

type BluetoothDevice struct {
	Path      string
	VendorID  string
	ProductID string
	Type      string
}

type server struct {
	proto.UnimplementedDeviceServiceServer
}

func (s *server) ListDevices(context.Context, *proto.DeviceRequest) (*proto.DeviceList, error) {

	ctxUsb := gousb.NewContext()
	defer ctxUsb.Close()

	usbDevices, err := ctxUsb.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true
	})

	if err != nil {
		log.Printf("Error opening devices: %v\n", err)
	}
	// Close the devices when done.
	defer func() {
		for _, d := range usbDevices {
			d.Close()
		}
	}()

	var devices = &proto.DeviceList{}

	devices.Device = make([]*proto.Device, len(usbDevices))

	for _, dev := range usbDevices {
		desc := dev.Desc
		devices.Device = append(devices.Device, &proto.Device{
			DevicePath: fmt.Sprintf("%v", desc.Path),
			VendorId:   fmt.Sprintf("%v", desc.Vendor),
			ProductId:  fmt.Sprintf("%v", desc.Product),
			DeviceType: getDeviceType(desc),
		})
	}

	devicesBT, err := getBluetoothDevices()
	if err != nil {
		log.Fatalf("Error getting Bluetooth devices: %v", err)
	}

	for _, dev := range devicesBT {
		devices.Device = append(devices.Device, &proto.Device{
			DevicePath: dev.Path,
			VendorId:   dev.VendorID,
			ProductId:  dev.ProductID,
			DeviceType: dev.Type,
		})
	}

	return devices, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterDeviceServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func getDeviceType(desc *gousb.DeviceDesc) string {
	switch desc.Class {
	case gousb.ClassPerInterface:
		return "Defined at Interface Level"
	case gousb.ClassAudio:
		return "Audio"
	case gousb.ClassHID:
		return "Human Interface Device"
	case gousb.ClassPhysical:
		return "Physical"
	case gousb.ClassImage:
		return "Imaging"
	case gousb.ClassPrinter:
		return "Printer"
	case gousb.ClassMassStorage:
		return "Mass Storage"
	case gousb.ClassHub:
		return "Hub"
	case gousb.ClassData:
		return "Data"
	case gousb.ClassSmartCard:
		return "Smart Card"
	case gousb.ClassContentSecurity:
		return "Content Security"
	case gousb.ClassVideo:
		return "Video"
	case gousb.ClassPersonalHealthcare:
		return "Personal Healthcare"
	case gousb.ClassDiagnosticDevice:
		return "Diagnostic Device"
	case gousb.ClassMiscellaneous:
		return "Miscellaneous"
	default:
		return "Unknown"
	}
}

func getBluetoothDevices() ([]BluetoothDevice, error) {
	btInfo, err := runCommand("system_profiler", "SPBluetoothDataType")
	if err != nil {
		log.Fatalf("Error running system_profiler SPBluetoothDataType: %v", err)
	}

	lines := strings.Split(btInfo, "\n")
	var devices []BluetoothDevice
	device := BluetoothDevice{}
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		log.Println("fields:", fields, "|", fields[len(fields)-1])

		if strings.Contains(fields[0], "Address") && len(fields) > 1 {
			log.Println("NB fields:", fields, "|", fields[len(fields)-1])
			device.Path = fields[len(fields)-1]
		}
		if strings.Contains(fields[0], "Minor") && len(fields) > 1 {
			log.Println("NB fields:", fields, "|", fields[len(fields)-1])
			device.Type = fields[len(fields)-1]
			devices = append(devices, device)
			device = BluetoothDevice{}
		}
		if strings.Contains(fields[0], "Vendor") && len(fields) > 1 {
			log.Println("NB fields:", fields, "|", fields[len(fields)-1])
			device.VendorID = fields[len(fields)-1]
		}
		if strings.Contains(fields[0], "Product") && len(fields) > 1 {
			log.Println("NB fields:", fields, "|", fields[len(fields)-1])
			device.ProductID = fields[len(fields)-1]
		}

	}

	return devices, nil
}
