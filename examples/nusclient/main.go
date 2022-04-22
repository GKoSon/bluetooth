package main

// This example implements a NUS (Nordic UART Service) client. See nusserver for
// details.

import (
	"log"
	"time"

	"tinygo.org/x/bluetooth"
)

var (
	serviceUUID = bluetooth.ServiceUUIDNordicUART
	rxUUID      = bluetooth.CharacteristicUUIDUARTRX
	txUUID      = bluetooth.CharacteristicUUIDUARTTX
)

var adapter = bluetooth.DefaultAdapter

const target_name = "M_KOSON"

var runCnt int64 = 0

func oneloop() {
	// Enable BLE interface.
	//STEP 准备通过这个接口 选择不同HCI
	err := adapter.Enable()
	if err != nil {
		log.Printf("could not enable the BLE stack:%v", err.Error())
		return
	}

	// The address to connect to. Set during scanning and read afterwards.
	var foundDevice bluetooth.ScanResult

	// Scan for NUS peripheral.
	runCnt++
	log.Printf("Scanning...%d", runCnt)
	err = adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		if !result.AdvertisementPayload.HasServiceUUID(serviceUUID) {
			return
		}
		/*加强停止Scan函数的限制条件*/
		/*条件1--从机必须有名字 而且符合约定*/
		log.Printf("Scanned name 1-%s 2-%s", foundDevice.AdvertisementPayload.LocalName(), foundDevice.LocalName())
		if target_name != foundDevice.AdvertisementPayload.LocalName() {
			log.Printf("Failed name")
			return
		}

		/*条件2--从机连接成功*/
		d, e := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
		if e != nil {
			log.Printf("Failed Connect %v %v", d, e)
			return
		}

		foundDevice = result

		// Stop the scan.
		err := adapter.StopScan()
		if err != nil {
			// Unlikely, but we can't recover from this.
			log.Printf("Failed StopScan %v ", err.Error())
			/*
				    &测试
					执行hciconfig hci0 down 以后在run程序
					程序每次都是掉到这里 疯狂输出
				    这里没有办法 应该panic
			*/
			log.Panic("hciconfig hci0 up")
		}
	})

	if err != nil {
		log.Printf("Failed Scan %v ", err.Error())
		/*
			    &测试
				执行run程序 随后ctrl+c退出了
				下次继续run程序
				程序每次都是掉到这里 疯狂输出a scan is already in progress
				因为我没有defer 进程杀死了 没有人关闭蓝牙扫描
				那就在这里简单除暴的stop吧

				也有时候 Failed Scan Resource Not Ready
				需要去执行hciconfig hci0 up
		*/
		adapter.StopScan()
		return
	}

	// Found a NUS peripheral. Connect to it.
	device, err := adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		log.Printf("TG Failed to connect: %v ", err.Error())
		return
	}

	// Connected. Look up the Nordic UART Service.
	log.Printf("Discovering service...")
	services, err := device.DiscoverServices([]bluetooth.UUID{serviceUUID})
	if err != nil {
		println("Failed to discover the Nordic UART Service:", err.Error())
		return
	}
	service := services[0]

	// Get the two characteristics present in this service.
	chars, err := service.DiscoverCharacteristics([]bluetooth.UUID{rxUUID, txUUID})
	if err != nil {
		println("Failed to discover RX and TX characteristics:", err.Error())
		return
	}
	var rx bluetooth.DeviceCharacteristic
	var tx bluetooth.DeviceCharacteristic
	if chars[0].UUID() == txUUID {
		tx = chars[0]
		rx = chars[1]
	} else {
		tx = chars[1]
		rx = chars[0]
	}
	log.Printf("RX %v\r\n", rx)
	log.Printf("TX %v\r\n", tx)
	log.Printf("DiscoverCharacteristics:%+v\r\n", chars)

	// Enable notifications to receive incoming data.
	err = tx.EnableNotifications(func(value []byte) {
		log.Printf("PI recv %d bytes: %+v\r\n", len(value), value)
	})
	if err != nil {
		log.Printf("Failed EnableNotifications %+v\r\n", err.Error())
		return
	}

	log.Printf("Connected.When NODE disconnect.This pid while Exit\r\n")
	/*等待从机断开 PI从不发消息*/
	for {
		if !device.IsConnected() {
			log.Printf("device GoodBye\r\n")
			return
		}
		time.Sleep(time.Microsecond * 500)
	}
}

func main() {
	for {
		oneloop()
	}
}

/*
程序逻辑
main是死循环 也就是一旦oneloop()调用return那就继续下一次
正常retun是等待BLE链路的标记位

*/
