import socket
import ssl
import struct


def https_get(socks_port, ip, port, host, path):
    sock = socket.create_connection(('127.0.0.1', socks_port), timeout=30)
    sock.sendall(b'\x05\x01\x00')
    assert sock.recv(2) == b'\x05\x00'
    sock.sendall(b'\x05\x01\x00\x01' + socket.inet_aton(ip) + struct.pack('!H', port))
    reply = sock.recv(10)
    assert reply[:2] == b'\x05\x00', reply
    context = ssl.create_default_context()
    context.check_hostname = False
    context.verify_mode = ssl.CERT_NONE
    tls = context.wrap_socket(sock, server_hostname=host)
    tls.sendall('GET {0} HTTP/1.1\r\nHost: {1}\r\nConnection: close\r\n\r\n'.format(path, host).encode())
    response = b''
    while True:
        chunk = tls.recv(4096)
        if not chunk:
            break
        response += chunk
    tls.close()
    return response.decode(errors='replace')
