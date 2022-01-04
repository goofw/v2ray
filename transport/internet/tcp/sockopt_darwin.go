// +build darwin

package tcp

import (
    "fmt"
    "os/exec"
    "strconv"
    "strings"

    "github.com/v2fly/v2ray-core/v5/common/net"
    "github.com/v2fly/v2ray-core/v5/transport/internet"
)

// GetOriginalDestination from tcp conn
func GetOriginalDestination(conn internet.Connection) (net.Destination, error) {
    addr, _ := lookup(conn.RemoteAddr().(*net.TCPAddr))
    dest := net.TCPDestination(net.ParseAddress(addr.IP.String()), net.Port(addr.Port))
    return dest, nil
}

func lookup(addr *net.TCPAddr) (*net.TCPAddr, error) {
    var new_addr net.TCPAddr
    
    out, _ := exec.Command("sudo", "-n", "/sbin/pfctl", "-s", "state").Output()
    for _, line := range strings.Split(string(out), "\n") {
        if strings.Contains(line, "ESTABLISHED:ESTABLISHED") {
            if strings.Contains(line, fmt.Sprintf("%s:%d", addr.IP.String(), addr.Port)) {
                fields := strings.Fields(line)
                if len(fields) > 4 {
                    addr := strings.Split(fields[4], ":")
                    new_addr.IP = net.ParseIP(addr[0])
                    new_addr.Port, _ = strconv.Atoi(addr[1])
                    break
                }
            }
        }
    }

    return &new_addr, nil
}
