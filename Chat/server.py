import socket
import threading
import sys

wrong_name = "sorry this name is already taken"
accept_name = "OK"


class Users:
    mode = "speaking"
    sock = 0
    name = ""

    def __init__(self, mode, name, sock):
        self.mode = mode
        self.sock = sock
        self.name = name


class Users_send:
    mode = "speaking"
    name = ""

    def __init__(self, mode, name):
        self.mode = mode
        self.name = name


class Server:
    def __init__(self):
        self.ip = socket.gethostbyname(socket.gethostname())
        while 1:
            try:
                self.port = 5454

                self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
                self.s.bind((self.ip, self.port))

                break
            except:
                print("Couldn't bind to that port", flush=True)
        self.usernames = set()
        self.connections = []
        self.accept_connections()

    def accept_connections(self):
        self.s.listen(100)

        print('Running on IP: ' + self.ip, flush=True)
        print('Running on port: ' + str(self.port), flush=True)

        while True:
            c, addr = self.s.accept()
            username = c.recv(1024).decode("utf-8")

            while username in self.usernames:
                c.send(wrong_name.encode())
                username = c.recv(1024).decode("utf-8")

            c.send(accept_name.encode())
            print("New client " + "'" + username + "' " + "connected: " + "speaking", flush=True)

            self.usernames.add(username)
            self.connections.append(Users(mode="speaking", name=username, sock=c))
            self.list_of_users()
            threading.Thread(target=self.handle_client, args=(c, addr,)).start()

    def broadcast(self, sock, data):
        for client in self.connections:
            if client.sock != self.s and client.sock != sock:
                try:
                    client.sock.send(data)
                except:
                    pass

    def change_mode(self, c, type):
        for user in self.connections:
            if user.sock == c:
                user.mode = type

    def disconnect_informating(self, sock):
        tmp = 0
        for i in self.connections:
            if i.sock == sock:
                print("Client " + "'" + i.name + "'" + " disconnected ", flush=True)
                self.usernames.remove(i.name)
                tmp = i
                break
        self.connections.remove(tmp)
        self.list_of_users()

    def list_of_users(self):
        print("", flush=True)
        print("----------------------------------------", flush=True)
        print("Clients information", flush=True)
        i = 1
        for user in self.connections:
            print(str(i) + ") " + "User " + user.name + " is " + user.mode, flush=True)
            i += 1
        print("----------------------------------------\n", flush=True)

    def handle_client(self, c, addr):
        while 1:
            try:
                data = c.recv(1024)
                if sys.getsizeof(data) < 70:
                    udata = data.decode()
                    if udata == "mute":
                        self.change_mode(c, "mute")
                    elif udata == "speaking":
                        self.change_mode(c, "speaking")
                    elif udata == "information":
                        self.list_of_users()
                    elif udata == "exit":
                        self.disconnect_informating(c)
                        c.close()
                else:
                    self.broadcast(c, data)

            except socket.error:
                c.close()


server = Server()
