package networks

import (
	"errors"
	"net"
)

func isValidIP(ip net.IP) bool {
	return ip != nil && !ip.IsLoopback() && ip.To4() != nil
}

// GetLocalIP retorna o primeiro endereço IP (no formato de string) das interfaces de rede não-loopback no sistema.
func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			if isValidIP(ipnet.IP) {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("nenhum IP válido encontrado")
}
