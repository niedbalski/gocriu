package gocriu

import (
	proto "github.com/golang/protobuf/proto"
	"net"
	"os"
	"path/filepath"
	"strconv"
)

type Criu struct {
	UnixSocketPath string
	ImagesDirPath  string
	ShellJob       bool
	LogLevel       int32
}

const (
	DefaultUnixSocketPath string = "/var/run/criu-service.socket"
	DefaultImagesDirPath  string = "/tmp/dumps/"
	DefaultShellJob       bool   = true
	DefaultLogLevel       int32  = 4
)

func CriuClient(unixSocketPath string, imagesDirPath string, shellJob bool) (*Criu, error) {

	if unixSocketPath == "" {
		unixSocketPath = DefaultUnixSocketPath
	}

	if imagesDirPath == "" {
		imagesDirPath = DefaultImagesDirPath
		if _, err := os.Stat(imagesDirPath); os.IsNotExist(err) {
			return nil, err
		}
	}

	return &Criu{
		UnixSocketPath: unixSocketPath,
		ImagesDirPath:  imagesDirPath,
		ShellJob:       shellJob,
	}, nil
}

func (c *Criu) CriuRequest(requestType CriuReqType, pid int32) (*CriuResp, error) {
	dumpDir, err := c.GetDumpDir(pid)
	response := CriuResp{}

	if err != nil {
		return nil, err
	}

	dir, err := os.Open(dumpDir)
	if err != nil {
		return nil, err
	}

	fd := int32(dir.Fd())

	options := &CriuOpts{
		ImagesDirFd: &fd,
		Pid:         &pid,
		ShellJob:    &c.ShellJob,
		LogLevel:    &c.LogLevel,
	}

	conn, err := net.DialUnix("unixpacket", nil,
		&net.UnixAddr{
			c.UnixSocketPath,
			"unixpacket",
		})

	if err != nil {
		return nil, err
	}

	request := &CriuReq{
		Type: &requestType,
		Opts: options,
	}

	serialized, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	conn.Write(serialized)

	buf := make([]byte, 1024)

	if _, err = conn.Read(buf); err != nil {
		return nil, err
	}

	for idx, element := range buf {
		if element == 0 && buf[idx+1] == 0 {
			if requestType != CriuReqType(CriuReqType_DUMP) {
				buf = buf[0:idx]
			} else {
				buf = buf[0 : idx+1]
			}
			break
		}
	}

	err = proto.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}

	conn.Close()
	dir.Close()

	return &response, nil
}

func (c *Criu) GetDumpDir(pid int32) (string, error) {
	dumpPath := filepath.Join(c.ImagesDirPath, strconv.Itoa(int(pid)))

	if _, err := os.Stat(dumpPath); os.IsNotExist(err) {
		os.MkdirAll(dumpPath, 0755)
	}

	return dumpPath, nil
}

func (c *Criu) Dump(pid int32) (*CriuResp, error) {
	response, err := c.CriuRequest(CriuReqType(CriuReqType_DUMP), pid)

	if err != nil {
		return nil, err
	}

	return response, nil

}

func (c *Criu) Restore(pid int32) (*CriuResp, error) {
	response, err := c.CriuRequest(CriuReqType(CriuReqType_RESTORE), pid)

	if err != nil {
		return nil, err
	}

	return response, nil
}
