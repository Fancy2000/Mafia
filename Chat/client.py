import socket
import threading
import pyaudio
import numpy as np
import cv2
import os


class Users:
    mode = "speaking"
    name = ""

    def __init__(self, mode, name):
        self.mode = mode
        self.name = name


class Client:
    def __init__(self):
        self.s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        while 1:
            try:
                self.target_ip = input('Enter IP address of server --> ')
                self.target_port = int(input('Enter target port of server --> '))
                self.username = input('Enter your name --> ')

                self.s.connect((self.target_ip, self.target_port))

                self.s.send(self.username.encode())
                err = self.s.recv(1024).decode()
                while err != "OK":
                    print(err, flush=True)
                    self.username = input('Enter your name --> ')
                    self.s.send(self.username.encode())
                    err = self.s.recv(1024).decode()
                self.mode = True

                break
            except:
                print("Couldn't connect to server", flush=True)

        chunk_size = 1024  # 512
        audio_format = pyaudio.paInt16
        channels = 1
        rate = 20000
        self.disconnect = False
        # initialise microphone recording
        self.p = pyaudio.PyAudio()
        self.playing_stream = self.p.open(format=audio_format, channels=channels, rate=rate, output=True,
                                          frames_per_buffer=chunk_size)
        self.recording_stream = self.p.open(format=audio_format, channels=channels, rate=rate, input=True,
                                            frames_per_buffer=chunk_size)

        print("Connected to Server", flush=True)


        # start threads
        self.receive_thread = threading.Thread(target=self.receive_server_data)
        self.receive_thread.start()
        self.send_data_to_server()

    def receive_server_data(self):
        while True:
            try:
                data = self.s.recv(1024)
                self.playing_stream.write(data)
            except:
                pass

    def print_users(self, info):
        print("----------------------------------------", flush=True)
        print("Clients information:", flush=True)
        i = 1
        for user in info:
            print(str(i) + ") " + "User " + user.name + " is " + user.mode, flush=True)
            i += 1
        print("----------------------------------------", flush=True)

    def send_data_to_server(self):
        cv2.imshow(self.username, np.zeros(shape=(1000, 1000, 3), dtype=np.uint8))
        while True:
            try:
                key = cv2.waitKey(1)
                if key == ord("m"):
                    self.s.send(b"mute")
                    self.recording_stream.stop_stream()
                    self.mode = False
                if key == ord("s"):
                    self.s.send(b"speaking")
                    self.recording_stream.start_stream()
                    self.mode = True
                if key == ord("i"):
                    self.s.send(b"information")
                if key == ord("q"):
                    self.s.send(b"exit")
                    os._exit(0)
                if self.mode:
                    data = self.recording_stream.read(1024)
                    self.s.sendall(data)
            except:
                pass


client = Client()

