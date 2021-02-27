package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tidwall/geojson"
	"github.com/tidwall/gjson"
	"github.com/tormoder/fit"
)

func FitFromGeojson(jsonData string) {
	g, _ := geojson.Parse(string(jsonData), &geojson.ParseOptions{RequireValid: false})
	fmt.Println(g)
	h := fit.NewHeader(fit.V20, true)
	file, err := fit.NewFile(fit.FileTypeCourse, h)
	if err != nil {
		fmt.Println(err)
	}
	fileIdMesg := fit.NewFileIdMsg()
	fileIdMesg.Manufacturer = fit.ManufacturerGarmin
	fileIdMesg.Type = fit.FileTypeCourse
	fileIdMesg.Product = 1
	fileIdMesg.SerialNumber = uint32(1234)
	file.FileId = *fileIdMesg

	devMsg := fit.NewDeveloperDataIdMsg()
	devMsg.ApplicationId = []byte{
		0x1, 0x1, 0x2, 0x3,
		0x5, 0x8, 0xD, 0x15,
		0x22, 0x37, 0x59, 0x90,
		0xE9, 0x79, 0x62, 0xDB,
	}

	course, err := file.Course()
	fmt.Println(course)
	if err != nil {
		fmt.Printf("course: got error, want none; error is: %v", err)
	}

	courseMsg := fit.NewCourseMsg()
	courseMsg.Name = "tempname"
	courseMsg.Sport = fit.SportHiking
	course.Course = courseMsg

	geometry := gjson.Get(g.JSON(), "geometry.coordinates")
	//len := len(geometry.Array())
	//coursePoints := make([]*fit.CoursePointMsg, 0)
	for _, v := range geometry.Array() {
		point := v.Array()
		long := point[0].Float()
		lat := point[1].Float()
		cp := fit.NewCoursePointMsg()
		cp.PositionLong = fit.NewLongitudeDegrees(long)
		cp.PositionLat = fit.NewLatitudeDegrees(lat)
		course.CoursePoints = append(course.CoursePoints, cp)
	}
	fmt.Println(course)
	//course.CoursePoints = coursePoints
	outBuf := &bytes.Buffer{}
	f, err := os.Create("/tmp/test.fit")
	err = fit.Encode(outBuf, file, binary.LittleEndian)
	n2, err := f.Write(outBuf.Bytes())
	fmt.Println(err)
	fmt.Printf("wrote %d bytes\n", n2)

	t, err := fit.Decode(bytes.NewReader(outBuf.Bytes()))
	if err != nil {
		fmt.Printf("decode: got error, want none; error is: %v", err)
	}
	fmt.Println(t)

	f2, err := os.Create("/tmp/test2.fit")
	w := bufio.NewWriter(f2)
	err = fit.Encode(w, file, binary.LittleEndian)
	f2.Close()
	fmt.Println(err)

}

func main() {
	data, err := ioutil.ReadFile("feature.json")
	if err != nil {
		fmt.Println("File reading error", err)
		return
	}
	jsonData := string(data)
	FitFromGeojson(jsonData)
}
